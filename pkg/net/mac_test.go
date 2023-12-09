// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"reflect"
	"strings"
	"testing"
)

var parseMACTests = []struct {
	in  string
	out HardwareAddr
	err string
}{
	// See RFC 7042, Section 2.1.1.
	{"00:00:5e:00:53:01", HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}, ""},
	{"00-00-5e-00-53-01", HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}, ""},
	{"0000.5e00.5301", HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}, ""},

	{"ab:cd:ef:AB:CD:EF", HardwareAddr{0xab, 0xcd, 0xef, 0xab, 0xcd, 0xef}, ""},

	{"00005e005301", HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}, ""},
	{"00005E005301", HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}, ""},

	{"01.02.03.04.05.06", nil, "invalid MAC address"},
	{"01:02:03:04:05:06:", nil, "invalid MAC address"},
	{"x1:02:03:04:05:06", nil, "invalid MAC address"},
	{"01002:03:04:05:06", nil, "invalid MAC address"},
	{"01:02003:04:05:06", nil, "invalid MAC address"},
	{"01:02:03004:05:06", nil, "invalid MAC address"},
	{"01:02:03:04005:06", nil, "invalid MAC address"},
	{"01:02:03:04:05006", nil, "invalid MAC address"},
	{"01-02:03:04:05:06", nil, "invalid MAC address"},
	{"01:02-03-04-05-06", nil, "invalid MAC address"},
	{"0123:4567:89AF", nil, "invalid MAC address"},
	{"0123-4567-89AF", nil, "invalid MAC address"},
}

func TestParseMAC(t *testing.T) {
	match := func(err error, s string) bool {
		if s == "" {
			return err == nil
		}
		return err != nil && strings.Contains(err.Error(), s)
	}

	for i, tt := range parseMACTests {
		out, err := ParseMAC(tt.in)
		if !reflect.DeepEqual(out, tt.out) || !match(err, tt.err) {
			t.Errorf("ParseMAC(%q) = %v, %v, want %v, %v", tt.in, out, err, tt.out, tt.err)
		}
		if tt.err == "" {
			// Verify that serialization works too, and that it round-trips.
			s := out.String()
			out2, err := ParseMAC(s)
			if err != nil {
				t.Errorf("%d. ParseMAC(%q) = %v", i, s, err)
				continue
			}
			if !reflect.DeepEqual(out2, out) {
				t.Errorf("%d. ParseMAC(%q) = %v, want %v", i, s, out2, out)
			}
		}
	}
}
