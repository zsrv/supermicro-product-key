package nonjson

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
)

// ParseEncodedProductKey returns a ProductKey, decoded from a base64-encoded string.
// The BMC MAC address associated with the product key is used in the decryption process.
// An error is returned instead if one occurs.
func ParseEncodedProductKey(encodedProductKey string, macAddress netinternal.HardwareAddr) (*ProductKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedProductKey)
	if err != nil {
		return nil, err
	}

	decrypted, err := decryptProductKey(decoded, macAddress)
	if err != nil {
		return nil, err
	}

	productKey, err := readDecryptedProductKey(decrypted)
	if err != nil {
		return nil, err
	}

	return productKey, nil
}

// decryptProductKey returns a decrypted product key, or an error if one occurs.
// The BMC MAC address is used to derive the decryption key and iv.
func decryptProductKey(encrypted []byte, macAddress netinternal.HardwareAddr) ([]byte, error) {
	keyingMaterial := macAddress.String() + "ej" + "mb"
	key, iv := generateKeyAndIV(keyingMaterial)

	plaintext, err := aesDecrypt(encrypted[16:], key, iv)
	if err != nil {
		return nil, err
	}

	plaintextWithHeader := append(encrypted[:16], plaintext...)
	return plaintextWithHeader, nil
}

// readDecryptedProductKey returns a ProductKey constructed from the decrypted
// product key bytes, or an error if one occurs.
func readDecryptedProductKey(decryptedProductKey []byte) (*ProductKey, error) {
	if len(decryptedProductKey) == 0 {
		msg := "product key byte array is null or empty"
		return nil, errors.New(msg)
	}
	if len(decryptedProductKey) > 255 {
		msg := fmt.Sprintf("max length of product key byte array is 255, actual length is %d", len(decryptedProductKey))
		return nil, errors.New(msg)
	}

	productKey := &ProductKey{}
	reader := bufio.NewReader(bytes.NewReader(decryptedProductKey))

	// Unencrypted block

	formatVersionUnencrypted, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if formatVersionUnencrypted != 0 {
		msg := fmt.Sprintf("format version must be 0x00, actual version is 0x%02X", productKey.FormatVersion)
		return nil, errors.New(msg)
	}
	productKey.FormatVersion = formatVersionUnencrypted

	softwareIdentifierUnencryptedByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	for i := 0; i < 14; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if b != 0 {
			return nil, errors.New("invalid data in unencrypted block")
		}
	}

	// Encrypted block

	formatVersion, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if formatVersion != formatVersionUnencrypted {
		msg := fmt.Sprintf("format versions do not match. version in unencrypted block is 0x%02X, version in encrypted block is 0x%02X",
			formatVersionUnencrypted, formatVersion)
		return nil, errors.New(msg)
	}

	softwareIdentifierByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if softwareIdentifierByte != softwareIdentifierUnencryptedByte {
		msg := fmt.Sprintf("software identifier byte in unencrypted block is 0x%02X, identifier in encrypted block is 0x%02X",
			softwareIdentifierUnencryptedByte, softwareIdentifierByte)
		return nil, errors.New(msg)
	}

	softwareIdentifierString, err := reader.ReadBytes(0)
	if err != nil {
		return nil, err
	}
	softwareIdentifierString = softwareIdentifierString[:len(softwareIdentifierString)-1]

	sid, err := SoftwareIdentifiers.ByID(softwareIdentifierByte)
	if err != nil {
		return nil, err
	}

	if string(softwareIdentifierString) != sid.DisplayName {
		msg := fmt.Sprintf("software identifier string '%s' does not match string '%s' that is associated with software identifier byte 0x%02X in the registry",
			softwareIdentifierString, sid.DisplayName, softwareIdentifierByte)
		return nil, errors.New(msg)
	}
	productKey.SoftwareIdentifier = *sid

	softwareVersion, err := reader.ReadBytes(0)
	if err != nil {
		return nil, err
	}
	productKey.SoftwareVersion = string(softwareVersion[:len(softwareVersion)-1])

	invoiceNumber, err := reader.ReadBytes(0)
	if err != nil {
		return nil, err
	}
	productKey.InvoiceNumber = string(invoiceNumber[:len(invoiceNumber)-1])

	creationDate := make([]byte, 4)
	if _, err := io.ReadFull(reader, creationDate); err != nil {
		return nil, err
	}
	productKey.CreationDate = bytesToDate(creationDate)

	expirationDate := make([]byte, 4)
	if _, err := io.ReadFull(reader, expirationDate); err != nil {
		return nil, err
	}
	productKey.ExpirationDate = bytesToDate(expirationDate)

	var keyProperty []byte
	keyPropertyLength, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if keyPropertyLength > 0 {
		keyProperty = make([]byte, keyPropertyLength)
		if _, err := io.ReadFull(reader, keyProperty); err != nil {
			return nil, err
		}
		productKey.Property = keyProperty
	}

	var secretData []byte
	secretDataLength, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if secretDataLength > 0 {
		secretData = make([]byte, secretDataLength)
		if _, err := io.ReadFull(reader, secretData); err != nil {
			return nil, err
		}
		productKey.SecretData = secretData
	}

	checksum, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	err = productKey.calculateChecksum()
	if err != nil {
		return nil, err
	}

	if checksum != productKey.Checksum {
		msg := fmt.Sprintf("checksum stored in product key is 0x%02X, calculated checksum is 0x%02X", checksum, productKey.Checksum)
		return nil, errors.New(msg)
	}

	remainingBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	for _, b := range remainingBytes {
		if b != 0 {
			msg := fmt.Sprintf("padding bytes must be 0x00, encountered 0x%02X", b)
			return nil, errors.New(msg)
		}
	}

	return productKey, nil
}
