package gogrs

import (
	"math/rand"
	"time"
)

// RandInt return random int.
func RandInt() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}
