package signatureverification

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/caspereijkens/cryptocurrency/internal/ellipticcurve"
	"github.com/caspereijkens/cryptocurrency/internal/utils"
)

// TestSignatureVerification checks the signature verification process.
func TestS256PointVerify(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		pointX string
		pointY string
		z      string
		r      string
		s      string
	}{
		{
			pointX: "0x04519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574",
			pointY: "0x82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4",
			z:      "0xbc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423",
			r:      "0x37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6",
			s:      "0x8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec",
		},
		{
			pointX: "0x887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c",
			pointY: "0x61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34",
			z:      "0xec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60",
			r:      "0xac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395",
			s:      "0x68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4",
		},
		{
			pointX: "0x887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c",
			pointY: "0x61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34",
			z:      "0x7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d",
			r:      "0xeff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c",
			s:      "0xc7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		z, r, s, x, y, err := parseSignatureParts(map[string]string{
			"z": tc.z,
			"r": tc.r,
			"s": tc.s,
			"x": tc.pointX,
			"y": tc.pointY,
		})
		if err != nil {
			t.Fatalf("Failed to parse signature parts: %v", err)
		}

		point, err := createEllipticCurvePoint(x, y)
		if err != nil {
			t.Fatalf("Failed to create elliptic curve point: %v", err)
		}

		sig := NewSignature(r, s)

		if !point.Verify(z, sig) {
			t.Error("Could not verify signature.")
		}
	}
}

func TestOrderofGenerator(t *testing.T) {
	// setup
	identity, _ := ellipticcurve.NewPoint(nil, nil, &A.FieldElement, &B.FieldElement)
	// Call your function
	result, err := G.ScalarMultiplication(N)

	if err != nil {
		t.Errorf("Calculation did not go correct.")
	}
	if !result.Equal(identity) {
		t.Errorf("Point should be the identity point")
	}
}

func TestSignatureString(t *testing.T) {
	a := NewSignature(big.NewInt(7), big.NewInt(17))
	expected := "Signature(7,11)"
	if a.String() != expected {
		t.Errorf("Signature String representation '%s' is not as expected '%s': ", a.String(), expected)
	}
}

func TestSignatureSerialize(t *testing.T) {
	expectedHexString := "3045022037206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c60221008ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec"
	r := "0x37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6"
	rInt, _ := new(big.Int).SetString(r, 0)
	s := "0x8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec"
	sInt, _ := new(big.Int).SetString(s, 0)
	sig := NewSignature(rInt, sInt)
	der := sig.Serialize()
	derHexString := hex.EncodeToString(der)
	if derHexString != expectedHexString {
		t.Errorf("failed to create DER format:\nGot:\n%s\nExpected:\n%s", derHexString, expectedHexString)
	}
}

func TestSignAndVerify(t *testing.T) {
	privKey, err := NewPrivateKey(Hash256ToBigInt("my secret"))
	if err != nil {
		t.Errorf("failed to create private key: %v", err)
	}
	z := Hash256ToBigInt("my message")
	k := big.NewInt(1234567890)
	sig, err := privKey.Sign(z)
	if err != nil {
		t.Errorf("Signing failed with private '%08x', signature hash '%08x' and random number k '%s'", privKey, z, k.String())
	}
	P, _ := G.ScalarMultiplication(privKey.Secret)
	if !P.Verify(z, sig) {
		t.Errorf("Signature verification failed  private '%08x', signature hash '%08x' and random number k '%s'", privKey, z, k.String())
	}
}

