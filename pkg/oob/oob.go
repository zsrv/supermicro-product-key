package oob

import (
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"hash"
	"slices"
	"sync"
	"unsafe"

	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
)

// oobHMACKey is the HMAC key used to generate the SFT-OOB-LIC product key.
var oobHMACKey = []byte{0x85, 0x44, 0xe3, 0xb4, 0x7e, 0xca, 0x58, 0xf9, 0x58, 0x30, 0x43, 0xf8}

// hsdcHMACKey is an HMAC key used for unknown purposes.
var hsdcHMACKey = []byte{0x39, 0xcb, 0x2a, 0x1a, 0x3d, 0x74, 0x8f, 0xf1, 0xde, 0xe4, 0x6b, 0x87}

const hexDigit = "0123456789ABCDEF"

type ProductKey []byte

func (pk ProductKey) String() string {
	if len(pk) == 0 {
		return ""
	}
	buf := make([]byte, 0, 17)
	for i, b := range pk {
		if i > 0 && i%2 == 0 {
			buf = append(buf, '-')
		}
		buf = append(buf, hexDigit[b>>4])
		buf = append(buf, hexDigit[b&0xF])
	}
	return *(*string)(unsafe.Pointer(&buf))
}

func (pk ProductKey) Equal(x ProductKey) bool {
	if len(pk) == len(x) {
		return string(pk) == string(x)
	}
	return false
}

func ParseProductKey(s string) (pk ProductKey, err error) {
	if len(s) < 24 {
		goto error
	}

	if s[4] == '-' {
		if (len(s)+1)%5 != 0 {
			goto error
		}
		n := 2 * (len(s) + 1) / 5
		if n != 12 {
			goto error
		}
		pk = make(ProductKey, n)
		for x, i := 0, 0; i < n; i += 2 {
			var ok bool
			if pk[i], ok = netinternal.Xtoi2(s[x:x+2], 0); !ok {
				goto error
			}
			if pk[i+1], ok = netinternal.Xtoi2(s[x+2:], s[4]); !ok {
				goto error
			}
			x += 5
		}
	} else if len(s) == 24 {
		pk = make(ProductKey, 12)
		for x, i := 0, 0; i < 12; i++ {
			var ok bool
			if pk[i], ok = netinternal.Xtoi2(s[x:x+2], 0); !ok {
				goto error
			}
			x += 2
		}
	} else {
		goto error
	}
	return pk, nil

error:
	return nil, errors.New("invalid product key: " + s)
}

type hmacData struct {
	key  []byte
	hash hash.Hash
}

var hashPool = sync.Pool{
	New: func() interface{} {
		return &hmacData{
			key:  nil,
			hash: nil,
		}
	},
}

// EncodeOOBProductKey returns an OOB product key for the given BMC MAC address.
func EncodeOOBProductKey(macAddress netinternal.HardwareAddr) (ProductKey, error) {
	h := hashPool.Get().(*hmacData)
	defer hashPool.Put(h)
	if !slices.Equal(h.key, oobHMACKey) {
		h.hash = hmac.New(sha1.New, oobHMACKey)
		h.key = oobHMACKey
	}
	h.hash.Reset()

	h.hash.Write(macAddress)
	sum := h.hash.Sum(nil)

	productKey := ProductKey(sum[:12])
	return productKey, nil
}
