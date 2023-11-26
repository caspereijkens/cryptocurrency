package util

import (
	"crypto/sha256"
	"math/big"
)

func Hash256ToBigInt(data string) *big.Int {
	// First SHA-256 hash
	hasher1 := sha256.New()
	hasher1.Write([]byte(data))
	firstHashBytes := hasher1.Sum(nil)

	// Second SHA-256 hash on the result of the first hash
	hasher2 := sha256.New()
	hasher2.Write(firstHashBytes)
	secondHashBytes := hasher2.Sum(nil)

	// Convert the second hash bytes to a big.Int
	bigInt := new(big.Int)
	bigInt.SetBytes(secondHashBytes)
	return bigInt
}
