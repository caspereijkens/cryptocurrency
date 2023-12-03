package util

import (
	"crypto/sha256"
	"encoding/hex"
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

// TestHmacSHA256 tests the hmacSHA256 function
func TestHmacSHA256(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		key         []byte
		data        []byte
		expectedHex string // Expected output in hexadecimal format
	}{
		{
			name:        "Empty key and data",
			key:         []byte(""),
			data:        []byte(""),
			expectedHex: "b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad",
		},
		{
			name:        "Example key and data",
			key:         []byte("key"),
			data:        []byte("The quick brown fox jumps over the lazy dog"),
			expectedHex: "f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd8",
		},
		// Add more test cases as necessary
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := HmacSHA256(tc.key, tc.data)
			actualHex := hex.EncodeToString(actual)
			if actualHex != tc.expectedHex {
				t.Errorf("hmacSHA256(%x, %s) = %s, want %s", tc.key, tc.data, actualHex, tc.expectedHex)
			}
		})
	}
}
