package utils

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"
)

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

func TestSerializeInt(t *testing.T) {
	tests := []struct {
		input    *big.Int
		expected []byte
	}{
		{big.NewInt(0), []byte{0x00}},                                         // Zero
		{big.NewInt(127), []byte{0x7F}},                                       // Positive number below 128 (0x80)
		{big.NewInt(128), []byte{0x00, 0x80}},                                 // Positive number above 127
		{new(big.Int).SetBytes([]byte{0x00, 0x81}), []byte{0x00, 0x81}},       // Bytes with high bit set
		{new(big.Int).SetBytes([]byte{0x00, 0x00, 0x00, 0x00}), []byte{0x00}}, // Null bytes
	}

	for _, test := range tests {
		result := SerializeInt(test.input)

		if !bytes.Equal(result, test.expected) {
			t.Errorf("SerializeInt(%v) returned %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestLstripNullBytes(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte{0x00, 0x00, 0x00, 0x01, 0x02}, []byte{0x01, 0x02}},
		{[]byte{0x00, 0x00, 0x00, 0x00}, []byte{}},
		{[]byte{0x01, 0x02, 0x03}, []byte{0x01, 0x02, 0x03}},
		{[]byte{}, []byte{}},
	}

	for _, testCase := range testCases {
		result := lstripNullBytes(testCase.input)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("lstripNullBytes(%v) = %v, expected %v", testCase.input, result, testCase.expected)
		}
	}
}

func TestEncodeBase58(t *testing.T) {
	testsBytes := []struct {
		input    []byte
		expected string
	}{
		{[]byte{0x00, 0x00, 0x00, 0x00}, "1111"},
		{[]byte{0x00, 0x00, 0x00, 0x01}, "1112"},
		{[]byte{0x00, 0x00, 0x00, 0x42}, "11129"},
		{[]byte{0x12, 0x34, 0x56, 0x78, 0x9a}, "348ALpH"},
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "11111111"},
		{[]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}, "C3CPq7c8PY"},
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, "11111112"},
	}

	for _, test := range testsBytes {
		result := EncodeBase58(test.input)
		if result != test.expected {
			t.Errorf("For input %v, expected %s, but got %s", test.input, test.expected, result)
		}
	}

	testsHex := []struct {
		input    string
		expected string
	}{
		{"7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", "9MA8fRQrT4u8Zj8ZRd6MAiiyaxb2Y1CMpvVkHQu5hVM6"},
		{"eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c", "4fE3H2E6XMp4SsxtwinF7w9a34ooUrwWe4WsW1458Pd"},
		{"c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6", "EQJsjkd6JaGwxrjEhfeqPenqHwrBmPQZjJGNSCHBkcF7"},
	}

	for _, test := range testsHex {
		sBytes, _ := hex.DecodeString(test.input)
		result := EncodeBase58(sBytes)
		if result != test.expected {
			t.Errorf("For input %v, expected %s, but got %s", test.input, test.expected, result)
		}
	}
}

func TestHash160(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello world",
			expected: "d7d5ee7824ff93f94c3055af9382c86c68b5ca92",
		},
		{
			input:    "Hi mom!",
			expected: "eab3813216e715e5830980f3532d44a50df3ce11",
		},
	}

	for _, test := range tests {
		inputBytes := []byte(test.input)
		result := Hash160(inputBytes)
		resultHex := hex.EncodeToString(result)

		if resultHex != test.expected {
			t.Errorf("For input '%s', expected %s but got %s", test.input, test.expected, resultHex)
		}
	}
}
