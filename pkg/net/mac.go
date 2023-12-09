// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"errors"
	"unsafe"
)

const hexDigit = "0123456789ABCDEF"

// A HardwareAddr represents a physical hardware address.
type HardwareAddr []byte

func (a HardwareAddr) String() string {
	if len(a) == 0 {
		return ""
	}
	buf := make([]byte, 0, len(a)*3-1)
	for _, b := range a {
		buf = append(buf, hexDigit[b>>4])
		buf = append(buf, hexDigit[b&0xF])
	}
	return *(*string)(unsafe.Pointer(&buf))
}

// ParseMAC parses s as an IEEE 802 MAC-48 or EUI-48 link-layer address
// using one of the following formats:
//
//	00005e005301
//	00:00:5e:00:53:01
//	00-00-5e-00-53-01
//	0000.5e00.5301
func ParseMAC(s string) (hw HardwareAddr, err error) {
	if len(s) < 12 {
		goto error
	}

	if s[2] == ':' || s[2] == '-' {
		if (len(s)+1)%3 != 0 {
			goto error
		}
		n := (len(s) + 1) / 3
		if n != 6 {
			goto error
		}
		hw = make(HardwareAddr, n)
		for x, i := 0, 0; i < n; i++ {
			var ok bool
			if hw[i], ok = Xtoi2(s[x:], s[2]); !ok {
				goto error
			}
			x += 3
		}
	} else if s[4] == '.' {
		if (len(s)+1)%5 != 0 {
			goto error
		}
		n := 2 * (len(s) + 1) / 5
		if n != 6 {
			goto error
		}
		hw = make(HardwareAddr, n)
		for x, i := 0, 0; i < n; i += 2 {
			var ok bool
			if hw[i], ok = Xtoi2(s[x:x+2], 0); !ok {
				goto error
			}
			if hw[i+1], ok = Xtoi2(s[x+2:], s[4]); !ok {
				goto error
			}
			x += 5
		}
	} else if len(s) == 12 {
		hw = make(HardwareAddr, 6)
		for x, i := 0, 0; i < 6; i++ {
			var ok bool
			if hw[i], ok = Xtoi2(s[x:x+2], 0); !ok {
				goto error
			}
			x += 2
		}
	} else {
		goto error
	}
	return hw, nil

error:
	return nil, errors.New("invalid MAC address: " + s)
}
