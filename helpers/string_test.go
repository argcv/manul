package helpers

import (
	"math/rand"
	"testing"
	"time"
)

func TestRandomString(t *testing.T) {
	rand.Seed(time.Now().Unix())
	s := RandomString(10, CharsetDigit)
	t.Logf("Random Digits: %s", s)
	if len(s) != 10 {
		t.Fatalf("Incorrect Size:: %v", len(s))
	}
	s = RandomString(10, CharsetCharLowerCase)
	t.Logf("Random Lowercase: %s", s)
	if len(s) != 10 {
		t.Fatalf("Incorrect Size:: %v", len(s))
	}
	s = RandomString(10, CharsetHex)
	t.Logf("Random HexDigits: %s", s)
	if len(s) != 10 {
		t.Fatalf("Incorrect Size:: %v", len(s))
	}

}
