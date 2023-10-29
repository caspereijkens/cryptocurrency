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

// Point represents a point on the Elliptic Curve y^2 = x^3 + 7
type Point struct {
	x *big.Int
	y *big.Int
	a *big.Int
	b *big.Int
}

func NewPoint(x, y, a, b *big.Int) (*Point, error) {
	// Check if this is the point at infinity
	if x == nil && y == nil {
		return &Point{nil, nil, a, b}, nil
	}
	// Check if the point (x, y) is on the elliptic curve y^2 = x^3 + ax + b
	xCubed := new(big.Int).Exp(x, big.NewInt(3), nil)
	ax := new(big.Int).Mul(a, x)
	rightHandSide := new(big.Int).Add(xCubed, ax)
	rightHandSide.Add(rightHandSide, b)

	ySquared := new(big.Int).Mul(y, y)

	if ySquared.Cmp(rightHandSide) != 0 {
		return nil, fmt.Errorf("Point ") //(%s, %s) does not exist on elliptic curve y^2 = x^3 + %s x + %s", x.String(), y.String(), a.String(), b.String()
	}

	return &Point{x, y, a, b}, nil
}

func (p *Point) Equal(q *Point) bool {
	return p.x.Cmp(q.x) == 0 &&
		p.y.Cmp(q.y) == 0 &&
		p.a.Cmp(q.a) == 0 &&
		p.b.Cmp(q.b) == 0
}

func (p *Point) NotEqual(q *Point) bool {
	return !p.Equal(q)
}

// String returns the string representation of a field element.
func (p *Point) String() string {
	return fmt.Sprintf("Point_%s_%s(%s,%s)", p.a.String(), p.b.String(), p.x.String(), p.y.String())
}

// Add performs the addition of two elliptic curve points (p and q).
// It returns the resulting point and an error if the operation is not valid.
func (p *Point) Add(q *Point) (*Point, error) {
	// Check if the points are on the same curve
	a, b := p.a, p.b
	if a.Cmp(q.a) != 0 || b.Cmp(q.b) != 0 {
		return nil, errors.New("points are on different curves")
	}

	// Check if either of the points is the point at infinity (identity)
	if p.x == nil && p.y == nil {
		return q, nil
	}
	if q.x == nil && q.y == nil {
		return p, nil
	}

	// Exception when the tangent line is vertical
	if p.Equal(q) && p.y.Cmp(big.NewInt(0)) == 0 {
		return NewPoint(nil, nil, a, b)
	}

	// Check if the points are additive inverses of each other
	if p.Equal(&Point{q.x, new(big.Int).Neg(q.y), a, b}) {
		return NewPoint(nil, nil, a, b)
	}

	// Calculate the sum of the points using the elliptic curve addition rules
	x1, y1 := p.x, p.y
	x2, y2 := q.x, q.y

	lambda := new(big.Int)
	lambda.Sub(y2, y1)
	lambda.Div(lambda, new(big.Int).Sub(x2, x1))

	x3 := new(big.Int)
	x3.Mul(lambda, lambda)
	x3.Sub(x3, x1)
	x3.Sub(x3, x2)

	y3 := new(big.Int)
	y3.Sub(x1, x3)
	y3.Mul(y3, lambda)
	y3.Sub(y3, y1)

	return NewPoint(x3, y3, a, b)
}
