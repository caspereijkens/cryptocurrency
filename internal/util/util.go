package util

import (
	"crypto/hmac"
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

// HmacSHA256 computes the HMAC SHA-256 digest of the data using the given key
func HmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func SerializeInt(i *big.Int) []byte {
	bytes := i.FillBytes(make([]byte, 32))

	// Trim leading null bytes
	bytes = lstripNullBytes(bytes)

	// Add a null byte if the high bit is set
	if len(bytes) > 0 && bytes[0]&0x80 != 0 {
		bytes = append([]byte{0x00}, bytes...)
	}

	// Ensure a non-empty byte slice
	if len(bytes) == 0 {
		bytes = []byte{0x00}
	}

	return bytes
}

// lstripNullBytes trims leading null bytes from a byte slice
func lstripNullBytes(data []byte) []byte {
	var i int
	for i = 0; i < len(data); i++ {
		if data[i] != 0 {
			break
		}
	}
	return data[i:]
}
