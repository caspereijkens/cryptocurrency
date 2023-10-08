package cryptography

import (
	"errors"
	"fmt"
	"math/big"
)

// FieldElement represents an element in a finite field.
type FieldElement struct {
	value *big.Int
	prime *big.Int
}

// NewFieldElement creates a new FieldElement with the given value and prime.
func NewFieldElement(value, prime *big.Int) (*FieldElement, error) {
	if value.Sign() < 0 || value.Cmp(prime) >= 0 {
		return nil, errors.New("value not in the range [0, prime-1]")
	}
	return &FieldElement{value: new(big.Int).Set(value), prime: new(big.Int).Set(prime)}, nil
}

// Add adds two field elements and returns a new field element.
func (a *FieldElement) Add(b *FieldElement) (*FieldElement, error) {
	if a.prime.Cmp(b.prime) != 0 {
		return nil, errors.New("field elements are from different fields")
	}
	result := new(big.Int).Add(a.value, b.value)
	return NewFieldElement(result, a.prime)
}

// Subtract subtracts two field elements and returns a new field element.
func (a *FieldElement) Subtract(b *FieldElement) (*FieldElement, error) {
	if a.prime.Cmp(b.prime) != 0 {
		return nil, errors.New("field elements are from different fields")
	}
	result := new(big.Int).Sub(a.value, b.value)
	if result.Sign() < 0 {
		// If the result is negative, add the prime to wrap it around within the field
		result = result.Add(result, a.prime)
	}
	return NewFieldElement(result, a.prime)
}

// Multiply multiplies two field elements and returns a new field element.
func (a *FieldElement) Multiply(b *FieldElement) (*FieldElement, error) {
	if a.prime.Cmp(b.prime) != 0 {
		return nil, errors.New("field elements are from different fields")
	}
	result := new(big.Int).Mul(a.value, b.value)
	return NewFieldElement(result.Mod(result, a.prime), a.prime)
}

// Exponentiate computes the exponentiation of a field element to a given power.
func (a *FieldElement) Exponentiate(power *big.Int) (*FieldElement, error) {
	result := new(big.Int).Exp(a.value, power, a.prime)
	return NewFieldElement(result.Mod(result, a.prime), a.prime)
}

// Equal checks if two field elements are equal.
func (a *FieldElement) Equal(b *FieldElement) bool {
	return a.value.Cmp(b.value) == 0 && a.prime.Cmp(b.prime) == 0
}

// String returns the string representation of a field element.
func (a *FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%s(%s)", a.prime.String(), a.value.String())
}

// Divide computes the division of two field elements (a / b).
func (a *FieldElement) Divide(b *FieldElement) (*FieldElement, error) {
	if a.prime.Cmp(b.prime) != 0 {
		return nil, errors.New("field elements are from different fields")
	}
	if b.value.Sign() == 0 {
		return nil, errors.New("division by zero")
	}
	// Compute the modular multiplicative inverse of b
	inverse := new(big.Int).ModInverse(b.value, a.prime)
	if inverse == nil {
		return nil, errors.New("division by non-invertible element")
	}
	result := new(big.Int).Mul(a.value, inverse)
	return NewFieldElement(result.Mod(result, a.prime), a.prime)
}
