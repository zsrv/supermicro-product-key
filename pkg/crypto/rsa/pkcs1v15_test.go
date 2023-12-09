// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsa_test

import (
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	rsa2 "github.com/zsrv/supermicro-product-key/pkg/crypto/rsa"
	"testing"
)

type signPKCS1v15Test struct {
	in, out string
}

// These vectors have been tested with
//
//	`openssl rsautl -verify -inkey pk -in signature | hexdump -C`
var signPKCS1v15Tests = []signPKCS1v15Test{
	{"Test.\n", "a4f3fa6ea93bcdd0c57be020c1193ecbfd6f200a3d95c409769b029578fa0e336ad9a347600e40d3ae823b8c7e6bad88cc07c1d54c3a1523cbbb6d58efc362ae"},
}

func TestVerifyPKCS1v15(t *testing.T) {
	for i, test := range signPKCS1v15Tests {
		h := sha1.New()
		h.Write([]byte(test.in))
		digest := h.Sum(nil)

		sig, _ := hex.DecodeString(test.out)

		err := rsa2.VerifyPKCS1v15(rsaPublicKey, crypto.SHA1, digest, sig)
		if err != nil {
			t.Errorf("#%d %s", i, err)
		}
	}
}

var rsaPublicKey = parsePublicKey(`-----BEGIN RSA PUBLIC KEY-----
MEgCQQCymQ9JxH36jNQArmpNG4o7ahNkKyPyiwA7+5d5Ct6aTMgriyqBdH3ewIti
luU6CMMxaH7yXEv0k2uhwOYEHp0VAgMBAAE=
-----END RSA PUBLIC KEY-----`)

func parsePublicKey(s string) *rsa2.PublicKey {
	p, _ := pem.Decode([]byte(s))
	k, err := x509.ParsePKCS1PublicKey(p.Bytes)
	if err != nil {
		panic(err)
	}
	return &rsa2.PublicKey{
		N: k.N,
		E: k.E,
	}
}
