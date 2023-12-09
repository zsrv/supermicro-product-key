package json

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	rsainternal "github.com/zsrv/supermicro-product-key/pkg/crypto/rsa"
)

// parseRSAPublicKey returns a rsa.PublicKey parsed from pemData.
func parseRSAPublicKey(pemData []byte) (*rsainternal.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubInternal := &rsainternal.PublicKey{
		N: pub.N,
		E: pub.E,
	}

	return pubInternal, nil
}
