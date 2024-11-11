package main

import (
	"testing"
	"unicode/utf8"
)

// go test -run=FuzzReverse
// go test -fuzz=Fuzz
// go test -fuzz=Fuzz -fuzztime 10s
// go test -run={FuzzTestName}/{filename}
func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev, err1 := ReverseString(orig)
		if err1 != nil {
			return
		}
		doubleRev, err2 := ReverseString(rev)
		if err2 != nil {
			t.Skip() // = return
		}
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
