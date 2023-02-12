package nonjson

import (
	"encoding/hex"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/zsrv/supermicro-product-key/pkg/oob"
)

// BruteForceProductKeyMACAddress finds the MAC address associated with an encrypted
// product key. The MAC address can then be used to decrypt the key.
func BruteForceProductKeyMACAddress(encodedProductKey string) (string, error) {
	brute := func(macBlock string, result chan string, done chan bool) {
		log.Debug().Msgf("searching mac address block '%s'", macBlock)

		for one := 0; one <= 255; one++ {
			hexOne := hex.EncodeToString([]byte{byte(one)})
			for two := 0; two <= 255; two++ {
				hexTwo := hex.EncodeToString([]byte{byte(two)})
				for three := 0; three <= 255; three++ {
					hexThree := hex.EncodeToString([]byte{byte(three)})

					mac := macBlock + hexOne + hexTwo + hexThree

					if _, err := DecodeProductKey(encodedProductKey, mac); err != nil {
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
			return resultMAC, nil
		case <-done:
			continue
		}
	}

	return "", errors.New("could not find a matching mac address")
}
