package nonjson

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/zsrv/supermicro-product-key/pkg/oob"
	"strings"
)

// encryptProductKey returns an encrypted product key, or an error if one occurs.
// The BMC MAC address is used to derive the encryption key and iv.
func encryptProductKey(decrypted []byte, macAddress string) ([]byte, error) {
	macAddress, err := oob.NormalizeMACAddress(macAddress)
	if err != nil {
		return nil, err
	}

	keyingMaterial := strings.ToUpper(macAddress) + "ej" + "mb"
	key, iv := generateKeyAndIV(keyingMaterial)

	ciphertext, err := aesEncrypt(decrypted[16:], key, iv)
	if err != nil {
		return nil, err
	}

	ciphertextWithHeader := append(decrypted[:16], ciphertext...)
	return ciphertextWithHeader, nil
}

// getEncryptedSecretData returns the secret data value associated with the given
// BMC MAC address and additional keying material derived from other content in
// the product key.
func getEncryptedSecretData(macAddress string, keyingMaterialInput string) ([]byte, error) {
	macAddress, err := oob.NormalizeMACAddress(macAddress)
	if err != nil {
		return nil, err
	}

	macAddress = strings.ToUpper(macAddress)

	keyingMaterial := macAddress + "am" + "ac" + keyingMaterialInput
	key, iv := generateKeyAndIV(keyingMaterial)

	encryptedMAC, err := aesEncrypt([]byte(macAddress), key, iv)
	if err != nil {
		return nil, err
	}
	secretData := fmt.Sprintf("%02x", md5.Sum(encryptedMAC))

	return []byte(secretData), nil
}

// generateKeyAndIV returns a key and iv derived from the input.
func generateKeyAndIV(input string) ([]byte, []byte) {
	digest := fmt.Sprintf("%02x", md5.Sum([]byte(input)))
	iv := digest[:16]
	key := digest[16:]

	return []byte(key), []byte(iv)
}

// aesEncrypt returns the plaintext encrypted using AES in CBC mode, or an error if one occurs.
// PKCS#7 padding will be added to the plaintext before encryption.
func aesEncrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	plaintext, err := pkcs7Pad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// aesDecrypt returns the ciphertext decrypted using AES in CBC mode, or an error if one occurs.
// PKCS#7 padding will be removed from the plaintext after decryption.
func aesDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// pkcs7Pad appends PKCS#7 padding.
// https://github.com/Luzifer/go-openssl
func pkcs7Pad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, fmt.Errorf("invalid block size %d", blockSize)
	}

	padLen := 1
	for ((len(data) + padLen) % blockSize) != 0 {
		padLen++
	}

	pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, pad...), nil
}

// pkcs7Unpad returns a slice of the original data with PKCS#7 padding removed.
// https://github.com/Luzifer/go-openssl
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, fmt.Errorf("invalid block size %d", blockSize)
	}

	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data length %d", len(data))
	}

	padlen := int(data[len(data)-1])
	if padlen > blockSize || padlen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padlen], nil
}
