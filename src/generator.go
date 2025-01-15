package main

import (
    "crypto/rand"
    "math/big"
)

// Generate bounded value between minimum and maximum
func randIntBound(min, max int) int {
    n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)-int64(min)+1))
    res := int64(min) + n.Int64()
    return int(res)
}

// Generate value between 0 and max including boundaries
func randInt(max int) int {
    return randIntBound(0, max)
}
