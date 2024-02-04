package block

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"testing"

	"github.com/caspereijkens/cryptocurrency/internal/utils"
)

func TestParse(t *testing.T) {
	blockRaw, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Fatalf("Failed to decode blockRaw hex: %v", err)
	}
	stream := bytes.NewReader(blockRaw)
	block, err := Parse(stream)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Test block fields
	wantVersion := uint32(0x20000002)
	if block.Version != wantVersion {
		t.Errorf("Version mismatch. Got: %x, Want: %x", block.Version, wantVersion)
	}

	wantPrevBlock, _ := hex.DecodeString("000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e")
	if !reflect.DeepEqual(block.PrevBlock[:], wantPrevBlock) {
		t.Errorf("PrevBlock mismatch. Got: %x, Want: %x", block.PrevBlock, wantPrevBlock)
	}

	wantMerkleRoot, _ := hex.DecodeString("be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b")
	if !reflect.DeepEqual(block.MerkleRoot[:], wantMerkleRoot) {
		t.Errorf("MerkleRoot mismatch. Got: %x, Want: %x", block.MerkleRoot, wantMerkleRoot)
	}

	wantTimestamp := uint32(0x59a7771e)
	if block.Timestamp != wantTimestamp {
		t.Errorf("Timestamp mismatch. Got: %x, Want: %x", block.Timestamp, wantTimestamp)
	}

	wantBits := uint32(0xe93c0118)
	if block.Bits != wantBits {
		t.Errorf("Bits mismatch. Got: %x, Want: %x", block.Bits, wantBits)
	}

	wantNonce := uint32(0xa4ffd71d)
	if block.Nonce != wantNonce {
		t.Errorf("Nonce mismatch. Got: %x, Want: %x", block.Nonce, wantNonce)
	}
}

func TestSerialize(t *testing.T) {
	// Create a byte slice representing the raw block data
	blockRaw, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Fatalf("Failed to decode blockRaw hex: %v", err)
	}

	// Create a bytes.Reader from the blockRaw data
	stream := bytes.NewReader(blockRaw)

	// Parse the block from the stream
	block, err := Parse(stream)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Serialize the block
	serializedBlock, err := block.Serialize()
	if err != nil {
		t.Fatalf("Serialization error: %v", err)
	}

	// Compare the serialized block with the original raw block data
	if !bytes.Equal(serializedBlock, blockRaw) {
		t.Errorf("Serialized block does not match original raw block data")
	}
}

func TestHash(t *testing.T) {
	// Create a byte slice representing the raw block data
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")

	// Create a bytes.Reader from the blockRaw data
	stream := bytes.NewReader(blockRaw)

	// Parse the block from the stream
	block, err := Parse(stream)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Calculate the hash of the block
	hash, err := block.Hash()

	if err != nil {
		t.Fatalf("Hash256 error: %v", err)
	}

	// Convert the expected hash from hexadecimal to byte slice
	expectedHash, _ := hex.DecodeString("0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523")

	// Compare the calculated hash with the expected hash
	if !bytes.Equal(hash, expectedHash) {
		t.Errorf("Hash mismatch. Got: %x, Want: %x", hash, expectedHash)
	}
}

func TestBIP9(t *testing.T) {
	// Block 1 (BIP9 signaled)
	block1Raw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block1Stream := bytes.NewReader(block1Raw)
	block1, _ := Parse(block1Stream)

	// Block 2 (BIP9 not signaled)
	block2Raw, _ := hex.DecodeString("0400000039fa821848781f027a2e6dfabbf6bda920d9ae61b63400030000000000000000ecae536a304042e3154be0e3e9a8220e5568c3433a9ab49ac4cbb74f8df8e8b0cc2acf569fb9061806652c27")
	block2Stream := bytes.NewReader(block2Raw)
	block2, _ := Parse(block2Stream)

	// Test BIP9 signaled block
	if !block1.BIP9() {
		t.Errorf("Block 1 should signal BIP9 but it didn't")
	}

	// Test BIP9 not signaled block
	if block2.BIP9() {
		t.Errorf("Block 2 should not signal BIP9 but it did")
	}
}

func TestBIP91(t *testing.T) {
	// Block 1 (BIP91 signaled)
	block1Raw, _ := hex.DecodeString("1200002028856ec5bca29cf76980d368b0a163a0bb81fc192951270100000000000000003288f32a2831833c31a25401c52093eb545d28157e200a64b21b3ae8f21c507401877b5935470118144dbfd1")
	block1Stream := bytes.NewReader(block1Raw)
	block1, _ := Parse(block1Stream)

	// Block 2 (BIP91 not signaled)
	block2Raw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block2Stream := bytes.NewReader(block2Raw)
	block2, _ := Parse(block2Stream)

	// Test BIP91 signaled block
	if !block1.BIP91() {
		t.Errorf("Block 1 should signal BIP91 but it didn't")
	}

	// Test BIP91 not signaled block
	if block2.BIP91() {
		t.Errorf("Block 2 should not signal BIP91 but it did")
	}
}

