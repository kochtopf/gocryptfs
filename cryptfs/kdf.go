package cryptfs

import (
	"fmt"
	"golang.org/x/crypto/scrypt"
	"math"
	"os"
)

const (
	// 1 << 16 uses 64MB of memory,
	// takes 4 seconds on my Atom Z3735F netbook
	SCRYPT_DEFAULT_LOGN = 16
)

type scryptKdf struct {
	Salt   []byte
	N      int
	R      int
	P      int
	KeyLen int
}

func NewScryptKdf(logN int) scryptKdf {
	var s scryptKdf
	s.Salt = RandBytes(KEY_LEN)
	if logN <= 0 {
		s.N = 1 << SCRYPT_DEFAULT_LOGN
	} else {
		if logN < 10 {
			fmt.Printf("Error: scryptn below 10 is too low to make sense. Aborting.\n")
			os.Exit(1)
		}
		s.N = 1 << uint32(logN)
	}
	s.R = 8 // Always 8
	s.P = 1 // Always 1
	s.KeyLen = KEY_LEN
	return s
}

func (s *scryptKdf) DeriveKey(pw string) []byte {
	k, err := scrypt.Key([]byte(pw), s.Salt, s.N, s.R, s.P, s.KeyLen)
	if err != nil {
		panic(fmt.Sprintf("DeriveKey failed: %s", err.Error()))
	}
	return k
}

// LogN - N is saved as 2^LogN, but LogN is much easier to work with.
// This function gives you LogN = Log2(N).
func (s *scryptKdf) LogN() int {
	return int(math.Log2(float64(s.N)) + 0.5)
}
