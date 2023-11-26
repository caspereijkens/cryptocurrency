package util

import (
	"crypto/sha256"
	"math/big"
	"testing"
)

// TestHash256ToBigInt tests the hash256ToBigInt function
func TestHash256ToBigInt(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"Empty String", "", sha256SumToBigIntString("")},
		{"Normal String", "my secret", sha256SumToBigIntString("my secret")},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Hash256ToBigInt(tt.input)
			if result.String() != tt.expect {
				t.Errorf("hash256ToBigInt(%s) = %v, want %v", tt.input, result, tt.expect)
			}
		})
	}
}

// sha256SumToBigIntString computes the SHA-256 hash of a string and returns it as a big integer in string format.
// This is used to generate expected values for test cases.
func sha256SumToBigIntString(data string) string {
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
	return bigInt.String()
}
