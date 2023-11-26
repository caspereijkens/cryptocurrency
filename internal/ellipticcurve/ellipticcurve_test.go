package ellipticcurve

import (
	"math/big"
	"testing"

	"github.com/caspereijkens/cryptocurrency/internal/finitefield"
)

func TestNewPoint(t *testing.T) {
	prime := big.NewInt(223)
	a, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	validPoints := [][]*big.Int{
		{big.NewInt(192), big.NewInt(105)},
		{big.NewInt(17), big.NewInt(56)},
		{big.NewInt(1), big.NewInt(193)},
	}
	invalidPoints := [][]*big.Int{
		{big.NewInt(200), big.NewInt(119)},
		{big.NewInt(42), big.NewInt(99)},
	}

	// Test case 1: Valid input values
	for _, point := range validPoints {
		x, _ := finitefield.NewFieldElement(point[0], prime)
		y, _ := finitefield.NewFieldElement(point[1], prime)
		_, err := NewPoint(x, y, a, b)
		if err != nil {
			t.Errorf("NewPoint returned an error for valid input: %v", err)
		}
	}

	// Test case 2: Invalid input values
	for _, point := range invalidPoints {
		x, _ := finitefield.NewFieldElement(point[0], prime)
		y, _ := finitefield.NewFieldElement(point[1], prime)
		_, err := NewPoint(x, y, a, b)
		if err == nil {
			t.Errorf("NewPoint did not return an error for invalid input: %v", err)
		}
	}

	// Test case 3: Point at infintity
	var inf *finitefield.FieldElement
	_, err := NewPoint(inf, inf, a, b)
	if err != nil {
		t.Errorf("NewPoint did not return an error for invalid input: %v", err)
	}

	// Test case 4: Mal-formed ecc parameters
	x, _ := finitefield.NewFieldElement(big.NewInt(17), prime)
	y, _ := finitefield.NewFieldElement(big.NewInt(56), prime)
	_, err = NewPoint(x, y, inf, inf)
	if err == nil {
		t.Errorf("ECC parameters are malformed and NewPoint did not return an error for invalid input: %v", err)
	}

}

func TestPointEqual(t *testing.T) {
	prime := big.NewInt(223)
	a1, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b1, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x1, _ := finitefield.NewFieldElement(big.NewInt(17), prime)
	y1, _ := finitefield.NewFieldElement(big.NewInt(56), prime)

	// Test case 1: equality of two equal field elements.
	p1, _ := NewPoint(x1, y1, a1, b1)
	q1, _ := NewPoint(x1, y1, a1, b1)
	if !p1.Equal(q1) {
		t.Error("Point Equal returned false for equal field elements")
	}

	// Test case 2: inequality of two inequal field elements
	x2, _ := finitefield.NewFieldElement(big.NewInt(192), prime)
	y2, _ := finitefield.NewFieldElement(big.NewInt(105), prime)
	q2, _ := NewPoint(x2, y2, a1, b1)
	if p1.Equal(q2) {
		t.Error("Point Equal returned true for different field elements")
	}

	// Test case 3: equality of two identity points
	var inf *finitefield.FieldElement
	pInf, _ := NewPoint(inf, inf, a1, b1)
	if !pInf.Equal(&Point{nil, nil, a1, b1}) {
		t.Error("Point is mistakenly not marked as equal to the point at infinity")
	}

	// Test case 4: inequality of identity and non-identity
	if p1.Equal(pInf) {
		t.Error("Point is marked as equal to the point at infinity")
	}

	// Test case 5: inequality of points on different elliptic curves
	a3, _ := finitefield.NewFieldElement(big.NewInt(5), prime)
	b3, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x3, _ := finitefield.NewFieldElement(big.NewInt(-1), prime)
	y3, _ := finitefield.NewFieldElement(big.NewInt(-1), prime)
	q3, _ := NewPoint(x3, y3, a3, b3)
	if p1.Equal(q3) {
		t.Error("Point Equal returned true for points on different elliptic curves")
	}
}

