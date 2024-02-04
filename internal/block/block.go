package block

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/big"
	"slices"

	"github.com/caspereijkens/cryptocurrency/internal/utils"
)

const (
	twoWeeks = int64(14 * 24 * 60 * 60)
)

// Block struct represents a Bitcoin block
type Block struct {
	Version    uint32
	PrevBlock  [32]byte
	MerkleRoot [32]byte
	Timestamp  uint32
	Bits       uint32
	Nonce      uint32
}

// Parse reads from the given byte stream and parses a block
func Parse(r io.Reader) (*Block, error) {
	block := &Block{}
	err := binary.Read(r, binary.LittleEndian, &block.Version)
	if err != nil {
		return nil, err
	}

	_, err = io.ReadFull(r, block.PrevBlock[:])
	if err != nil {
		return nil, err
	}
	slices.Reverse(block.PrevBlock[:])

	_, err = io.ReadFull(r, block.MerkleRoot[:])
	if err != nil {
		return nil, err
	}
	slices.Reverse(block.MerkleRoot[:])

	err = binary.Read(r, binary.LittleEndian, &block.Timestamp)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &block.Bits)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &block.Nonce)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// Serialize serializes the block into a byte slice
func (b *Block) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Serialize version
	err := binary.Write(&buf, binary.LittleEndian, b.Version)
	if err != nil {
		return nil, err
	}

	// Serialize prev_block (reversed)
	_, err = buf.Write(utils.ReverseBytes(b.PrevBlock[:]))
	if err != nil {
		return nil, err
	}

	// Serialize merkle_root (reversed)
	_, err = buf.Write(utils.ReverseBytes(b.MerkleRoot[:]))
	if err != nil {
		return nil, err
	}

	// Serialize timestamp
	err = binary.Write(&buf, binary.LittleEndian, b.Timestamp)
	if err != nil {
		return nil, err
	}

	// Serialize bits
	err = binary.Write(&buf, binary.BigEndian, b.Bits)
	if err != nil {
		return nil, err
	}

	// Serialize nonce
	err = binary.Write(&buf, binary.BigEndian, b.Nonce)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Hash returns the hash256 interpreted little endian of the block
func (b *Block) Hash() ([]byte, error) {
	data, err := b.Serialize()
	if err != nil {
		return nil, err
	}
	return utils.ReverseBytes(utils.Hash256(data)), nil
}

// BIP9 returns whether this block is signaling readiness for BIP9
func (b *Block) BIP9() bool {
	return b.Version>>29 == 0x001
}

// BIP91 returns whether this block is signaling readiness for BIP91
func (b *Block) BIP91() bool {
	return b.Version>>4&1 == 1
}

// BIP141 returns whether this block is signaling readiness for BIP141
func (b *Block) BIP141() bool {
	return b.Version>>1&1 == 1
}

func (b *Block) Target() *big.Int {
	return BitsToTarget(b.Bits)
}

// BitsToTarget converts the bits representation to target
func BitsToTarget(bits uint32) *big.Int {
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, bits)
	exponent := byteSlice[3]
	byteSlice[3] = 0x00
	coefficient := big.NewInt(int64(binary.LittleEndian.Uint32(byteSlice)))
	target := new(big.Int).Exp(big.NewInt(256), big.NewInt(int64(exponent-3)), nil)
	target.Mul(target, coefficient)
	return target
}

// Difficulty returns the block difficulty based on the bits
func (b *Block) Difficulty() *big.Int {
	lowestDifficultyBits := uint32(0xffff001d)
	lowestTarget := BitsToTarget(lowestDifficultyBits)

	currentTarget := b.Target()

	difficulty := new(big.Int).Quo(lowestTarget, currentTarget)

	return difficulty
}

// CheckPOW returns whether this block satisfies proof of work
func (b *Block) CheckPOW() bool {
	hash, _ := b.Hash()
	hashInt := new(big.Int).SetBytes(hash)
	target := b.Target()
	return target.Cmp(hashInt) == 1
}

func TargetToBits(target *big.Int) uint32 {
	var exponent int
	var coefficient []byte

	rawBytes := make([]byte, 32)
	target.FillBytes(rawBytes)
	rawBytes = utils.LstripNullBytes(rawBytes)
	if rawBytes[0] > 0x7f {
		exponent = len(rawBytes) + 1
		coefficient = append([]byte{0x00}, rawBytes[:2]...)
	} else {
		exponent = len(rawBytes)
		coefficient = rawBytes[:3]
	}
	bits := append(utils.ReverseBytes(coefficient), byte(exponent))

	return binary.BigEndian.Uint32(bits)
}

func CalculateNewBits(previousBits uint32, timeDifferential int64) uint32 {
	previousTarget := BitsToTarget(previousBits)

	// Ensure time differential is within the specified range between 0.5 week and 8 weeks
	if timeDifferential > twoWeeks*4 {
		timeDifferential = twoWeeks * 4
	}
	if timeDifferential < twoWeeks/4 {
		timeDifferential = twoWeeks / 4
	}

	newTarget := new(big.Int).Mul(previousTarget, big.NewInt(timeDifferential))
	newTarget.Div(newTarget, big.NewInt(twoWeeks))

	return TargetToBits(newTarget)
}
