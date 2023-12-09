package net

// Bigger than we need, not too big to worry about overflow
const big = 0xFFFFFF

// Hexadecimal to integer.
// Returns number, characters consumed, success.
func Xtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s); i++ {
		if '0' <= s[i] && s[i] <= '9' {
			n *= 16
			n += int(s[i] - '0')
		} else if 'a' <= s[i] && s[i] <= 'f' {
			n *= 16
			n += int(s[i]-'a') + 10
		} else if 'A' <= s[i] && s[i] <= 'F' {
			n *= 16
			n += int(s[i]-'A') + 10
		} else {
			break
		}
		if n >= big {
			return 0, i, false
		}
	}
	if i == 0 {
		return 0, i, false
	}
	return n, i, true
}

// Xtoi2 converts the next two hex digits of s into a byte.
// If s is longer than 2 bytes then the third byte must be e.
// If the first two bytes of s are not hex digits or the third byte
// does not match e, false is returned.
func Xtoi2(s string, e byte) (byte, bool) {
	if len(s) > 2 && s[2] != e {
		return 0, false
	}
	n, ei, ok := Xtoi(s[:2])
	return byte(n), ok && ei == 2
}
