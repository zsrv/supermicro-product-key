package oob

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// oobHMACKey is the HMAC key used to generate the SFT-OOB-LIC product key.
var oobHMACKey = []byte{0x85, 0x44, 0xe3, 0xb4, 0x7e, 0xca, 0x58, 0xf9, 0x58, 0x30, 0x43, 0xf8}

// hsdcHMACKey is an HMAC key used for unknown purposes.
var hsdcHMACKey = []byte{0x39, 0xcb, 0x2a, 0x1a, 0x3d, 0x74, 0x8f, 0xf1, 0xde, 0xe4, 0x6b, 0x87}

// EncodeOOBProductKey returns an OOB product key for the given BMC MAC address.
func EncodeOOBProductKey(macAddress string) (string, error) {
	macAddress, err := NormalizeMACAddress(macAddress)
	if err != nil {
		return "", err
	}

	macAddressBytes, err := hex.DecodeString(macAddress)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha1.New, oobHMACKey)
	mac.Write(macAddressBytes)
	sum := mac.Sum(nil)

	productKey := fmt.Sprintf("%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X",
		sum[0], sum[1], sum[2], sum[3], sum[4], sum[5],
		sum[6], sum[7], sum[8], sum[9], sum[10], sum[11],
	)

	return productKey, nil
}

// NormalizeProductKey returns the legacy OOB product key in normalized form.
func NormalizeProductKey(key string) (string, error) {
	key = strings.Map(
		func(r rune) rune {
			if (r < '0' || r > '9') && (r < 'A' || r > 'F') && (r < 'a' || r > 'f') {
				return -1
			}
			return r
		},
		strings.ToUpper(key),
	)

	if len(key) != 24 {
		return "", errors.New("product key without separators must have a length of 24")
	}

	var keySplit []string
	for i := 0; i < len(key); i += 4 {
		keySplit = append(keySplit, key[i:i+4])
	}

	return strings.Join(keySplit, "-"), nil
}

// NormalizeMACAddress returns the MAC address in lowercase with separators removed,
// or an error if the MAC address is not valid.
func NormalizeMACAddress(macAddress string) (string, error) {
	macAddress = strings.Map(
		func(r rune) rune {
			if (r < '0' || r > '9') && (r < 'A' || r > 'F') && (r < 'a' || r > 'f') {
				return -1
			}
			return r
		},
		macAddress,
	)

	if len(macAddress) != 12 {
		return "", errors.New("mac address without separators must have a length of 12")
	}

	return strings.ToLower(macAddress), nil
}
