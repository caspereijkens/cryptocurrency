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

// Squared computes the exponentiation of a field element to a given power.
// TODO add unittest
func (a *FieldElement) Squared() (*FieldElement, error) {
	return a.Exponentiate(big.NewInt(2))
}

// Equal checks if two field elements are equal.
func (a *FieldElement) Equal(b *FieldElement) bool {
	return a.value.Cmp(b.value) == 0 && a.prime.Cmp(b.prime) == 0
}

// Negate returns a new FieldElement with the negated value of the current FieldElement.
func (a *FieldElement) Negate() *FieldElement {
	// Calculate the negated value as (prime - value) % prime
	negatedValue := new(big.Int).Sub(a.prime, a.value)
	negatedValue.Mod(negatedValue, a.prime)

	// Create a new FieldElement with the negated value and the same prime
	negatedFieldElement, _ := NewFieldElement(negatedValue, a.prime)

	return negatedFieldElement
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

// Point represents a point on the Elliptic Curve y^2 = x^3 + 7
type Point struct {
	x *FieldElement
	y *FieldElement
	a *FieldElement
	b *FieldElement
}

func NewPoint(x, y, a, b *FieldElement) (*Point, error) {
	// Check if a and b are well defined
	if a == nil || b == nil {
		return nil, fmt.Errorf("elliptic curve parameters are not well-defined")
	}

	// Check if this is the point at infinity
	if x == nil && y == nil {
		return &Point{nil, nil, a, b}, nil
	}

	// Check if the point (x, y) is on the elliptic curve y^2 = x^3 + ax + b
	xCubed, err := x.Exponentiate(big.NewInt(3))
	if err != nil {
		return nil, err
	}

	ax, err := a.Multiply(x)
	if err != nil {
		return nil, err
	}

	rightHandSide, err := xCubed.Add(ax)
	if err != nil {
		return nil, err
	}

	rightHandSide, err = rightHandSide.Add(b)
	if err != nil {
		return nil, err
	}

	ySquared, err := y.Exponentiate(big.NewInt(2))
	if err != nil {
		return nil, err
	}

	if !ySquared.Equal(rightHandSide) {
		return nil, fmt.Errorf("Point (%s, %s) does not exist on elliptic curve y^2 = x^3 + %s x + %s", x.String(), y.String(), a.String(), b.String()) //
	}

	return &Point{x, y, a, b}, nil
}

func (p *Point) IsIdentityElement() bool {
	return p.x == nil && p.y == nil
}

func (p *Point) Equal(q *Point) bool {
	// Check that the points are on the same line
	if !p.a.Equal(q.a) || !p.b.Equal(q.b) {
		return false
	}

	// Check if they're both identity elements
	if p.IsIdentityElement() && q.IsIdentityElement() {
		return true
	}

	// Check if they're both identity elements
	if p.IsIdentityElement() || q.IsIdentityElement() {
		return false
	}

	return p.x.Equal(q.x) && p.y.Equal(q.y)
}

func (p *Point) EqualEllipticCurve(q *Point) bool {
	return p.a.Equal(q.a) && p.b.Equal(q.b)
}

// func (p *Point) NotEqual(q *Point) bool {
// 	return !p.Equal(q)
// }

// String returns the string representation of a field element.
func (p *Point) String() string {
	return fmt.Sprintf("Point_%s_%s(%s,%s) FieldElement(%s)", p.a.value.String(), p.b.value.String(), p.x.value.String(), p.y.value.String(), p.x.prime.String())
}

// Copy returns a new Point with the same values as the current Point.
func (p *Point) Copy() (*Point, error) {
	// TODO add unittest
	return NewPoint(p.x, p.y, p.a, p.b)
}

// Add performs the addition of two elliptic curve points (p and q).
// It returns the resulting point and an error if the operation is not valid.
func (p *Point) Add(q *Point) (*Point, error) {
	// Check if the points are on the same curve
	if !p.EqualEllipticCurve(q) {
		return nil, fmt.Errorf("points are on different curves")
	}

	// Check if either of the points is the identity point
	if p.IsIdentityElement() {
		return q.Copy()
	}

	if q.IsIdentityElement() {
		return p.Copy()
	}

	// Exception when the tangent line is vertical, then return the identity point
	a := p.a
	b := p.b
	if p.Equal(q) && p.y.value == big.NewInt(0) {
		return NewPoint(nil, nil, a, b)
	}

	// Check if the points are additive inverses of each other, then return point at infinity (identity)
	if p.Equal(&Point{q.x, q.y.Negate(), a, b}) {
		return NewPoint(nil, nil, a, b)
	}

	// Calculate the sum of the points using the elliptic curve addition rules
	x1, y1 := p.x, p.y
	x2, y2 := q.x, q.y

	// Calculate geometric difference
	dy, err := y2.Subtract(y1)
	if err != nil {
		return nil, err
	}
	dx, err := x2.Subtract(x1)
	if err != nil {
		return nil, err
	}
	slope, err := dy.Divide(dx)
	if err != nil {
		return nil, err
	}

	slopeSquared, err := slope.Squared()
	if err != nil {
		return nil, err
	}

	xTotal, err := x1.Add(x2)
	if err != nil {
		return nil, err
	}

	x3, err := slopeSquared.Subtract(xTotal)
	if err != nil {
		return nil, err
	}

	dx13, err := x1.Subtract(x3)
	if err != nil {
		return nil, err
	}

	slopedx13, err := slope.Multiply(dx13)
	if err != nil {
		return nil, err
	}

	y3, err := slopedx13.Subtract(y1)
	if err != nil {
		return nil, err
	}

	return NewPoint(x3, y3, a, b)
}
