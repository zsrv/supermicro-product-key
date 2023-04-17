package json

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/pkg/errors"
)

func VerifyProductKeySignature(productKey string, macAddress string) error {
	license, err := unmarshalLicenseFromJSON(productKey)
	if err != nil {
		return err
	}

	return license.VerifySignature(macAddress)
}

func unmarshalLicenseFromJSON(productKey string) (License, error) {
	var license License
	err := json.Unmarshal([]byte(productKey), &license)
	if err != nil {
		return License{}, errors.WithMessage(err, "could not unmarshal license JSON")
	}

	return license, nil
}

// parseRSAPublicKey returns a rsa.PublicKey parsed from pemData.
func parseRSAPublicKey(pemData []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub, nil
}