func TestSerializeS256Point(t *testing.T) {
	// Define the test case struct
	type testCase struct {
		name        string
		secret      *big.Int
		compressed  bool
		expectedSec []byte
		secLength   int
	}

	// Create a slice of test cases
	testCases := []testCase{
		{
			name:        "Test Case 1",
			secret:      big.NewInt(5000),
			compressed:  false,
			expectedSec: []byte{4, 255, 229, 88, 227, 136, 133, 47, 1, 32, 228, 106, 242, 209, 179, 112, 248, 88, 84, 168, 235, 8, 65, 129, 30, 206, 14, 62, 3, 210, 130, 213, 124, 49, 93, 199, 40, 144, 164, 241, 10, 20, 129, 192, 49, 176, 59, 53, 27, 13, 199, 153, 1, 202, 24, 160, 12, 240, 9, 219, 219, 21, 122, 29, 16},
			secLength:   65,
		},
		{
			name:        "Test Case 2",
			secret:      new(big.Int).Exp(big.NewInt(2018), big.NewInt(5), nil),
			compressed:  false,
			expectedSec: []byte{4, 2, 127, 61, 161, 145, 132, 85, 224, 60, 70, 246, 89, 38, 106, 27, 181, 32, 78, 149, 157, 183, 54, 77, 47, 71, 59, 223, 143, 10, 19, 204, 157, 255, 135, 100, 127, 208, 35, 193, 59, 74, 73, 148, 241, 118, 145, 137, 88, 6, 225, 180, 11, 87, 244, 253, 34, 88, 26, 79, 70, 133, 31, 59, 6},
			secLength:   65,
		},
		{
			name:        "Test Case 3",
			secret:      big.NewInt(3917405024756549),
			compressed:  false,
			expectedSec: []byte{4, 217, 12, 214, 37, 238, 135, 221, 56, 101, 109, 217, 92, 247, 159, 101, 246, 15, 114, 115, 182, 125, 48, 150, 230, 139, 216, 30, 79, 83, 66, 105, 31, 132, 46, 250, 118, 47, 213, 153, 97, 208, 233, 152, 3, 198, 30, 219, 168, 179, 227, 247, 220, 58, 52, 24, 54, 249, 119, 51, 174, 191, 152, 113, 33},
			secLength:   65,
		},
		{
			name:        "Test Case 4",
			secret:      big.NewInt(5001),
			compressed:  true,
			expectedSec: []byte{3, 87, 164, 243, 104, 134, 138, 138, 109, 87, 41, 145, 228, 132, 230, 100, 129, 15, 241, 76, 5, 192, 250, 2, 50, 117, 37, 17, 81, 254, 14, 83, 209},
			secLength:   33,
		},
		{
			name:        "Test Case 5",
			secret:      new(big.Int).Exp(big.NewInt(2019), big.NewInt(5), nil),
			compressed:  true,
			expectedSec: []byte{2, 147, 62, 194, 210, 177, 17, 185, 39, 55, 236, 18, 241, 197, 210, 15, 50, 51, 160, 173, 33, 205, 139, 54, 208, 188, 167, 160, 207, 165, 203, 135, 1},
			secLength:   33,
		},
		{
			name:        "Test Case 6",
			secret:      big.NewInt(3917405025026849),
			compressed:  true,
			expectedSec: []byte{2, 150, 190, 91, 18, 146, 246, 200, 86, 179, 197, 101, 78, 136, 111, 193, 53, 17, 70, 32, 89, 8, 156, 223, 156, 71, 150, 35, 191, 203, 231, 118, 144},
			secLength:   33,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			privKey, err := NewPrivateKey(tc.secret)
			if err != nil {
				t.Fatalf("failed to create private key: %v", err)
			}

			sec := privKey.Point.Serialize(tc.compressed)

			if !bytes.Equal(sec, tc.expectedSec) {
				t.Errorf("Incorrect SEC: %v", sec)
			}

			if len(sec) != tc.secLength {
				t.Errorf("Incorrect SEC length: %d", len(sec))
			}
		})
	}
}

func TestS256PointAddress(t *testing.T) {
	privateKey1, _ := NewPrivateKey(big.NewInt(5002))
	privateKey2, _ := NewPrivateKey(new(big.Int).Exp(big.NewInt(2020), big.NewInt(5), nil))
	secret3, _ := new(big.Int).SetString("0x12345deadbeef", 0)
	privateKey3, _ := NewPrivateKey(secret3)

	// Test cases
	testCases := []struct {
		point       *S256Point
		compressed  bool
		testnet     bool
		expected    string
		description string
	}{
		{privateKey1.Point, false, true, "mmTPbXQFxboEtNRkwfh6K51jvdtHLxGeMA", "Uncompressed SEC on testnet"},
		{privateKey2.Point, true, true, "mopVkxp8UhXqRYbCYJsbeE1h1fiF64jcoH", "Compressed SEC on testnet"},
		{privateKey3.Point, true, false, "1F1Pn2y6pDb68E5nYJJeba4TLg2U7B6KF1", "Compressed SEC on mainnet"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := tc.point.Address(tc.compressed, tc.testnet)
			if result != tc.expected {
				t.Errorf("Address() returned %s, expected %s", result, tc.expected)
			}
		})
	}
}

func TestS256PointHash160(t *testing.T) {
	privateKey1, _ := NewPrivateKey(big.NewInt(5002))
	privateKey2, _ := NewPrivateKey(new(big.Int).Exp(big.NewInt(2020), big.NewInt(5), nil))
	secret3, _ := new(big.Int).SetString("0x12345deadbeef", 0)
	privateKey3, _ := NewPrivateKey(secret3)

	expected1, _ := hex.DecodeString("41243614aecd13819d7a7f348a4a07fbcb29d8e5")
	expected2, _ := hex.DecodeString("5b1257a7bb81398208766db21a4959ae068310ea")
	expected3, _ := hex.DecodeString("99a4c61750789253f69fd750ac0d021263373305")

	// Test cases
	testCases := []struct {
		point       *S256Point
		compressed  bool
		expected    []byte
		description string
	}{
		{privateKey1.Point, false, expected1, "Uncompressed hash160"},
		{privateKey2.Point, true, expected2, "Compressed hash160"},
		{privateKey3.Point, true, expected3, "Compressed hash160"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := tc.point.Hash160(tc.compressed)
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("Address() returned %x, expected %x", result, tc.expected)
			}
		})
	}
}

