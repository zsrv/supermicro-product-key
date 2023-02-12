package nonjson

import (
	"time"
)

// bytesToDate returns a time.Time from four bytes that represent
// the number of seconds since January 1, 1970, UTC.
// Dates are stored as four bytes in non-JSON product keys.
func bytesToDate(in []byte) time.Time {
	unixSeconds := bytesToUnixSeconds(in)
	return time.Unix(unixSeconds, 0).UTC()
}

// dateToBytes returns four bytes that represent the number of seconds
// since January 1, 1970, UTC.
// Dates are stored as four bytes in non-JSON product keys.
func dateToBytes(in time.Time) []byte {
	unixSeconds := in.Unix()
	return unixSecondsToBytes(unixSeconds)
}

// bytesToUnixSeconds returns the number of seconds since January 1, 1970, UTC
// extracted from four bytes.
// Dates are stored as four bytes in non-JSON product keys.
func bytesToUnixSeconds(in []byte) int64 {
	b3 := int64(in[3] & 255)
	b2 := int64(in[2] & 255)
	b1 := int64(in[1] & 255)
	b0 := int64(in[0] & 255)

	return (b3 << 24) + (b2 << 16) + (b1 << 8) + b0
}

// unixSecondsToBytes returns four bytes that store an int64 representing
// the number of seconds since January 1, 1970, UTC.
// Dates are stored as four bytes in non-JSON product keys.
func unixSecondsToBytes(in int64) []byte {
	out := make([]byte, 4)
	out[0] = byte(in)
	out[1] = byte(in >> 8)
	out[2] = byte(in >> 16)
	out[3] = byte(in >> 24)

	return out
}
