package signatureverification

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/caspereijkens/cryptocurrency/internal/ellipticcurve"
	"github.com/caspereijkens/cryptocurrency/internal/util"
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

func TestSignAndVerify(t *testing.T) {
	privKey, err := NewPrivateKey(util.Hash256ToBigInt("my secret"))
	if err != nil {
		t.Errorf("failed to create private key: %v", err)
	}
	z := util.Hash256ToBigInt("my message")
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
// 	e := util.Hash256ToBigInt("my secret")
// 	z1 := util.Hash256ToBigInt("my first message")
// 	z2 := util.Hash256ToBigInt("my second message")
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
			secret:       util.Hash256ToBigInt("my secret"),
			z:            util.Hash256ToBigInt("Hi Mom!"),                                      // Example value for z
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
