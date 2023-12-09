package json

import (
	"crypto"
	"crypto/sha256"
	"encoding/json"

	rsainternal "github.com/zsrv/supermicro-product-key/pkg/crypto/rsa"
	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
)

// LicenseSigningPublicKey is used to verify JSON format licenses
// that have been digitally signed by Supermicro.
var LicenseSigningPublicKey = []byte(`-----BEGIN RSA PUBLIC KEY-----
MIIBCAKCAQEAvDb4MxWw/FTi8pscP6S2YAl/3gmVOj/StG0lu3PdCSmdpzmzbOU9
KBS3t0yPZ0ynUQj/qXOwaVLvBJ+uCE0pGIRWkzBersVUzmXXN8Dza5yOzlLIdsVn
amUrKcRHgC+otRE/gnCxIiioacy9TkA96otbAvztCl1j1W8oCixazpfwZrayy12y
CcOyquZr3prngLCgOWa9e9cLSekKuvYXPKPC0CogLz0ueg+y+gcWGHGwbMtERLGf
WYDrXD1mdlV0EL4A5H4v9bqzn0yHe9dIb7+tdMHiN+qLjMzcVSzwfXu7Abk7rPOz
wBY6gMvfhfaOTj8aU5JFQ3DEBdTPOwgtSwIBAw==
-----END RSA PUBLIC KEY-----`)

type License struct {
	ProductKey ProductKey
}

type ProductKey struct {
	Node      Node
	Signature []byte
}

type Node struct {
	LicenseID   string
	LicenseName string
	CreateDate  string
}

func ParseLicense(productKey string) (License, error) {
	var license License
	err := json.Unmarshal([]byte(productKey), &license)
	if err != nil {
		return License{}, err
	}

	return license, nil
}

// Verify verifies the License signature.
//
// An error is returned if the License data has been altered, the macAddress
// is not the one that was in the License message that was signed, or
// LicenseSigningPublicKey is not the one that was used to sign the License data.
func (l *License) Verify(macAddress netinternal.HardwareAddr) error {
	publicKey, err := parseRSAPublicKey(LicenseSigningPublicKey)
	if err != nil {
		return err
	}

	return l.VerifyWithKey(publicKey, macAddress)
}

// VerifyWithKey verifies the License signature against publicKey.
//
// An error is returned if the License data has been altered, the macAddress
// is not the one that was in the License message that was signed, or the
// publicKey is not the one that was used to sign the License data.
func (l *License) VerifyWithKey(publicKey *rsainternal.PublicKey, macAddress netinternal.HardwareAddr) error {
	message := []byte(macAddress.String() +
		l.ProductKey.Node.LicenseID +
		l.ProductKey.Node.LicenseName +
		l.ProductKey.Node.CreateDate)
	hashed := sha256.Sum256(message)

	err := rsainternal.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], l.ProductKey.Signature)
	if err != nil {
		return err
	}

	return nil
}
