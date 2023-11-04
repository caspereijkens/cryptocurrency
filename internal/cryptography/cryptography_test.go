package cryptography

import (
	"math/big"
	"testing"
)

func TestNewFieldElement(t *testing.T) {
	// Test case 1: Valid input values
	value := big.NewInt(7)
	prime := big.NewInt(17)
	fe, err := NewFieldElement(value, prime)
	if err != nil {
		t.Errorf("NewFieldElement returned an error for valid input: %v", err)
	}
	if !fe.Equal(&FieldElement{value, prime}) {
		t.Errorf("NewFieldElement did not create the expected FieldElement")
	}

	// Test case 2: Value out of range
	value = big.NewInt(17)
	_, err = NewFieldElement(value, prime)
	if err == nil {
		t.Error("NewFieldElement did not return an error for out-of-range value")
	}
}

func TestFieldElementAdd(t *testing.T) {
	// Test case 1: Add two field elements with the same prime
	a, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	b, _ := NewFieldElement(big.NewInt(8), big.NewInt(17))
	result, err := a.Add(b)
	if err != nil {
		t.Errorf("FieldElement Add returned an error: %v", err)
	}
	expected, _ := NewFieldElement(big.NewInt(15), big.NewInt(17))
	if !result.Equal(expected) {
		t.Errorf("FieldElement Add result is not as expected")
	}

	// Test case 2: Add two field elements with different primes
	c, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	d, _ := NewFieldElement(big.NewInt(8), big.NewInt(19))
	_, err = c.Add(d)
	if err == nil {
		t.Error("FieldElement Add did not return an error for different primes")
	}
}

func TestFieldElementSubtract(t *testing.T) {
	// Test case 1: Subtract two field elements with the same prime
	a, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	b, _ := NewFieldElement(big.NewInt(8), big.NewInt(17))
	result, err := a.Subtract(b)
	if err != nil {
		t.Errorf("FieldElement Subtract returned an error: %v", err)
	}
	expected, _ := NewFieldElement(big.NewInt(16), big.NewInt(17))
	if !result.Equal(expected) {
		t.Errorf("FieldElement Subtract result is not as expected")
	}

	// Test case 2: Subtract two field elements with different primes
	c, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	d, _ := NewFieldElement(big.NewInt(8), big.NewInt(19))
	_, err = c.Subtract(d)
	if err == nil {
		t.Error("FieldElement Subtract did not return an error for different primes")
	}
}

func TestFieldElementMultiply(t *testing.T) {
	// Test case 1: Add two field elements with the same prime
	a, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	b, _ := NewFieldElement(big.NewInt(8), big.NewInt(17))
	result, err := a.Multiply(b)
	if err != nil {
		t.Errorf("FieldElement Multiply returned an error: %v", err)
	}
	expected, _ := NewFieldElement(big.NewInt(5), big.NewInt(17))
	if !result.Equal(expected) {
		t.Errorf("FieldElement Multiply result is not as expected: %s", result.String())
	}

	// Test case 2: Multiply two field elements with different primes
	c, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	d, _ := NewFieldElement(big.NewInt(8), big.NewInt(19))
	_, err = c.Multiply(d)
	if err == nil {
		t.Error("FieldElement Multiply did not return an error for different primes")
	}
}

func TestFieldElementExponentiate(t *testing.T) {
	// Test exponentiation of a field element.
	base, _ := NewFieldElement(big.NewInt(3), big.NewInt(17))
	power := big.NewInt(3)
	result, err := base.Exponentiate(power)
	if err != nil {
		t.Errorf("FieldElement Exponentiate returned an error: %v", err)
	}
	expected, _ := NewFieldElement(big.NewInt(10), big.NewInt(17))
	if !result.Equal(expected) {
		t.Errorf("FieldElement Multiply result is not as expected: %s", result.String())
	}
}

func TestFieldElementEqual(t *testing.T) {
	// Test equality of two field elements.
	a, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	b, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	if !a.Equal(b) {
		t.Error("FieldElement Equal returned false for equal field elements")
	}

	c, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	d, _ := NewFieldElement(big.NewInt(8), big.NewInt(17))
	if c.Equal(d) {
		t.Error("FieldElement Equal returned true for different field elements")
	}
}

func TestFieldElementString(t *testing.T) {
	// Test string representation of a field element.
	a, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	expected := "FieldElement_17(7)"
	if a.String() != expected {
		t.Errorf("FieldElement String representation is not as expected")
	}
}