func TestPointString(t *testing.T) {
	prime := big.NewInt(223)
	a, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x, _ := finitefield.NewFieldElement(big.NewInt(17), prime)
	y, _ := finitefield.NewFieldElement(big.NewInt(56), prime)
	p, _ := NewPoint(x, y, a, b)
	// Test string representation of a point.
	expected := "Point_0_7(17,56) Field_223"
	if p.String() != expected {
		t.Errorf("Point String representation (%s) is not as expected (%s)", p.String(), expected)
	}
}

func TestPointIsIdentityElement(t *testing.T) {
	var inf *finitefield.FieldElement
	prime := big.NewInt(223)
	a, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b, _ := finitefield.NewFieldElement(big.NewInt(7), prime)

	// Test 1 Verify identity point returns true
	pInf, _ := NewPoint(inf, inf, a, b)
	if !pInf.IsIdentityElement() {
		t.Error("Point is not marked as the point at infinity")
	}

	// Test 2 Verify non-identity point returns false
	x, _ := finitefield.NewFieldElement(big.NewInt(17), prime)
	y, _ := finitefield.NewFieldElement(big.NewInt(56), prime)
	p, _ := NewPoint(x, y, a, b)
	if p.IsIdentityElement() {
		t.Error("Point is mistakenly marked as the point at infinity")
	}
}

func TestEqualEllipticCurve(t *testing.T) {
	// Create two points for testing
	prime := big.NewInt(223)
	a1, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b1, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	a2, _ := finitefield.NewFieldElement(big.NewInt(5), prime)
	b2, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x1, _ := finitefield.NewFieldElement(big.NewInt(17), prime)
	y1, _ := finitefield.NewFieldElement(big.NewInt(56), prime)
	x2, _ := finitefield.NewFieldElement(big.NewInt(-1), prime)
	y2, _ := finitefield.NewFieldElement(big.NewInt(-1), prime)
	p, _ := NewPoint(x1, y1, a1, b1)
	q, _ := NewPoint(x2, y2, a2, b2)

	// Test equal points
	if !p.EqualEllipticCurve(p) {
		t.Errorf("Expected points to be equal, but they are not.")
	}

	// Test different points
	if p.EqualEllipticCurve(q) {
		t.Errorf("Expected points to be different, but they are equal.")
	}
}

// TestCalculatedxdy tests the calculatedxdy function
func TestCalculatedxdy(t *testing.T) {
	prime := big.NewInt(223)
	a, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x1, _ := finitefield.NewFieldElement(big.NewInt(17), prime)
	y1, _ := finitefield.NewFieldElement(big.NewInt(56), prime)
	x2, _ := finitefield.NewFieldElement(big.NewInt(49), prime)
	y2, _ := finitefield.NewFieldElement(big.NewInt(71), prime)
	p, _ := NewPoint(x1, y1, a, b)
	q, _ := NewPoint(x2, y2, a, b)

	// Test with the point equal to itself
	_, _, err := p.calculatedxdy(p)
	if err != nil {
		t.Fatalf("Failed to calculate dx and dy: %s", err)
	}
	// Add assertions to check if dx and dy are correct
	// For example:
	// expectedDx := /* some expected value */
	// expectedDy := /* some expected value */
	// if !dx.Equal(expectedDx) || !dy.Equal(expectedDy) {
	// 	t.Errorf("calculatedxdy did not calculate the expected results for equal points")
	// }

	// Test case 2:  point not equal to itself
	// You would create a different point q with its x and y field elements here

	_, _, err = p.calculatedxdy(q)
	if err != nil {
		t.Fatalf("Failed to calculate dx and dy: %s", err)
	}
	// Add assertions to check if dx and dy are correct for unequal points
}

