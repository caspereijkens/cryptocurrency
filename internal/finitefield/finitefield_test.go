package finitefield

import (
	"fmt"
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

	// Test case 3: Value nil
	value = nil
	fe, err = NewFieldElement(value, prime)
	if err != nil {
		t.Error("NewFieldElement did not return an error for out-of-range value")
	}
	if fe != nil {
		t.Error("NewFieldElement did not return desired nil field element")
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
	// Test case 1: Multiply two field elements with the same prime
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
	c, _ := NewFieldElement(big.NewInt(8), big.NewInt(19))
	_, err = a.Multiply(c)
	if err == nil {
		t.Error("FieldElement Multiply did not return an error for different primes")
	}

	// Test case 3: Multiply a field element by zero
	d, _ := NewFieldElement(big.NewInt(7), big.NewInt(17))
	zeroFieldElement, _ := NewFieldElement(big.NewInt(0), big.NewInt(17))
	result, err = d.Multiply(zeroFieldElement)
	if err != nil {
		t.Errorf("FieldElement Multiply returned an error: %v", err)
	}
	expectedZero, _ := NewFieldElement(big.NewInt(0), big.NewInt(17))
	if !result.Equal(expectedZero) {
		t.Errorf("FieldElement Multiply result is not zero: %s", result.String())
	}

	// Test case 4: Multiply zero by a field element
	result, err = zeroFieldElement.Multiply(d)
	if err != nil {
		t.Errorf("FieldElement Multiply returned an error: %v", err)
	}
	if !result.Equal(expectedZero) {
		t.Errorf("FieldElement Multiply result is not zero: %s", result.String())
	}
}

func TestFieldElementExponentiate(t *testing.T) {
	// Test setup with different base values and powers
	testCases := []struct {
		base     int64
		power    int64
		expected int64
		prime    int64 // Using a prime number for the field
	}{
		{3, 3, 27, 53},         // Normal case
		{0, 5, 0, 53},          // Exponentiating zero
		{5, 0, 1, 53},          // Power of zero
		{6, 2, 36, 53},         // Square
		{2, 3, 8, 53},          // Cube
		{12, 1, 12, 53},        // Power of one
		{15, 3, 3375 % 53, 53}, // Larger numbers
		// Add more cases as necessary
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Base%dPower%d", tc.base, tc.power), func(t *testing.T) {
			base, _ := NewFieldElement(big.NewInt(tc.base), big.NewInt(tc.prime))
			power := big.NewInt(tc.power)
			expected, _ := NewFieldElement(big.NewInt(tc.expected), big.NewInt(tc.prime))

			result, err := base.Exponentiate(power)
			if err != nil {
				t.Errorf("FieldElement Exponentiate returned an error: %v", err)
			}
			if !result.Equal(expected) {
				t.Errorf("FieldElement Exponentiate result is not as expected: got %s, want %s", result.String(), expected.String())
			}
		})
	}
}

func TestFieldElementSquared(t *testing.T) {
	// Test setup
	prime := big.NewInt(17)
	testCases := []struct {
		name     string
		input    *big.Int
		expected *big.Int
	}{
		{"Square of 2", big.NewInt(2), big.NewInt(4)},
		{"Square of 0", big.NewInt(0), big.NewInt(0)},
		{"Square of 5", big.NewInt(5), big.NewInt(8)},
		{"Square of 6", big.NewInt(6), big.NewInt(2)},
		// Add more test cases...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a, _ := NewFieldElement(tc.input, prime)
			result, err := a.Squared()
			if err != nil {
				t.Errorf("FieldElement Exponentiate returned an error: %v", err)
			}
			expected, _ := NewFieldElement(tc.expected, prime)
			if !result.Equal(expected) {
				t.Errorf("FieldElement Multiply result is not as expected: %s", result.String())
			}
		})
	}
}

func TestFieldElementCubed(t *testing.T) {
	// Test setup
	prime := big.NewInt(17)
	testCases := []struct {
		name     string
		input    *big.Int
		expected *big.Int
	}{
		{"Cube of 2", big.NewInt(2), big.NewInt(8)},
		{"Cube of 0", big.NewInt(0), big.NewInt(0)},
		{"Cube of 5", big.NewInt(5), big.NewInt(6)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a, _ := NewFieldElement(tc.input, prime)
			result, err := a.Cubed()
			if err != nil {
				t.Errorf("FieldElement Exponentiate returned an error: %v", err)
			}
			expected, _ := NewFieldElement(tc.expected, prime)
			if !result.Equal(expected) {
				t.Errorf("FieldElement Multiply result is not as expected: %s", result.String())
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		value  string
		prime  string
		result string
	}{
		{"4", "17", "2"},
		{"9", "17", "14"},
		{"16", "17", "4"},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%s", test.value, test.prime), func(t *testing.T) {
			value, _ := new(big.Int).SetString(test.value, 10)
			prime, _ := new(big.Int).SetString(test.prime, 10)
			expected, _ := new(big.Int).SetString(test.result, 10)

			fieldElement := &FieldElement{Value: value, Prime: prime}
			result, err := fieldElement.Sqrt()

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.Value.Cmp(expected) != 0 {
				t.Errorf("Expected square root: %s, got: %s", expected, result.Value)
			}
		})
	}
}

func TestGetEvenOddSquareRoots(t *testing.T) {
	tests := []struct {
		value string
		prime string
		even  string
		odd   string
	}{
		{"4", "17", "2", "15"},
		{"9", "17", "14", "3"},
		{"16", "17", "4", "13"},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%s", test.value, test.prime), func(t *testing.T) {
			value, _ := new(big.Int).SetString(test.value, 10)
			prime, _ := new(big.Int).SetString(test.prime, 10)
			expectedEven, _ := new(big.Int).SetString(test.even, 10)
			expectedOdd, _ := new(big.Int).SetString(test.odd, 10)

			fieldElement := &FieldElement{Value: value, Prime: prime}
			even, odd, err := fieldElement.GetEvenOddSquareRoots()

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if even.Cmp(expectedEven) != 0 || odd.Cmp(expectedOdd) != 0 {
				t.Errorf("Expected even: %s, odd: %s; got even: %s, odd: %s", expectedEven, expectedOdd, even, odd)
			}
		})
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

		negatedFe, err := fe.Negate()
		if err != nil {
			t.Errorf("Error negating FieldElement: %v", fe)
		}
		if negatedFe.Value.Cmp(test.expectedValue) != 0 {
			t.Errorf("Negate(%v) => %v, expected %v", test.inputValue, negatedFe.Value, test.expectedValue)
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
	if result.Value.Cmp(expected.Value) != 0 {
		t.Errorf("Expected result %v, but got %v", expected.Value, result.Value)
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
