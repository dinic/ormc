package utils

import (
	"bytes"
	"strings"
)

const (
	HungarianPkgSeq = byte('_')
	HungarianDirSeq = byte('-')
)

// Hungarian2Camel ASCII-only strings
func Hungarian2Camel(w string) string {
	if len(w) == 0 {
		return w
	}

	s := []byte{HungarianDirSeq, HungarianPkgSeq}
	var b strings.Builder
	b.Grow(len(w))

	st := 0
	if c := w[st]; c >= 'a' && c <= 'z' {
		c -= 'a' - 'A'
		st = 1
		b.WriteByte(c)
	}

	found := false
	for i := st; i < len(w); i++ {
		c := w[i]
		if bytes.ContainsRune(s, rune(c)) {
			found = true
			continue
		}

		if found {
			found = false
			if c >= 'a' && c <= 'z' {
				c -= 'a' - 'A'
			}
		}

		b.WriteByte(c)
	}

	return b.String()
}

// Camel2Hungarian ASCII-only strings
func Camel2Hungarian(w string, seq ...byte) string {
	s := HungarianPkgSeq // default seq _
	if len(seq) != 0 {
		s = seq[0]
	}

	var b strings.Builder
	pre := false
	cur := false
	for i := 0; i < len(w); i++ {
		pre = cur
		cur = false
		c := w[i]
		if c >= 'A' && c <= 'Z' {
			cur = true
			c += 'a' - 'A'
			if i != 0 && pre == false {
				b.WriteByte(s)
			}
		}

		b.WriteByte(c)
	}

	return b.String()
}
