package cryptography

import (
	"errors"
	"fmt"
	"math/big"
)

var (
	S256Prime = getS256Prime()
	A, _      = NewS256FieldElement(big.NewInt(0))
	B, _      = NewS256FieldElement(big.NewInt(7))
	N         = getOrder()
	G         = getS256Generator()
)

// FieldElement represents an element in a finite field.
type FieldElement struct {
	value *big.Int
	prime *big.Int
}

func getOrder() *big.Int {
	// Since the generator Point is known, the group that it generates and so its order are also known.
	// This is the hex value of this order.
	orderHex := "0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"

	// Set the big int to your large number
	order, _ := new(big.Int).SetString(orderHex, 0) // [2:] to remove the '0x' prefix

	return order

}

func getS256Generator() *S256Point {
	// https://crypto.stackexchange.com/questions/60420/what-does-the-special-form-of-the-base-point-of-secp256k1-allow
	xHex := "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	yHex := "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"

	x, _ := new(big.Int).SetString(xHex, 0) // [2:] to remove the '0x' prefix
	y, _ := new(big.Int).SetString(yHex, 0) // [2:] to remove the '0x' prefix

	xF, _ := NewS256FieldElement(x)
	yF, _ := NewS256FieldElement(y)

	generator, _ := NewS256Point(xF, yF)
	return generator
}

func getS256Prime() *big.Int {
	two := big.NewInt(2)
	p := new(big.Int)

	// Calculate 2^256
	twoToThe256 := new(big.Int).Exp(two, big.NewInt(256), nil)

	// Calculate 2^32
	twoToThe32 := new(big.Int).Exp(two, big.NewInt(32), nil)

	// Subtract 2^32 and 977 from 2^256
	p.Sub(twoToThe256, twoToThe32)
	p.Sub(p, big.NewInt(977))
	return p
}

// NewFieldElement creates a new FieldElement with the given value and prime.
func NewFieldElement(value, prime *big.Int) (*FieldElement, error) {
	if value == nil {
		return nil, nil
	}
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
	result := new(big.Int).Mod(new(big.Int).Add(a.value, b.value), a.prime)
	return NewFieldElement(result, a.prime)
}

// Subtract subtracts two field elements and returns a new field element.
func (a *FieldElement) Subtract(b *FieldElement) (*FieldElement, error) {
	if a.prime.Cmp(b.prime) != 0 {
		return nil, errors.New("field elements are from different fields")
	}
	result := new(big.Int).Sub(a.value, b.value)
	if result.Sign() < 0 {
		result.Add(result, a.prime)
	}
	return NewFieldElement(result, a.prime)
}

// Multiply multiplies two field elements and returns a new field element.
func (a *FieldElement) Multiply(b *FieldElement) (*FieldElement, error) {
	if a.prime.Cmp(b.prime) != 0 {
		return nil, errors.New("field elements are from different fields")
	}
	result := new(big.Int).Mul(a.value, b.value)
	result.Mod(result, a.prime)
	return NewFieldElement(result.Mod(result, a.prime), a.prime)
}

// Exponentiate computes the exponentiation of a field element to a given power.
func (a *FieldElement) Exponentiate(power *big.Int) (*FieldElement, error) {
	result := new(big.Int).Exp(a.value, power, a.prime)
	return NewFieldElement(result.Mod(result, a.prime), a.prime)
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
	return a.value.Cmp(b.value) == 0 && a.prime.Cmp(b.prime) == 0
}

// Negate returns a new FieldElement with the negated value of the current FieldElement.
func (a *FieldElement) Negate() (*FieldElement, error) {
	// Calculate the negated value as (prime - value) % prime
	negatedValue := new(big.Int).Sub(a.prime, a.value)
	return NewFieldElement(negatedValue.Mod(negatedValue, a.prime), a.prime)
}

