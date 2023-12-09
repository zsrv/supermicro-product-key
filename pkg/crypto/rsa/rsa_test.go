// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsa_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	rsainternal "github.com/zsrv/supermicro-product-key/pkg/crypto/rsa"
	"strings"
	"testing"
)

func testingKey(s string) string { return strings.ReplaceAll(s, "TESTING KEY", "PRIVATE KEY") }

func parseKey(s string) *rsa.PrivateKey {
	p, _ := pem.Decode([]byte(s))
	k, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		panic(err)
	}
	return k
}

var test2048Key = parseKey(testingKey(`-----BEGIN RSA TESTING KEY-----
MIIEnwIBAAKCAQBxY8hCshkKiXCUKydkrtQtQSRke28w4JotocDiVqou4k55DEDJ
akvWbXXDcakV4HA8R2tOGgbxvTjFo8EK470w9O9ipapPUSrRRaBsSOlkaaIs6OYh
4FLwZpqMNBVVEtguVUR/C34Y2pS9kRrHs6q+cGhDZolkWT7nGy5eSEvPDHg0EBq1
1hu6HmPmI3r0BInONqJg2rcK3U++wk1lnbD3ysCZsKOqRUms3n/IWKeTqXXmz2XK
J2t0NSXwiDmA9q0Gm+w0bXh3lzhtUP4MlzS+lnx9hK5bjzSbCUB5RXwMDG/uNMQq
C4MmA4BPceSfMyAIFjdRLGy/K7gbb2viOYRtAgEDAoIBAEuX2tchZgcGSw1yGkMf
OB4rbZhSSiCVvB5r1ew5xsnsNFCy1ducMo7zo9ehG2Pq9X2E8jQRWfZ+JdkX1gdC
fiCjSkHDxt+LceDZFZ2F8O2bwXNF7sFAN0rvEbLNY44MkB7jgv9c/rs8YykLZy/N
HH71mteZsO2Q1JoSHumFh99cwWHFhLxYh64qFeeH6Gqx6AM2YVBWHgs7OuKOvc8y
zUbf8xftPht1kMwwDR1XySiEYtBtn74JflK3DcT8oxOuCZBuX6sMJHKbVP41zDj+
FJZBmpAvNfCEYJUr1Hg+DpMLqLUg+D6v5vpliburbk9LxcKFZyyZ9QVe7GoqMLBu
eGsCgYEAummUj4MMKWJC2mv5rj/dt2pj2/B2HtP2RLypai4et1/Ru9nNk8cjMLzC
qXz6/RLuJ7/eD7asFS3y7EqxKxEmW0G8tTHjnzR/3wnpVipuWnwCDGU032HJVd13
LMe51GH97qLzuDZjMCz+VlbCNdSslMgWWK0XmRnN7Yqxvh6ao2kCgYEAm7fTRBhF
JtKcaJ7d8BQb9l8BNHfjayYOMq5CxoCyxa2pGBv/Mrnxv73Twp9Z/MP0ue5M5nZt
GMovpP5cGdJLQ2w5p4H3opcuWeYW9Yyru2EyCEAI/hD/Td3QVP0ukc19BDuPl5Wg
eIFs218uiVOU4pw3w+Et5B1PZ/F+ZLr5LGUCgYB8RmMKV11w7CyRnVEe1T56Ru09
Svlp4qQt0xucHr8k6ovSkTO32hd10yxw/fyot0lv1T61JHK4yUydhyDHYMQ81n3O
IUJqIv/qBpuOxvQ8UqwIQ3iU69uOk6TIhSaNlqlJwffQJEIgHf7kOdbOjchjMA7l
yLpmETPzscvUFGcXmwKBgGfP4i1lg283EvBp6Uq4EqQ/ViL6l5zECXce1y8Ady5z
xhASqiHRS9UpN9cU5qiCoyae3e75nhCGym3+6BE23Nede8UBT8G6HuaZZKOzHSeW
IVrVW1QLVN6T4DioybaI/gLSX7pjwFBWSJI/dFuNDexoJS1AyUK+NO/2VEMnUMhD
AoGAOsdn3Prnh/mjC95vraHCLap0bRBSexMdx77ImHgtFUUcSaT8DJHs+NZw1RdM
SZA0J+zVQ8q7B11jIgz5hMz+chedwoRjTL7a8VRTKHFmmBH0zlEuV7L79w6HkRCQ
VRg10GUN6heGLv0aOHbPdobcuVDH4sgOqpT1QnOuce34sQs=
-----END RSA TESTING KEY-----`))

var test2048PublicKey = parsePublicKey(`-----BEGIN RSA PUBLIC KEY-----
MIIBBwKCAQBxY8hCshkKiXCUKydkrtQtQSRke28w4JotocDiVqou4k55DEDJakvW
bXXDcakV4HA8R2tOGgbxvTjFo8EK470w9O9ipapPUSrRRaBsSOlkaaIs6OYh4FLw
ZpqMNBVVEtguVUR/C34Y2pS9kRrHs6q+cGhDZolkWT7nGy5eSEvPDHg0EBq11hu6
HmPmI3r0BInONqJg2rcK3U++wk1lnbD3ysCZsKOqRUms3n/IWKeTqXXmz2XKJ2t0
NSXwiDmA9q0Gm+w0bXh3lzhtUP4MlzS+lnx9hK5bjzSbCUB5RXwMDG/uNMQqC4Mm
A4BPceSfMyAIFjdRLGy/K7gbb2viOYRtAgED
-----END RSA PUBLIC KEY-----`)

func BenchmarkVerifyPKCS1v15(b *testing.B) {
	b.Run("2048", func(b *testing.B) {
		hashed := sha256.Sum256([]byte("testing"))
		s, err := rsa.SignPKCS1v15(rand.Reader, test2048Key, crypto.SHA256, hashed[:])
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := rsa.VerifyPKCS1v15(&test2048Key.PublicKey, crypto.SHA256, hashed[:], s)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkVerifyPKCS1v15Internal(b *testing.B) {
	b.Run("2048", func(b *testing.B) {
		hashed := sha256.Sum256([]byte("testing"))
		s, err := rsa.SignPKCS1v15(rand.Reader, test2048Key, crypto.SHA256, hashed[:])
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := rsainternal.VerifyPKCS1v15(test2048PublicKey, crypto.SHA256, hashed[:], s)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
