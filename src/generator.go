package main

import (
	"crypto/rand"
	"math/big"
)

// Generate bounded value between minimum and maximum
func randIntBound(minV, maxV int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(maxV)-int64(minV)+1))
	res := int64(minV) + n.Int64()
	return int(res)
}

// Generate value between 0 and max including boundaries
func randInt(maxV int) int {
	return randIntBound(0, maxV)
}