// String returns the string representation of a field element.
func (a *FieldElement) String() string {
	if a.prime.Cmp(S256Prime) == 0 {
		return fmt.Sprintf("%064x", a.value)
	}
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
	xCubed, err := x.Cubed()
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

	ySquared, err := y.Squared()
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

// String returns the string representation of a field element.
func (p *Point) String() string {
	var aVal, bVal, xVal, yVal, xPrime string

	if p == nil {
		return "Point(nil)"
	}

	if p.a != nil && p.a.value != nil {
		aVal = p.a.value.String()
	} else {
		aVal = "<nil>"
	}

	if p.b != nil && p.b.value != nil {
		bVal = p.b.value.String()
	} else {
		bVal = "<nil>"
	}

	if p.a.prime != nil {
		xPrime = p.a.prime.String()
	} else {
		xPrime = "<nil>"
	}

	if p.x != nil && p.x.value != nil {
		xVal = p.x.value.String()
	} else {
		xVal = "inf"
	}

	if p.y != nil && p.y.value != nil {
		yVal = p.y.value.String()
	} else {
		yVal = "inf"
	}

	return fmt.Sprintf("Point_%s_%s(%s,%s) Field_%s", aVal, bVal, xVal, yVal, xPrime)
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

	// Handle special cases
	// Exception when the tangent line is vertical, then return the identity point
	if p.Equal(q) && p.isVerticalTangent(q) {
		return NewPoint(nil, nil, p.a, p.b)
	}
	// Check if the points are additive inverses of each other, then return point at infinity (identity)
	y2_neg, err := q.y.Negate()
	if err != nil {
		return nil, err
	}
	if p.Equal(&Point{q.x, y2_neg, p.a, p.b}) {
		return NewPoint(nil, nil, p.a, p.b)
	}

	// Calculate the sum of the points using the elliptic curve addition rules
	slope, err := p.calculateSlope(q)
	if err != nil {
		return nil, err
	}

	x3, err := p.calculateX3(q, slope)
	if err != nil {
		return nil, err
	}

	y3, err := p.calculateY3(q, x3, slope)
	if err != nil {
		return nil, err
	}

	return NewPoint(x3, y3, p.a, p.b)
}

func (p *Point) calculateSlope(q *Point) (*FieldElement, error) {
	dx, dy, err := p.calculatedxdy(q)
	if err != nil {
		return nil, err
	}
	slope, err := dy.Divide(dx)
	if err != nil {
		return nil, err
	}
	return slope, nil
}

func (p *Point) isVerticalTangent(q *Point) bool {
	return p.Equal(q) && p.y.value.Cmp(big.NewInt(0)) == 0
}

func (p *Point) calculateX3(q *Point, slope *FieldElement) (*FieldElement, error) {
	slopeSquared, err := slope.Squared()
	if err != nil {
		return nil, err
	}

	xTotal, err := p.x.Add(q.x)
	if err != nil {
		return nil, err
	}

	x3, err := slopeSquared.Subtract(xTotal)
	if err != nil {
		return nil, err
	}

	return x3, nil
}

func (p *Point) calculateY3(q *Point, x3 *FieldElement, slope *FieldElement) (*FieldElement, error) {
	dx13, err := p.x.Subtract(x3)
	if err != nil {
		return nil, err
	}

	slopedx13, err := slope.Multiply(dx13)
	if err != nil {
		return nil, err
	}

	y3, err := slopedx13.Subtract(p.y)
	if err != nil {
		return nil, err
	}

	return y3, nil
}

// Calculates dx and dy needed to compute the slope.
func (p *Point) calculatedxdy(q *Point) (*FieldElement, *FieldElement, error) {
	if p.Equal(q) {
		// In this case we need to compute the differential
		three, err := NewFieldElement(big.NewInt(3), p.x.prime)
		if err != nil {
			return nil, nil, err
		}
		dy, err := p.x.Squared()
		if err != nil {
			return nil, nil, err
		}
		dy, err = dy.Multiply(three)
		if err != nil {
			return nil, nil, err
		}
		dy, err = dy.Add(p.a)
		if err != nil {
			return nil, nil, err
		}
		dx, err := p.y.Add(p.y)
		if err != nil {
			return nil, nil, err
		}
		return dx, dy, nil
	}
	dy, err := q.y.Subtract(p.y)
	if err != nil {
		return nil, nil, err
	}

	dx, err := q.x.Subtract(p.x)
	if err != nil {
		return nil, nil, err
	}
	return dx, dy, nil
}

// ScalarMult performs scalar multiplication of a point on an elliptic curve.
func (p *Point) ScalarMultiplication(coefficient *big.Int) (*Point, error) {
	if coefficient.Sign() == -1 {
		return nil, fmt.Errorf("coefficient must be positive")
	}
	// We start the result at the identity element
	result, err := NewPoint(nil, nil, p.a, p.b)
	if err != nil {
		return nil, err
	}
	// current represents the point at the current bit.
	current, err := p.Copy()
	if err != nil {
		return nil, err
	}
	// Binary expansion, allows to do multiplication in log_2(n) loops
	for coef := coefficient; coef.Cmp(big.NewInt(0)) > 0; coef.Rsh(coef, 1) {
		// Check if the rightmost bit is a 1.
		if coef.Bit(0) == 1 {
			// Add the value of the current bit
			result, err = result.Add(current)
			if err != nil {
				return nil, err
			}
		}
		// In effect, this doubles current
		// The first time through the loop it represents  1 x p
		// The second time through the loop it represents 2 x p
		// The third time through the loop it represents  4 x p
		current, err = current.Add(current)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// S256FieldElement represents a field element with a fixed prime for secp256k1.
type S256FieldElement struct {
	FieldElement
}

// NewS256FieldElement creates a new Secp256k1FieldElement with the fixed prime.
func NewS256FieldElement(value *big.Int) (*S256FieldElement, error) {
	f, err := NewFieldElement(value, S256Prime)
	if err != nil {
		return nil, err
	}
	return &S256FieldElement{*f}, nil
}

type S256Point struct {
	Point
}

func NewS256Point(x, y *S256FieldElement) (*S256Point, error) {
	p, err := NewPoint(&x.FieldElement, &y.FieldElement, &A.FieldElement, &B.FieldElement)
	if err != nil {
		return nil, err
	}
	return &S256Point{*p}, nil
}

func (p256 *S256Point) ScalarMultiplication(coefficient *big.Int) (*S256Point, error) {
	// We mod by N because that is the order of the Group generated by this specific point.
	// In other words, every n times we cycle back to the identiy.
	p, err := p256.Point.ScalarMultiplication(new(big.Int).Mod(coefficient, N))
	if err != nil {
		return nil, err
	}
	return &S256Point{*p}, nil
}
