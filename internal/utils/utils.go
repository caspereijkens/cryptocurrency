package utils

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
	"strings"

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

	return EncodeBase58(dataWithChecksum)
}

func DecodeBase58(s string) ([]byte, error) {
	num := new(big.Int)

	for _, c := range s {
		num.Mul(num, big.NewInt(58))
		num.Add(num, big.NewInt(int64(strings.IndexByte(base58Alphabet, byte(c)))))
	}

	combined := make([]byte, 25)
	copy(combined[25-len(num.Bytes()):], num.Bytes())

	checksum := combined[21:]
	if !bytes.Equal(Hash256(combined[:21])[:4], checksum) {
		return nil, fmt.Errorf("bad address: %x %x", checksum, Hash256(combined[:21])[:4])
	}

	return combined[1:21], nil
}

// TODO Make unit test
func Hash256(data []byte) []byte {
	sha256Digest := Sha256Hash(data)
	return Sha256Hash(sha256Digest)
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
	sha256Digest := Sha256Hash(s)
	ripemd160Digest := Ripemd160Hash(sha256Digest)
	return ripemd160Digest
}

func Sha1Hash(s []byte) []byte {
	sha1Hash := sha1.New()
	sha1Hash.Write(s)
	return sha1Hash.Sum(nil)
}

func Sha256Hash(s []byte) []byte {
	sha256Hash := sha256.New()
	sha256Hash.Write(s)
	return sha256Hash.Sum(nil)
}

func Ripemd160Hash(s []byte) []byte {
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(s)
	return ripemd160Hash.Sum(nil)
}

func FormatWithUnderscore(n int) string {
	str := strconv.Itoa(n)
	result := ""
	for i := 0; i < len(str); i++ {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += "_"
		}
		result += string(str[i])
	}
	return result
}

func EncodeVarint(i uint64) ([]byte, error) {
	if i < 0xfd {
		return []byte{byte(i)}, nil
	} else if i < 0x10000 {
		result := make([]byte, 3)
		result[0] = 0xfd
		binary.LittleEndian.PutUint16(result[1:], uint16(i))
		return result, nil
	} else if i < 0x100000000 {
		result := make([]byte, 5)
		result[0] = 0xfe
		binary.LittleEndian.PutUint32(result[1:], uint32(i))
		return result, nil
	} else if i < 0xffffffffffffffff {
		result := make([]byte, 9)
		result[0] = 0xff
		binary.LittleEndian.PutUint64(result[1:], i)
		return result, nil
	} else {
		return nil, fmt.Errorf("integer too large: %d", i)
	}
}

func ReadVarint(reader *bufio.Reader) (uint64, error) {
	buf := make([]byte, 1)
	_, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}

	i := uint64(buf[0])
	switch i {
	case 0xfd:
		// 0xfd means the next two bytes are the number
		return readLittleEndianUint16(reader)
	case 0xfe:
		// 0xfe means the next four bytes are the number
		return readLittleEndianUint32(reader)
	case 0xff:
		// 0xff means the next eight bytes are the number
		return readLittleEndianUint64(reader)
	default:
		// anything else is just the integer
		return i, nil
	}
}

// readLittleEndianUint16 reads a little-endian uint16 from the reader
func readLittleEndianUint16(reader *bufio.Reader) (uint64, error) {
	buf := make([]byte, 2)
	_, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}
	return uint64(binary.LittleEndian.Uint16(buf)), nil
}

// readLittleEndianUint32 reads a little-endian uint32 from the reader
func readLittleEndianUint32(reader *bufio.Reader) (uint64, error) {
	buf := make([]byte, 4)
	_, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}
	return uint64(binary.LittleEndian.Uint32(buf)), nil
}

// readLittleEndianUint64 reads a little-endian uint64 from the reader
func readLittleEndianUint64(reader *bufio.Reader) (uint64, error) {
	buf := make([]byte, 8)
	_, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}
