package nonjson

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/zsrv/supermicro-product-key/pkg/oob"
	"strconv"
	"time"
)

// NewProductKey returns a new ProductKey with default values set.
func NewProductKey() *ProductKey {
	return &ProductKey{
		FormatVersion:      0,
		SoftwareIdentifier: *SoftwareIdentifiers.Reserved,
		SoftwareVersion:    "none",
		InvoiceNumber:      "none",
		CreationDate:       time.Now().UTC(),
		ExpirationDate:     time.Unix(0, 0).UTC(),
		Property:           nil,
		SecretData:         nil,
		Checksum:           0,
	}
}

type ProductKey struct {
	// FormatVersion is the format version of the product key.
	// The only valid value is 0.
	FormatVersion byte
	// SoftwareIdentifier is the license type that the product
	// key is valid for.
	SoftwareIdentifier SoftwareIdentifier
	// SoftwareVersion is the software version the product key
	// is valid for.
	SoftwareVersion string
	// InvoiceNumber is the invoice number associated with the
	// product key.
	InvoiceNumber string
	// CreationDate is the date and time the product key was
	// generated, in UTC.
	CreationDate time.Time
	// ExpirationDate is the date and time the product key will
	// expire, in UTC. A value of 1970-01-01T00:00:00Z
	// (Unix time 0), indicates that the product key does not expire.
	ExpirationDate time.Time
	// The purpose of Property is not currently known. The value has
	// not yet been found to have been set in any product keys.
	Property []byte
	// SecretData is a value generated using the MAC address of the BMC
	// that the license is to be activated on as input.
	SecretData []byte
	// Checksum is a checksum value calculated using the other values
	// of the product key as input.
	Checksum byte
}

// Encode returns the encrypted, base64-encoded ProductKey associated with
// the given BMC MAC address.
func (pk *ProductKey) Encode(macAddress string) (string, error) {
	macAddress, err := oob.NormalizeMACAddress(macAddress)
	if err != nil {
		return "", err
	}

	buffer := bytes.NewBuffer(make([]byte, 0, 255))
	writer := bufio.NewWriterSize(buffer, 255)

	encryptedBuffer := bytes.NewBuffer(make([]byte, 0, 239))
	encryptedWriter := bufio.NewWriterSize(encryptedBuffer, 239)

	// Unencrypted block

	err = writer.WriteByte(pk.FormatVersion)
	if err != nil {
		return "", err
	}

	err = writer.WriteByte(pk.SoftwareIdentifier.ID)
	if err != nil {
		return "", err
	}

	_, err = writer.Write(bytes.Repeat([]byte{0}, 14))
	if err != nil {
		return "", err
	}

	// Encrypted block

	err = encryptedWriter.WriteByte(pk.FormatVersion)
	if err != nil {
		return "", err
	}

	err = encryptedWriter.WriteByte(pk.SoftwareIdentifier.ID)
	if err != nil {
		return "", err
	}

	_, err = encryptedWriter.WriteString(pk.SoftwareIdentifier.DisplayName)
	if err != nil {
		return "", err
	}
	err = encryptedWriter.WriteByte(0)
	if err != nil {
		return "", err
	}

	_, err = encryptedWriter.WriteString(pk.SoftwareVersion)
	if err != nil {
		return "", err
	}
	err = encryptedWriter.WriteByte(0)
	if err != nil {
		return "", err
	}

	_, err = encryptedWriter.WriteString(pk.InvoiceNumber)
	if err != nil {
		return "", err
	}
	err = encryptedWriter.WriteByte(0)
	if err != nil {
		return "", err
	}

	_, err = encryptedWriter.Write(dateToBytes(pk.CreationDate))
	if err != nil {
		return "", err
	}

	_, err = encryptedWriter.Write(dateToBytes(pk.ExpirationDate))
	if err != nil {
		return "", err
	}

	err = encryptedWriter.WriteByte(byte(len(pk.Property)))
	if err != nil {
		return "", err
	}
	if len(pk.Property) > 0 {
		_, err = encryptedWriter.Write(pk.Property)
		if err != nil {
			return "", err
		}
	}

	secretData, err := getEncryptedSecretData(macAddress, pk.getSecretDataKeyingMaterial())
	if err != nil {
		return "", err
	}
	err = encryptedWriter.WriteByte(byte(len(secretData)))
	if err != nil {
		return "", err
	}
	_, err = encryptedWriter.Write(secretData)
	if err != nil {
		return "", err
	}
	pk.SecretData = secretData

	err = pk.calculateChecksum()
	if err != nil {
		return "", err
	}
	err = encryptedWriter.WriteByte(pk.Checksum)
	if err != nil {
		return "", err
	}

	_, err = encryptedWriter.Write(bytes.Repeat([]byte{0}, encryptedWriter.Size()-encryptedWriter.Buffered()))
	if err != nil {
		return "", err
	}

	err = encryptedWriter.Flush()
	if err != nil {
		return "", err
	}

	_, err = writer.Write(encryptedBuffer.Bytes())
	if err != nil {
		return "", err
	}

	err = writer.Flush()
	if err != nil {
		return "", err
	}

	encrypted, err := encryptProductKey(buffer.Bytes(), macAddress)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(encrypted)

	return encoded, nil
}

