package helpers

import "math/rand"

func DistinctStrings(elems ...string) []string {
	m := map[string]bool{}
	for _, e := range elems {
		m[e] = true
	}
	var retElems []string
	for e, _ := range m {
		retElems = append(retElems, e)
	}
	return retElems
}

var (
	CharsetCharUpperCase = []rune("ABCDEFGHIJKLMNOPARSTUVWXYZ")
	CharsetCharLowerCase = []rune("abcdefghijklmnopqrstuvwxyz")
	CharsetChars         = append(CharsetCharUpperCase, CharsetCharLowerCase...)
	CharsetDigit         = []rune("0123456789")
	CharsetCharDigit     = append(CharsetChars, CharsetDigit...)
	CharsetHex           = []rune("0123456789abcdef")
)

func RandomString(size int, charset []rune) string {
	szChar := len(charset)
	if szChar == 0 {
		return ""
	}
	var buff []rune
	for i := 0; i < size; i++ {
		idx := rand.Intn(szChar)
		buff = append(buff, charset[idx])
	}
	return string(buff)
}