func TestFieldElementNegate(t *testing.T) {
	prime := big.NewInt(13)
	tests := []struct {
		inputValue    *big.Int
		expectedValue *big.Int
	}{
		{big.NewInt(7), big.NewInt(6)},
		{big.NewInt(0), big.NewInt(0)},
		{big.NewInt(12), big.NewInt(1)},
	}

	for _, test := range tests {
		fe, err := NewFieldElement(test.inputValue, prime)
		if err != nil {
			t.Fatalf("Error creating FieldElement: %v", err)
		}

		negatedFe := fe.Negate()
		if negatedFe.value.Cmp(test.expectedValue) != 0 {
			t.Errorf("Negate(%v) => %v, expected %v", test.inputValue, negatedFe.value, test.expectedValue)
		}
	}
}

func TestFieldElementDivide(t *testing.T) {
	// Create two FieldElement instances
	prime := big.NewInt(19) // Replace with your desired prime value
	a, _ := NewFieldElement(big.NewInt(2), prime)
	b, _ := NewFieldElement(big.NewInt(7), prime)

	// Test division of a by b
	result, err := a.Divide(b)

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Expected result for a / b in this case should be 6 / 2 = 3
	expected, _ := NewFieldElement(big.NewInt(3), prime)

	// Compare the result with the expected value
	if result.value.Cmp(expected.value) != 0 {
		t.Errorf("Expected result %v, but got %v", expected.value, result.value)
	}

	// Test division by zero
	zero, _ := NewFieldElement(big.NewInt(0), prime)
	_, err = a.Divide(zero)

	// Check for the division by zero error
	if err == nil {
		t.Error("Expected division by zero error, but got no error")
	}

	// Test division by different fields
	otherPrime := big.NewInt(17) // Different prime value
	c, _ := NewFieldElement(big.NewInt(3), otherPrime)

	_, err = a.Divide(c)

	// Check for the different fields error
	if err == nil {
		t.Error("Expected different fields error, but got no error")
	}
}

func TestNewPoint(t *testing.T) {
	prime := big.NewInt(223)
	a, _ := NewFieldElement(big.NewInt(0), prime)
	b, _ := NewFieldElement(big.NewInt(7), prime)
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
		x, _ := NewFieldElement(point[0], prime)
		y, _ := NewFieldElement(point[1], prime)
		_, err := NewPoint(x, y, a, b)
		if err != nil {
			t.Errorf("NewPoint returned an error for valid input: %v", err)
		}
	}

	// Test case 2: Invalid input values
	for _, point := range invalidPoints {
		x, _ := NewFieldElement(point[0], prime)
		y, _ := NewFieldElement(point[1], prime)
		_, err := NewPoint(x, y, a, b)
		if err == nil {
			t.Errorf("NewPoint did not return an error for invalid input: %v", err)
		}
	}

	// Test case 3: Point at infintity
	var inf *FieldElement
	_, err := NewPoint(inf, inf, a, b)
	if err != nil {
		t.Errorf("NewPoint did not return an error for invalid input: %v", err)
	}

	// Test case 4: Mal-formed ecc parameters
	x, _ := NewFieldElement(big.NewInt(17), prime)
	y, _ := NewFieldElement(big.NewInt(56), prime)
	_, err = NewPoint(x, y, inf, inf)
	if err == nil {
		t.Errorf("ECC parameters are malformed and NewPoint did not return an error for invalid input: %v", err)
	}

}

