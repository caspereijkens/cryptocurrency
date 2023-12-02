package signatureverification

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/caspereijkens/cryptocurrency/internal/ellipticcurve"
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

// func TestSignatureAlgorithm(t *testing.T) {
// 	// This function creates a signature.
// 	// It computes the parts of the signature r, s and z.
// 	// z is simply a message that is hashed. It is called the signature hash.
// 	// r is a random number that is created by scalar multiplication of G, the generator point G.
// 	// G generates a finite cyclic group. This operation on this group is point addition.
// 	// The identity element of this group is of course the identity element of the elliptic curve.

// 	e := util.Hash256ToBigInt("my secret")
// 	z := util.Hash256ToBigInt("my message")
// 	k := big.NewInt(1234567890)

// 	// Calculate the target R
// 	R, _ := G.ScalarMultiplication(k)

// 	// Calculate r, the x-value of target R
// 	r := R.X.Value

// 	// Calculate r * e
// 	re := new(big.Int).Mul(r, e)

// 	// Calculate re + z
// 	rePlusZ := new(big.Int).Add(re, z)

// 	// Calculate (re + z) * kInv
// 	kInv := new(big.Int).ModInverse(k, N)
// 	product := new(big.Int).Mul(rePlusZ, kInv)

// 	// Modulo with N to get the final result
// 	s := new(big.Int).Mod(product, N)

// 	point, _ := G.ScalarMultiplication(e)
// 	fmt.Printf("%s", point.String())
// 	fmt.Printf("%x", z)
// 	fmt.Printf("%x", r)
// 	fmt.Printf("%x", s)
// }
