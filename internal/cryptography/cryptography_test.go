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

func TestDivide(t *testing.T) {
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