func TestParseSEC(t *testing.T) {
	type testCase struct {
		name string
		sec  []byte
	}

	// Create a slice of test cases
	testCases := []testCase{
		{
			name: "Test Case 1",
			sec:  []byte{4, 255, 229, 88, 227, 136, 133, 47, 1, 32, 228, 106, 242, 209, 179, 112, 248, 88, 84, 168, 235, 8, 65, 129, 30, 206, 14, 62, 3, 210, 130, 213, 124, 49, 93, 199, 40, 144, 164, 241, 10, 20, 129, 192, 49, 176, 59, 53, 27, 13, 199, 153, 1, 202, 24, 160, 12, 240, 9, 219, 219, 21, 122, 29, 16},
		},
		{
			name: "Test Case 2",
			sec:  []byte{4, 2, 127, 61, 161, 145, 132, 85, 224, 60, 70, 246, 89, 38, 106, 27, 181, 32, 78, 149, 157, 183, 54, 77, 47, 71, 59, 223, 143, 10, 19, 204, 157, 255, 135, 100, 127, 208, 35, 193, 59, 74, 73, 148, 241, 118, 145, 137, 88, 6, 225, 180, 11, 87, 244, 253, 34, 88, 26, 79, 70, 133, 31, 59, 6},
		},
		{
			name: "Test Case 3",
			sec:  []byte{4, 217, 12, 214, 37, 238, 135, 221, 56, 101, 109, 217, 92, 247, 159, 101, 246, 15, 114, 115, 182, 125, 48, 150, 230, 139, 216, 30, 79, 83, 66, 105, 31, 132, 46, 250, 118, 47, 213, 153, 97, 208, 233, 152, 3, 198, 30, 219, 168, 179, 227, 247, 220, 58, 52, 24, 54, 249, 119, 51, 174, 191, 152, 113, 33},
		},
		{
			name: "Test Case 4",
			sec:  []byte{3, 87, 164, 243, 104, 134, 138, 138, 109, 87, 41, 145, 228, 132, 230, 100, 129, 15, 241, 76, 5, 192, 250, 2, 50, 117, 37, 17, 81, 254, 14, 83, 209},
		},
		{
			name: "Test Case 5",
			sec:  []byte{2, 147, 62, 194, 210, 177, 17, 185, 39, 55, 236, 18, 241, 197, 210, 15, 50, 51, 160, 173, 33, 205, 139, 54, 208, 188, 167, 160, 207, 165, 203, 135, 1},
		},
		{
			name: "Test Case 6",
			sec:  []byte{2, 150, 190, 91, 18, 146, 246, 200, 86, 179, 197, 101, 78, 136, 111, 193, 53, 17, 70, 32, 89, 8, 156, 223, 156, 71, 150, 35, 191, 203, 231, 118, 144},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseSEC(tc.sec)
			if err != nil {
				t.Errorf("SEC '%v' could not be parsed into a point: %v", tc.sec, err)
			}

		})
	}
}

// This test shows the importance of choosing a random k every time you sign.
// If our secret is e and we are reusing k to sign z1 and z2:
// kG = (r,y)
// s1 = (z1 + re)/k, s2 = (z2 + re)/k
// s1/s2 = (z1 + re) / (z2 + re)
// s1(z2 + re) = s2(z1 +re)
// s1z2 + s1re = s2z1 + s2re
// s1re - s2re = s2z1 - s1z2
// e = (s2z1 - s1z2) / (s1r - s2r)
// func TestImportanceOfUniqueK(t *testing.T) {
// 	e := utils.Hash256ToBigInt("my secret")
// 	z1 := utils.Hash256ToBigInt("my first message")
// 	z2 := utils.Hash256ToBigInt("my second message")
// 	k := big.NewInt(1234567890)
// 	sig1, _ := Sign(e, z1)
// 	sig2, _ := Sign(e, z2)
// 	if sig1.R.Cmp(sig2.R) != 0 {
// 		t.Error("Same random number k should lead to same r")
// 	}
// 	r := sig1.R
// 	s2z1 := new(big.Int).Mul(sig2.S, z1)
// 	s1z2 := new(big.Int).Mul(sig1.S, z2)
// 	rs1 := new(big.Int).Mul(r, sig1.S)
// 	rs2 := new(big.Int).Mul(r, sig2.S)
// 	num := new(big.Int).Sub(s2z1, s1z2)
// 	denom := new(big.Int).ModInverse(new(big.Int).Sub(rs1, rs2), N)
// 	found_e := new(big.Int).Mul(num, denom)
// 	if e.Cmp(new(big.Int).Mod(found_e, N)) != 0 {
// 		t.Error("Could not retrieve private key even though same k was used to sign two messages.")
// 	}
// }

