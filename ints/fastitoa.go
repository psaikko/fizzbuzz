// This implementation adapted from Go source strconv/itoa.go:

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ints

const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

// CopyItoa converts a positive integer to base10 string
//
// adapted from Go source function formatBits
//
func CopyItoa(dst []byte, i int, u uint64) {

	for u >= 100 {
		is := u % 100 * 2
		u /= 100
		i -= 2
		dst[i+1] = smallsString[is+1]
		dst[i+0] = smallsString[is+0]
	}

	// u < 100
	is := u * 2
	i--
	dst[i] = smallsString[is+1]
	if u >= 10 {
		i--
		dst[i] = smallsString[is]
	}
}