// getSecretDataKeyingMaterial returns data that is used to derive an
// encryption key and iv for the secret data portion of the product key.
func (pk *ProductKey) getSecretDataKeyingMaterial() string {
	softwareIdentifier := fmt.Sprintf("%02X", pk.SoftwareIdentifier.ID)
	keyingMaterial := softwareIdentifier + pk.SoftwareVersion + strconv.FormatInt(pk.ExpirationDate.Unix(), 10)
	if pk.Property != nil {
		for _, b := range pk.Property {
			keyingMaterial += fmt.Sprintf("%02X", b)
		}
	}

	return keyingMaterial
}

// calculateChecksum calculates and stores the checksum of the product key.
func (pk *ProductKey) calculateChecksum() error {
	buf := bytes.NewBuffer(make([]byte, 0, 239))
	w := bufio.NewWriter(buf)

	// Unencrypted block

	err := w.WriteByte(pk.FormatVersion)
	if err != nil {
		return err
	}

	err = w.WriteByte(pk.SoftwareIdentifier.ID)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes.Repeat([]byte{0}, 14))
	if err != nil {
		return err
	}

	// Encrypted block

	err = w.WriteByte(pk.FormatVersion)
	if err != nil {
		return err
	}

	err = w.WriteByte(pk.SoftwareIdentifier.ID)
	if err != nil {
		return err
	}

	_, err = w.WriteString(pk.SoftwareIdentifier.DisplayName)
	if err != nil {
		return err
	}

	err = w.WriteByte(0)
	if err != nil {
		return err
	}

	_, err = w.WriteString(pk.SoftwareVersion)
	if err != nil {
		return err
	}

	err = w.WriteByte(0)
	if err != nil {
		return err
	}

	_, err = w.WriteString(pk.InvoiceNumber)
	if err != nil {
		return err
	}

	err = w.WriteByte(0)
	if err != nil {
		return err
	}

	_, err = w.Write(dateToBytes(pk.CreationDate))
	if err != nil {
		return err
	}

	_, err = w.Write(dateToBytes(pk.ExpirationDate))
	if err != nil {
		return err
	}

	err = w.WriteByte(byte(len(pk.Property)))
	if err != nil {
		return err
	}
	if len(pk.Property) > 0 {
		_, err = w.Write(pk.Property)
		if err != nil {
			return err
		}
	}

	err = w.WriteByte(byte(len(pk.SecretData)))
	if err != nil {
		return err
	}
	if len(pk.SecretData) > 0 {
		_, err = w.Write(pk.SecretData)
		if err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	var checksum byte = 0
	for _, b := range buf.Bytes() {
		checksum += b
	}
	pk.Checksum = checksum

	return nil
}