func TestNewPrivateKey(t *testing.T) {
	// Test with valid input
	validSecret := big.NewInt(12345)
	expectedPoint, _ := G.ScalarMultiplication(validSecret)
	privKey, err := NewPrivateKey(validSecret)
	if err != nil {
		t.Errorf("NewPrivateKey with valid input returned an error: %v", err)
	}
	if privKey == nil {
		t.Errorf("NewPrivateKey with valid input returned a nil PrivateKey")
	}
	if !privKey.Point.Equal(&expectedPoint.Point) {
		t.Errorf("NewPrivateKey with valid input returned an incorrect public point")
	}
}

func TestGetDeterministicK(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name         string
		secret       *big.Int
		z            *big.Int
		expectedKHex string
	}{
		{
			name:         "Test Case 1",
			secret:       Hash256ToBigInt("my secret"),
			z:            Hash256ToBigInt("Hi Mom!"),                                           // Example value for z
			expectedKHex: "0x5a36ac7d11fc415802c6049fda6ced159feb2044ba9bc61ecb18c8366b64ac65", // Expected output (this should be pre-calculated for a known input)
		},
		// Add more test cases as necessary
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := NewPrivateKey(tc.secret)
			if err != nil {
				t.Errorf("failed to create private key: %v", err)
			}
			actualK := e.GetDeterministicK(tc.z)
			expectedK, _ := new(big.Int).SetString(tc.expectedKHex, 0)
			if actualK.Cmp(expectedK) != 0 {
				t.Errorf("GetDeterministicK(%v, %v) = %v, want %v", tc.secret, tc.z, actualK, tc.expectedKHex)
			}
		})
	}
}

func TestPrivateKeySerialize(t *testing.T) {
	privateKey1, _ := NewPrivateKey(big.NewInt(5003))
	privateKey2, _ := NewPrivateKey(new(big.Int).Exp(big.NewInt(2021), big.NewInt(5), nil))
	secret3, _ := new(big.Int).SetString("0x54321deadbeef", 0)
	privateKey3, _ := NewPrivateKey(secret3)

	expected1 := "cMahea7zqjxrtgAbB7LSGbcQUr1uX1ojuat9jZodMN8rFTv2sfUK"
	expected2 := "91avARGdfge8E4tZfYLoxeJ5sGBdNJQH4kvjpWAxgzczjbCwxic"
	expected3 := "KwDiBf89QgGbjEhKnhXJuH7LrciVrZi3qYjgiuQJv1h8Ytr2S53a"

	// Test cases
	testCases := []struct {
		privateKey  *PrivateKey
		compressed  bool
		testnet     bool
		expected    string
		description string
	}{
		{privateKey1, true, true, expected1, "Compressed, testnet"},
		{privateKey2, false, true, expected2, "Uncompressed, testnet"},
		{privateKey3, true, false, expected3, "Compressed, mainnet"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := tc.privateKey.Serialize(tc.compressed, tc.testnet)
			if result != tc.expected {
				t.Errorf("Serialize() returned \n\n%s, expected \n\n%s", result, tc.expected)
			}
		})
	}
}

// parseSignatureParts parses the signature components from hex strings.
func parseSignatureParts(parts map[string]string) (z, r, s, x, y *big.Int, err error) {
	z, ok := new(big.Int).SetString(parts["z"], 0)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("invalid z value")
	}
	r, ok = new(big.Int).SetString(parts["r"], 0)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("invalid r value")
	}
	s, ok = new(big.Int).SetString(parts["s"], 0)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("invalid s value")
	}
	x, ok = new(big.Int).SetString(parts["x"], 0)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("invalid x value")
	}
	y, ok = new(big.Int).SetString(parts["y"], 0)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("invalid y value")
	}
	return z, r, s, x, y, nil
}

// createEllipticCurvePoint creates a new elliptic curve point from big.Ints.
func createEllipticCurvePoint(x, y *big.Int) (*S256Point, error) {
	px, err := NewS256FieldElement(x)
	if err != nil {
		return nil, err
	}
	py, err := NewS256FieldElement(y)
	if err != nil {
		return nil, err
	}
	return NewS256Point(px, py)
}

func Hash256ToBigInt(data string) *big.Int {
	// First SHA-256 hash
	hash256 := utils.Hash256([]byte(data))

	// Convert the second hash bytes to a big.Int
	bigInt := new(big.Int)
	bigInt.SetBytes(hash256)
	return bigInt
}