func TestIsVerticalTangent(t *testing.T) {
	// Create a test case where p is a point with y = 0
	prime := big.NewInt(223)
	a, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x1, _ := finitefield.NewFieldElement(big.NewInt(6), prime)
	y1, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	p, _ := NewPoint(x1, y1, a, b)

	// Test if p is a vertical tangent to q
	result := p.isVerticalTangent(p)

	// Expected result is true, as p and q have the same x value and y = 0
	expectedResult := true

	if result != expectedResult {
		t.Errorf("Expected %v to be a vertical tangent to %v, but got %v", p, p, result)
	}

	// Create a test case where q has y != 0
	x2, _ := finitefield.NewFieldElement(big.NewInt(1), prime)
	y2, _ := finitefield.NewFieldElement(big.NewInt(30), prime)
	q, _ := NewPoint(x2, y2, a, b)

	// Test if p is a vertical tangent to q
	result = p.isVerticalTangent(q)

	// Expected result is false, as q has y != 0
	expectedResult = false

	if result != expectedResult {
		t.Errorf("Expected %v to not be a vertical tangent to %v, but got %v", p, q, result)
	}
}

func TestPointAdd(t *testing.T) {
	prime := big.NewInt(223)
	a, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x1, _ := finitefield.NewFieldElement(big.NewInt(192), prime)
	y1, _ := finitefield.NewFieldElement(big.NewInt(105), prime)
	p1, _ := NewPoint(x1, y1, a, b)
	y1_neg, _ := y1.Negate()
	p1_inv, _ := NewPoint(x1, y1_neg, a, b)
	identity, _ := NewPoint(nil, nil, a, b)

	// Test case 1: Add inverse of point and check if this adds up to point at infinity.
	result, err := p1.Add(p1_inv)
	if err != nil {
		t.Errorf("Point Add returned an error: %v", err)
	}
	if !result.Equal(identity) {
		t.Errorf("Point Add result is not as expected")
	}

	// Test case 2: Add two field elements on a different Elliptic Curve and verify that this raises an error
	a2, _ := finitefield.NewFieldElement(big.NewInt(5), prime)
	b2, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x2, _ := finitefield.NewFieldElement(big.NewInt(-1), prime)
	y2, _ := finitefield.NewFieldElement(big.NewInt(-1), prime)
	p2, _ := NewPoint(x2, y2, a2, b2)
	_, err = p1.Add(p2)
	if err == nil {
		t.Errorf("Point Add returned no error, but Points are from different curves")
	}

	// Test case 3: Add point at infinity to point
	result, err = p1.Add(identity)
	if err != nil {
		t.Errorf("Point Add to identiy returned an error: %v", err)
	}
	if !result.Equal(p1) {
		t.Errorf("Point Add result is not as expected")
	}

	// Test case 4: Add two points on the same elliptic curve
	x2, _ = finitefield.NewFieldElement(big.NewInt(17), prime)
	y2, _ = finitefield.NewFieldElement(big.NewInt(56), prime)
	p2, _ = NewPoint(x2, y2, a, b)
	result, err = p1.Add(p2)
	if err != nil {
		t.Errorf("Point Add returned an error: %v", err)
	}
	x3, _ := finitefield.NewFieldElement(big.NewInt(170), prime)
	y3, _ := finitefield.NewFieldElement(big.NewInt(142), prime)
	p3, _ := NewPoint(x3, y3, a, b)
	if !result.Equal(p3) {
		t.Errorf("Point Add result is not as expected")
	}

	// Test case 5: Test Commutative property
	result, err = p2.Add(p1)
	if err != nil {
		t.Errorf("Point Add returned an error: %v", err)
	}
	if !result.Equal(p3) {
		t.Errorf("Point Add result is not as expected")
	}

	// Test case 6: Add point to itself
	x4, _ := finitefield.NewFieldElement(big.NewInt(49), prime)
	y4, _ := finitefield.NewFieldElement(big.NewInt(71), prime)
	p4, _ := NewPoint(x4, y4, a, b)
	x5, _ := finitefield.NewFieldElement(big.NewInt(66), prime)
	y5, _ := finitefield.NewFieldElement(big.NewInt(111), prime)
	p5, _ := NewPoint(x5, y5, a, b)
	result, err = p4.Add(p4)
	if err != nil {
		t.Errorf("Point Add returned an error: %v", err)
	}
	if !result.Equal(p5) {
		t.Errorf("Point Add result is not as expected")
	}

	// // Test case 6: Test Associative property (Test was abandoned because I don't know good cases.)
	// p, _ = NewPoint(big.NewInt(-1), big.NewInt(-1), big.NewInt(5), big.NewInt(7))
	// q, _ = NewPoint(big.NewInt(2), big.NewInt(5), big.NewInt(5), big.NewInt(7))
	// r, _ := NewPoint(big.NewInt(0.25), big.NewInt(2.875), big.NewInt(5), big.NewInt(7))
	// result1, _ := p.Add(q)
	// result1, _ = result1.Add(r) // (p + q) + r
	// result2, _ := q.Add(r)
	// result2, _ = p.Add(result2) // p + (q + r)
	// if !result1.Equal(result2) {
	// 	t.Errorf("Point Add result is not as expected")
	// }

	// Test case 7: Add a point to itself (P_1=P_2)
}

