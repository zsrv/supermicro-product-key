// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsa

import (
	"errors"
	"github.com/zsrv/supermicro-product-key/pkg/crypto/internal/bigmod"
	"math/big"
	"sync"
)

// A PublicKey represents the public part of an RSA key.
type PublicKey struct {
	once    sync.Once
	N       *big.Int // modulus
	E       int      // public exponent
	modulus *bigmod.Modulus
}

// Size returns the modulus size in bytes. Raw signatures and ciphertexts
// for or by this public key will have the same size.
func (pub *PublicKey) Size() int {
	return (pub.N.BitLen() + 7) / 8
}

func encrypt(pub *PublicKey, plaintext *[]byte) ([]byte, error) {

	// Most of the CPU time for encryption and verification is spent in this
	// NewModulusFromBig call, because PublicKey doesn't have a Precomputed
	// field. If performance becomes an issue, consider placing a private
	// sync.Once on PublicKey to compute this.

	var err error
	pub.once.Do(func() {
		pub.modulus, err = bigmod.NewModulusFromBig(pub.N)
	})
	if err != nil {
		return nil, err
	}
	N := pub.modulus

	m, err := bigmod.NewNat().SetBytes(*plaintext, N)
	if err != nil {
		return nil, err
	}
	e := uint(pub.E)

	return bigmod.NewNat().ExpShort(m, e, N).Bytes(N), nil
}

// ErrVerification represents a failure to verify a signature.
// It is deliberately vague to avoid adaptive attacks.
var ErrVerification = errors.New("crypto/rsa: verification error")
