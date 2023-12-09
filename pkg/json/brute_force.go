package json

import (
	"errors"
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"
	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
	"github.com/zsrv/supermicro-product-key/pkg/oob"
)

// BruteForceMACAddressFromString finds the MAC address associated with the productKey.
func BruteForceMACAddressFromString(productKey string) (string, error) {
	license, err := ParseLicense(productKey)
	if err != nil {
		return "", err
	}

	return BruteForceMACAddress(license)
}

// BruteForceMACAddress finds the MAC address associated with the license.
func BruteForceMACAddress(license License) (string, error) {
	mac := make(netinternal.HardwareAddr, 6)
	gen := func(macBlocks [][3]byte, numGoroutines int) <-chan netinternal.HardwareAddr {
		out := make(chan netinternal.HardwareAddr, numGoroutines*100)
		go func() {
			for _, macBlock := range macBlocks {
				log.Debug().Msgf("started generating macs for block %X", macBlock)
				for one := 0; one <= 255; one++ {
					for two := 0; two <= 255; two++ {
						for three := 0; three <= 255; three++ {
							mac[0] = macBlock[0]
							mac[1] = macBlock[1]
							mac[2] = macBlock[2]
							mac[3] = byte(one)
							mac[4] = byte(two)
							mac[5] = byte(three)

							m := make(netinternal.HardwareAddr, len(mac))
							copy(m, mac)
							out <- m
						}
					}
				}
				log.Debug().Msgf("finished generating macs for block %X", macBlock)
			}
			close(out)
		}()
		return out
	}

	publicKey, err := parseRSAPublicKey(LicenseSigningPublicKey)
	if err != nil {
		return "", err
	}

	proc := func(in <-chan netinternal.HardwareAddr) <-chan string {
		out := make(chan string, 1)
		go func() {
			for macAddr := range in {
				if err := license.VerifyWithKey(publicKey, macAddr); err != nil {
					continue
				}
				out <- macAddr.String()
			}
			close(out)
		}()
		return out
	}

	merge := func(cs ...<-chan string) <-chan string {
		var wg sync.WaitGroup
		out := make(chan string, 1)

		// Start an output goroutine for each input channel in cs. output
		// copies values from c to out until c is closed, then calls wg.Done.
		output := func(c <-chan string) {
			for s := range c {
				out <- s
			}
			wg.Done()
		}
		wg.Add(len(cs))
		for _, c := range cs {
			go output(c)
		}

		// Start a goroutine to close out once all the output goroutines are
		// done. This must start after the wg.Add call.
		go func() {
			wg.Wait()
			close(out)
		}()
		return out
	}

	numCPU := runtime.NumCPU()
	in := gen(oob.SupermicroMACAddressBlocks, numCPU)

	// Fan out
	var procChannels []<-chan string
	for i := 0; i < numCPU-1; i++ {
		procChannels = append(procChannels, proc(in))
	}

	// Consume the merged output from all goroutines
	for result := range merge(procChannels...) {
		return result, nil
	}

	return "", errors.New("could not find a matching mac address")
}