func TestAddingtoInf(t *testing.T) {
	prime := big.NewInt(223)
	a, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x, _ := finitefield.NewFieldElement(big.NewInt(49), prime)
	y, _ := finitefield.NewFieldElement(big.NewInt(71), prime)
	p, _ := NewPoint(x, y, a, b)
	identity, _ := NewPoint(nil, nil, a, b)

	result := p
	for i := 1; i <= 20; i++ {
		result, _ = result.Add(p)
		// fmt.Printf("%d: %s\n", i, result.String())
	}
	if !result.Equal(identity) {
		t.Errorf("Point should be the identity point")
	}
}

// TODO Add unit test for copy

func TestScalarMultiplication(t *testing.T) {
	// setup
	prime := big.NewInt(223)
	a, _ := finitefield.NewFieldElement(big.NewInt(0), prime)
	b, _ := finitefield.NewFieldElement(big.NewInt(7), prime)
	x, _ := finitefield.NewFieldElement(big.NewInt(49), prime)
	y, _ := finitefield.NewFieldElement(big.NewInt(71), prime)
	p, _ := NewPoint(x, y, a, b)
	// Define test cases
	testCases := []struct {
		coefficient *big.Int
		expectedX   *big.Int
		expectedY   *big.Int
		expectError bool
	}{
		// Test cases with normal coefficients
		{big.NewInt(1), big.NewInt(49), big.NewInt(71), false},
		{big.NewInt(2), big.NewInt(66), big.NewInt(111), false},
		{big.NewInt(4), big.NewInt(207), big.NewInt(51), false},
		{big.NewInt(21), nil, nil, false},

		// Test case with coefficient 0 (edge case)
		{big.NewInt(0), nil, nil, false},

		// Test case with negative coefficient (error case)
		{big.NewInt(-1), nil, nil, true},

		// ... other test cases ...
	}

	for _, tc := range testCases {
		// Perform Scalar Multiplication
		result, err := p.ScalarMultiplication(tc.coefficient)
		if (err != nil) != tc.expectError {
			t.Errorf("Expected error: %v, got: %v for coefficient %d", tc.expectError, err, tc.coefficient)
		}
		expectedX, _ := finitefield.NewFieldElement(tc.expectedX, prime)
		expectedY, _ := finitefield.NewFieldElement(tc.expectedY, prime)
		expected, _ := NewPoint(expectedX, expectedY, a, b)

		// Assert results
		if !tc.expectError && !result.Equal(expected) {
			t.Errorf("Expected result %s, got %s for coefficient %s", expected.String(), result.String(), tc.coefficient.String())
		}
	}
}

func TestOrderofGenerator(t *testing.T) {
	// setup
	identity, _ := NewPoint(nil, nil, &A.FieldElement, &B.FieldElement)
	// Call your function
	result, err := G.ScalarMultiplication(N)

	if err != nil {
		t.Errorf("Calculation did not go correct.")
	}
	if !result.Equal(identity) {
		t.Errorf("Point should be the identity point")
	}
}
