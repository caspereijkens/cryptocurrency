package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func EncodeBase58(s []byte) string {
	var mod *big.Int
	var result []byte

	count := 0
	for _, c := range s {
		if c != 0 {
			break
		}
		count++
	}

	num := new(big.Int).SetBytes(s)
	prefix := bytes.Repeat([]byte{'1'}, count)
	for num.Cmp(big.NewInt(0)) > 0 {
		num, mod = new(big.Int).DivMod(num, big.NewInt(58), new(big.Int))
		result = append([]byte{base58Alphabet[mod.Int64()]}, result...)
	}

	return string(append(prefix, result...))
}

// TODO Make unit test
func EncodeBase58Checksum(data []byte) string {
	// Calculate the hash256
	hash256 := Hash256(data)

	dataWithChecksum := append(data, hash256[:4]...)

	base58Encoded := EncodeBase58(dataWithChecksum)

	return base58Encoded
}

// TODO Make unit test
func Hash256(data []byte) []byte {
	// First SHA-256 hash
	hasher1 := sha256.New()
	hasher1.Write(data)
	firstHashBytes := hasher1.Sum(nil)

	// Second SHA-256 hash on the result of the first hash
	hasher2 := sha256.New()
	hasher2.Write(firstHashBytes)
	secondHashBytes := hasher2.Sum(nil)
	return secondHashBytes
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

// sha256 followed by ripemd160
func Hash160(s []byte) []byte {
	sha256hash := sha256.New()
	sha256hash.Write(s)
	sha256Digest := sha256hash.Sum(nil)

	ripemd160hash := ripemd160.New()
	ripemd160hash.Write(sha256Digest)
	ripemd160Digest := ripemd160hash.Sum(nil)

	return ripemd160Digest
}