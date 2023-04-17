package json

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"github.com/zsrv/supermicro-product-key/pkg/oob"
	"strings"
)

// licenseSigningPublicKey is used to validate JSON format licenses
// that have been digitally signed by Supermicro.
var licenseSigningPublicKey = []byte(`-----BEGIN RSA PUBLIC KEY-----
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
	Signature string
}

type Node struct {
	LicenseID   string
	LicenseName string
	CreateDate  string
}

// VerifySignature returns nil if the license signature is valid when the given
// macAddress is used as input in the message being hashed and signed, or an error
// if one occurs.
func (l *License) VerifySignature(macAddress string) error {
	macAddress, err := oob.NormalizeMACAddress(macAddress)
	if err != nil {
		return err
	}
	macAddress = strings.ToUpper(macAddress)

	publicKey, err := parseRSAPublicKey(licenseSigningPublicKey)
	if err != nil {
		return err
	}

	signature, err := base64.StdEncoding.DecodeString(l.ProductKey.Signature)
	if err != nil {
		return err
	}

	message := []byte(macAddress +
		l.ProductKey.Node.LicenseID +
		l.ProductKey.Node.LicenseName +
		l.ProductKey.Node.CreateDate)
	hashed := sha256.Sum256(message)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return err
	}

	return nil
}
