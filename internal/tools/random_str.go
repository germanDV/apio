package tools

import "golang.org/x/exp/rand"

// RandomStr produces a random string containing ascii printable characters of the given length.
func RandomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		// 33-94 are ASCII printable characters.
		b[i] = byte(rand.Intn(62) + 33)
	}
	return string(b)
}
