package gogrs

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

// RandInt return random int.
func RandInt() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

func RandCryptInt() int64 {
	result, _ := crand.Int(crand.Reader, big.NewInt(102400))
	return result.Int64()
}
