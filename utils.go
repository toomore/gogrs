package gogrs

import (
	"math/rand"
	"time"
)

func RandInt() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}