func TestBIP141(t *testing.T) {
	// Block 1 (BIP141 signaled)
	block1Raw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block1Stream := bytes.NewReader(block1Raw)
	block1, _ := Parse(block1Stream)

	// Block 2 (BIP141 not signaled)
	block2Raw, _ := hex.DecodeString("0000002066f09203c1cf5ef1531f24ed21b1915ae9abeb691f0d2e0100000000000000003de0976428ce56125351bae62c5b8b8c79d8297c702ea05d60feabb4ed188b59c36fa759e93c0118b74b2618")
	block2Stream := bytes.NewReader(block2Raw)
	block2, _ := Parse(block2Stream)

	// Test BIP141 signaled block
	if !block1.BIP141() {
		t.Errorf("Block 1 should signal BIP141 but it didn't")
	}

	// Test BIP141 not signaled block
	if block2.BIP141() {
		t.Errorf("Block 2 should not signal BIP141 but it did")
	}
}

func TestProofOfWork(t *testing.T) {
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	blockHash := utils.ReverseBytes(utils.Hash256(blockRaw))
	blockId := hex.EncodeToString(blockHash)
	if blockId != "0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523" {
		t.Errorf("This block doesn't satisfy the expected target")
	}
}

func TestTarget(t *testing.T) {
	// Parse the block
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	blockStream := bytes.NewReader(blockRaw)
	block, _ := Parse(blockStream)
	expectedTarget := "13ce9000000000000000000000000000000000000000000"
	calculatedTarget := block.Target()
	gotTarget := fmt.Sprintf("%x", calculatedTarget)
	if gotTarget != expectedTarget {
		t.Errorf("Expected target: %s, got: %s", expectedTarget, gotTarget)
	}
}

func TestDifficulty(t *testing.T) {
	blockRaw, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Fatalf("Error decoding blockRaw: %v", err)
	}

	blockStream := bytes.NewReader(blockRaw)
	block, _ := Parse(blockStream)

	difficulty := block.Difficulty()

	expectedDifficulty := big.NewInt(888171856257)

	if difficulty.Cmp(expectedDifficulty) != 0 {
		t.Errorf("Difficulty calculation incorrect, got: %v, want: %v", difficulty, expectedDifficulty)
	}
}

func TestCheckPow(t *testing.T) {
	blockRaw1, _ := hex.DecodeString("04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec1")
	stream1 := bytes.NewReader(blockRaw1)
	block1, err := Parse(stream1)
	if err != nil {
		t.Error(err)
	}
	if !block1.CheckPOW() {
		t.Errorf("Block 1's Proof of Work check failed")
	}

	blockRaw2, _ := hex.DecodeString("04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec0")
	stream2 := bytes.NewReader(blockRaw2)
	block2, err := Parse(stream2)
	if err != nil {
		t.Error(err)
	}
	if block2.CheckPOW() {
		t.Errorf("Block 2's Proof of Work check passed, but it shouldn't have")
	}
}

func TestTargetToBits(t *testing.T) {
	// Parse the block
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	blockStream := bytes.NewReader(blockRaw)
	block, _ := Parse(blockStream)
	calculatedTarget := block.Target()
	gotBits := TargetToBits(calculatedTarget)
	if gotBits != block.Bits {
		t.Errorf("Expected target: %d, got: %d", gotBits, block.Bits)
	}

	block471744Raw, _ := hex.DecodeString("00000020fdf740b0e49cf75bb3d5168fb3586f7613dcc5cd89675b0100000000000000002e37b144c0baced07eb7e7b64da916cd3121f2427005551aeb0ec6a6402ac7d7f0e4235954d801187f5da9f5")
	block471744Stream := bytes.NewReader(block471744Raw)
	block471744, _ := Parse(block471744Stream)
	bits471744 := TargetToBits(block471744.Target())
	fmt.Println(bits471744)

	block473759Raw, _ := hex.DecodeString("000000201ecd89664fd205a37566e694269ed76e425803003628ab010000000000000000bfcade29d080d9aae8fd461254b041805ae442749f2a40100440fc0e3d5868e55019345954d80118a1721b2e")
	block473759Stream := bytes.NewReader(block473759Raw)
	block473759, _ := Parse(block473759Stream)
	bits473759 := TargetToBits(block473759.Target())

	if bits471744 != bits473759 {
		t.Errorf("Bits should be same")
	}
}

func TestCalculateNewBits(t *testing.T) {
	prevBits, _ := strconv.ParseUint("54d80118", 16, 32)
	timeDifferential := int64(302400)
	wantBits, _ := strconv.ParseUint("00157617", 16, 32)

	gotBits := CalculateNewBits(uint32(prevBits), timeDifferential)

	if uint32(wantBits) != gotBits {
		t.Errorf("calculateNewBits() = %d, want %d", gotBits, wantBits)
	}
}
