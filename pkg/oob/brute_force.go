package oob

import (
	"encoding/hex"
	"errors"
	"github.com/rs/zerolog/log"
)

// SupermicroMACAddressBlocks contains MAC address blocks
// that have been assigned to Supermicro.
// https://gitlab.com/wireshark/wireshark/-/blob/master/manuf
var SupermicroMACAddressBlocks = []string{
	"002590",
	"003048",
	"0cc47a",
	"3cecef",
	"7cc255",
	"ac1f6b",
}

// BruteForceProductKeyMACAddress returns the MAC address associated with the
// product key, or an error if one occurs or the MAC address was not found.
func BruteForceProductKeyMACAddress(productKey string) (string, error) {
	productKey, err := NormalizeProductKey(productKey)
	if err != nil {
		return "", err
	}

	brute := func(macBlock string, result chan string, done chan bool) {
		log.Debug().Msgf("searching mac address block '%s'", macBlock)

		for one := 0; one <= 255; one++ {
			hexOne := hex.EncodeToString([]byte{byte(one)})
			for two := 0; two <= 255; two++ {
				hexTwo := hex.EncodeToString([]byte{byte(two)})
				for three := 0; three <= 255; three++ {
					hexThree := hex.EncodeToString([]byte{byte(three)})

					mac := macBlock + hexOne + hexTwo + hexThree

					generatedKey, _ := EncodeOOBProductKey(mac)
					if generatedKey != productKey {
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

	for _, macBlock := range SupermicroMACAddressBlocks {
		go brute(macBlock, result, done)
	}

	for range SupermicroMACAddressBlocks {
		select {
		case resultMAC := <-result:
			return resultMAC, nil
		case <-done:
			continue
		}
	}

	return "", errors.New("could not find a matching mac address")
}
