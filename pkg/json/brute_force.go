package json

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zsrv/supermicro-product-key/pkg/oob"
	"strings"
)

// BruteForceMACAddressFromString finds the MAC address associated with the productKey.
func BruteForceMACAddressFromString(productKey string) (string, error) {
	license, err := unmarshalLicenseFromJSON(productKey)
	if err != nil {
		return "", err
	}

	return BruteForceMACAddress(license)
}

// BruteForceMACAddress finds the MAC address associated with the license.
func BruteForceMACAddress(license License) (string, error) {
	brute := func(macBlock string, result chan string, done chan bool) {
		log.Debug().Msgf("searching mac address block '%s'", macBlock)

		for one := 0; one <= 255; one++ {
			hexOne := hex.EncodeToString([]byte{byte(one)})
			for two := 0; two <= 255; two++ {
				hexTwo := hex.EncodeToString([]byte{byte(two)})
				for three := 0; three <= 255; three++ {
					hexThree := hex.EncodeToString([]byte{byte(three)})

					mac := strings.ToUpper(macBlock + hexOne + hexTwo + hexThree)
					if err := license.VerifySignature(mac); err != nil {
						continue
					}

					result <- mac
					done <- true
				}
			}
		}

		log.Debug().Msgf("finished searching mac address block '%s' with no matches", macBlock)
		done <- true
	}

	result := make(chan string)
	done := make(chan bool)

	for _, macBlock := range oob.SupermicroMACAddressBlocks {
		go brute(macBlock, result, done)
	}

	for range oob.SupermicroMACAddressBlocks {
		select {
		case resultMAC := <-result:
			return strings.ToLower(resultMAC), nil
		case <-done:
			continue
		}
	}

	return "", errors.New("could not find a matching mac address")
}
