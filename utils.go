package gogrs

import (
	"crypto/rand"
	"math/big"
)

// RandInt return random int.
func RandInt() int64 {
	result, _ := rand.Int(rand.Reader, big.NewInt(102400))
	return result.Int64()
}
