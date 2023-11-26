package finitefield

import (
	"fmt"
	"math/big"
)

// FieldElement represents an element in a finite field.
type FieldElement struct {
	Value *big.Int
	Prime *big.Int
}

// NewFieldElement creates a new FieldElement with the given value and prime.
func NewFieldElement(value, prime *big.Int) (*FieldElement, error) {
	if value == nil {
		return nil, nil
	}
	if value.Sign() < 0 || value.Cmp(prime) >= 0 {
		return nil, fmt.Errorf("value not in the range [0, prime-1]")
	}
	return &FieldElement{Value: new(big.Int).Set(value), Prime: new(big.Int).Set(prime)}, nil
}

// Add adds two field elements and returns a new field element.
func (a *FieldElement) Add(b *FieldElement) (*FieldElement, error) {
	if a.Prime.Cmp(b.Prime) != 0 {
		return nil, fmt.Errorf("field elements are from different fields")
	}
	result := new(big.Int).Mod(new(big.Int).Add(a.Value, b.Value), a.Prime)
	return NewFieldElement(result, a.Prime)
}

// Subtract subtracts two field elements and returns a new field element.
func (a *FieldElement) Subtract(b *FieldElement) (*FieldElement, error) {
	if a.Prime.Cmp(b.Prime) != 0 {
		return nil, fmt.Errorf("field elements are from different fields")
	}
	result := new(big.Int).Sub(a.Value, b.Value)
	if result.Sign() < 0 {
		result.Add(result, a.Prime)
	}
	return NewFieldElement(result, a.Prime)
}

// Multiply multiplies two field elements and returns a new field element.
func (a *FieldElement) Multiply(b *FieldElement) (*FieldElement, error) {
	if a.Prime.Cmp(b.Prime) != 0 {
		return nil, fmt.Errorf("field elements are from different fields")
	}
	result := new(big.Int).Mul(a.Value, b.Value)
	result.Mod(result, a.Prime)
	return NewFieldElement(result.Mod(result, a.Prime), a.Prime)
}

// Exponentiate computes the exponentiation of a field element to a given power.
func (a *FieldElement) Exponentiate(power *big.Int) (*FieldElement, error) {
	result := new(big.Int).Exp(a.Value, power, a.Prime)
	return NewFieldElement(result.Mod(result, a.Prime), a.Prime)
}

// Squared computes the square of a field element.
func (a *FieldElement) Squared() (*FieldElement, error) {
	return a.Exponentiate(big.NewInt(2))
}

func (a *FieldElement) Cubed() (*FieldElement, error) {
	return a.Exponentiate(big.NewInt(3))
}

// Equal checks if two field elements are equal.
func (a *FieldElement) Equal(b *FieldElement) bool {
	return a.Value.Cmp(b.Value) == 0 && a.Prime.Cmp(b.Prime) == 0
}

// Negate returns a new FieldElement with the negated value of the current FieldElement.
func (a *FieldElement) Negate() (*FieldElement, error) {
	// Calculate the negated value as (prime - value) % prime
	negatedValue := new(big.Int).Sub(a.Prime, a.Value)
	return NewFieldElement(negatedValue.Mod(negatedValue, a.Prime), a.Prime)
}

// String returns the string representation of a field element.
func (a *FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%s(%s)", a.Prime.String(), a.Value.String())
}

// Divide computes the division of two field elements (a / b).
func (a *FieldElement) Divide(b *FieldElement) (*FieldElement, error) {
	if a.Prime.Cmp(b.Prime) != 0 {
		return nil, fmt.Errorf("field elements are from different fields")
	}
	if b.Value.Sign() == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	// Compute the modular multiplicative inverse of b
	inverse := new(big.Int).ModInverse(b.Value, a.Prime)
	if inverse == nil {
		return nil, fmt.Errorf("division by non-invertible element")
	}
	result := new(big.Int).Mul(a.Value, inverse)
	return NewFieldElement(result.Mod(result, a.Prime), a.Prime)
}
