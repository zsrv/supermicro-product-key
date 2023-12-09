package nonjson

import (
	"errors"

	"github.com/rs/zerolog/log"
	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
	"github.com/zsrv/supermicro-product-key/pkg/oob"
)

// BruteForceMACAddress finds the MAC address associated with an encrypted
// product key. The MAC address can then be used to decrypt the key.
func BruteForceMACAddress(encodedProductKey string) (string, error) {
	brute := func(macBlock [3]byte, result chan netinternal.HardwareAddr, done chan bool) {
		log.Debug().Msgf("searching mac address block %X", macBlock)

		mac := make(netinternal.HardwareAddr, 6)
		for one := 0; one <= 255; one++ {
			for two := 0; two <= 255; two++ {
				for three := 0; three <= 255; three++ {
					mac[0] = macBlock[0]
					mac[1] = macBlock[1]
					mac[2] = macBlock[2]
					mac[3] = byte(one)
					mac[4] = byte(two)
					mac[5] = byte(three)

					if _, err := ParseEncodedProductKey(encodedProductKey, mac); err != nil {
						continue
					}

					m := make(netinternal.HardwareAddr, len(mac))
					copy(m, mac)
					result <- m
					done <- true
				}
			}
		}

		log.Debug().Msgf("finished searching mac address block %X with no matches", macBlock)
		done <- true
	}

	result := make(chan netinternal.HardwareAddr)
	done := make(chan bool)

	for _, macBlock := range oob.SupermicroMACAddressBlocks {
		go brute(macBlock, result, done)
	}

	for range oob.SupermicroMACAddressBlocks {
		select {
		case resultMAC := <-result:
			return resultMAC.String(), nil
		case <-done:
			continue
		}
	}

	return "", errors.New("could not find a matching mac address")
}