func TestPointEqual(t *testing.T) {
	prime := big.NewInt(223)
	a1, _ := NewFieldElement(big.NewInt(0), prime)
	b1, _ := NewFieldElement(big.NewInt(7), prime)
	x1, _ := NewFieldElement(big.NewInt(17), prime)
	y1, _ := NewFieldElement(big.NewInt(56), prime)

	// Test case 1: equality of two equal field elements.
	p1, _ := NewPoint(x1, y1, a1, b1)
	q1, _ := NewPoint(x1, y1, a1, b1)
	if !p1.Equal(q1) {
		t.Error("Point Equal returned false for equal field elements")
	}

	// Test case 2: inequality of two inequal field elements
	x2, _ := NewFieldElement(big.NewInt(192), prime)
	y2, _ := NewFieldElement(big.NewInt(105), prime)
	q2, _ := NewPoint(x2, y2, a1, b1)
	if p1.Equal(q2) {
		t.Error("Point Equal returned true for different field elements")
	}

	// Test case 3: equality of two identity points
	var inf *FieldElement
	pInf, _ := NewPoint(inf, inf, a1, b1)
	if !pInf.Equal(&Point{nil, nil, a1, b1}) {
		t.Error("Point is mistakenly not marked as equal to the point at infinity")
	}

	// Test case 4: inequality of identity and non-identity
	if p1.Equal(pInf) {
		t.Error("Point is marked as equal to the point at infinity")
	}

	// Test case 5: inequality of points on different elliptic curves
	a3, _ := NewFieldElement(big.NewInt(5), prime)
	b3, _ := NewFieldElement(big.NewInt(7), prime)
	x3, _ := NewFieldElement(big.NewInt(-1), prime)
	y3, _ := NewFieldElement(big.NewInt(-1), prime)
	q3, _ := NewPoint(x3, y3, a3, b3)
	if p1.Equal(q3) {
		t.Error("Point Equal returned true for points on different elliptic curves")
	}
}

func TestPointString(t *testing.T) {
	prime := big.NewInt(223)
	a, _ := NewFieldElement(big.NewInt(0), prime)
	b, _ := NewFieldElement(big.NewInt(7), prime)
	x, _ := NewFieldElement(big.NewInt(17), prime)
	y, _ := NewFieldElement(big.NewInt(56), prime)
	p, _ := NewPoint(x, y, a, b)
	// Test string representation of a point.
	expected := "Point_0_7(17,56) FieldElement(223)"
	if p.String() != expected {
		t.Errorf("Point String representation (%s) is not as expected", p.String())
	}
}

func TestPointIsIdentityElement(t *testing.T) {
	var inf *FieldElement
	prime := big.NewInt(223)
	a, _ := NewFieldElement(big.NewInt(0), prime)
	b, _ := NewFieldElement(big.NewInt(7), prime)

	// Test 1 Verify identity point returns true
	pInf, _ := NewPoint(inf, inf, a, b)
	if !pInf.IsIdentityElement() {
		t.Error("Point is not marked as the point at infinity")
	}

	// Test 2 Verify non-identity point returns false
	x, _ := NewFieldElement(big.NewInt(17), prime)
	y, _ := NewFieldElement(big.NewInt(56), prime)
	p, _ := NewPoint(x, y, a, b)
	if p.IsIdentityElement() {
		t.Error("Point is mistakenly marked as the point at infinity")
	}
}

func TestEqualEllipticCurve(t *testing.T) {
	// Create two points for testing
	prime := big.NewInt(223)
	a1, _ := NewFieldElement(big.NewInt(0), prime)
	b1, _ := NewFieldElement(big.NewInt(7), prime)
	a2, _ := NewFieldElement(big.NewInt(5), prime)
	b2, _ := NewFieldElement(big.NewInt(7), prime)
	x1, _ := NewFieldElement(big.NewInt(17), prime)
	y1, _ := NewFieldElement(big.NewInt(56), prime)
	x2, _ := NewFieldElement(big.NewInt(-1), prime)
	y2, _ := NewFieldElement(big.NewInt(-1), prime)
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

func TestPointAdd(t *testing.T) {
	prime := big.NewInt(223)
	a, _ := NewFieldElement(big.NewInt(0), prime)
	b, _ := NewFieldElement(big.NewInt(7), prime)
	x1, _ := NewFieldElement(big.NewInt(192), prime)
	y1, _ := NewFieldElement(big.NewInt(105), prime)
	p1, _ := NewPoint(x1, y1, a, b)
	p1_inv, _ := NewPoint(x1, y1.Negate(), a, b)
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
	a2, _ := NewFieldElement(big.NewInt(5), prime)
	b2, _ := NewFieldElement(big.NewInt(7), prime)
	x2, _ := NewFieldElement(big.NewInt(-1), prime)
	y2, _ := NewFieldElement(big.NewInt(-1), prime)
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
	x2, _ = NewFieldElement(big.NewInt(17), prime)
	y2, _ = NewFieldElement(big.NewInt(56), prime)
	p2, _ = NewPoint(x2, y2, a, b)
	result, err = p1.Add(p2)
	if err != nil {
		t.Errorf("Point Add returned an error: %v", err)
	}
	x3, _ := NewFieldElement(big.NewInt(170), prime)
	y3, _ := NewFieldElement(big.NewInt(142), prime)
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
