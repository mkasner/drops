package drops

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

const (
	minRandLength = 32 // min length for random string
	maxRandLength = 64 // max length for random string
)

// generates random string specified with minLength and maxLength
func Random() string {
	return randString(randLength(minRandLength, maxRandLength))
}

// gets random string in range
func RandomInRange(min, max int) string {
	return randString(randLength(min, max))
}

func randLength(min, max int) int {
	rl, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(rl.Int64()) + min
}

// returns random hex string specified by length
func randString(length int) string {
	buff := make([]byte, length/2)
	io.ReadFull(rand.Reader, buff)
	return fmt.Sprintf("%x", buff)
}
