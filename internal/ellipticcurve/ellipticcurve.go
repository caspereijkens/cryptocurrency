package ellipticcurve

import (
	"fmt"
	"math/big"

	"github.com/caspereijkens/cryptocurrency/internal/finitefield"
)

// Point represents a point on the Elliptic Curve y^2 = x^3 + 7
type Point struct {
	X *finitefield.FieldElement
	Y *finitefield.FieldElement
	A *finitefield.FieldElement
	B *finitefield.FieldElement
}

func NewPoint(x, y, a, b *finitefield.FieldElement) (*Point, error) {
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
	return p.X == nil && p.Y == nil
}

func (p *Point) Equal(q *Point) bool {
	// Check that the points are on the same line
	if !p.A.Equal(q.A) || !p.B.Equal(q.B) {
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

	return p.X.Equal(q.X) && p.Y.Equal(q.Y)
}

func (p *Point) EqualEllipticCurve(q *Point) bool {
	return p.A.Equal(q.A) && p.B.Equal(q.B)
}

// String returns the string representation of a field element.
func (p *Point) String() string {
	var aVal, bVal, xVal, yVal, xPrime string

	if p == nil {
		return "Point(nil)"
	}

	if p.A != nil && p.A.Value != nil {
		aVal = p.A.Value.String()
	} else {
		aVal = "<nil>"
	}

	if p.B != nil && p.B.Value != nil {
		bVal = p.B.Value.String()
	} else {
		bVal = "<nil>"
	}

	if p.A.Prime != nil {
		xPrime = p.A.Prime.String()
	} else {
		xPrime = "<nil>"
	}

	if p.X != nil && p.X.Value != nil {
		xVal = p.X.Value.String()
	} else {
		xVal = "inf"
	}

	if p.Y != nil && p.Y.Value != nil {
		yVal = p.Y.Value.String()
	} else {
		yVal = "inf"
	}

	return fmt.Sprintf("Point_%s_%s(%s,%s) Field_%s", aVal, bVal, xVal, yVal, xPrime)
}

// Copy returns a new Point with the same values as the current Point.
func (p *Point) Copy() (*Point, error) {
	// TODO add unittest
	return NewPoint(p.X, p.Y, p.A, p.B)
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
		return NewPoint(nil, nil, p.A, p.B)
	}
	// Check if the points are additive inverses of each other, then return point at infinity (identity)
	y2_neg, err := q.Y.Negate()
	if err != nil {
		return nil, err
	}
	if p.Equal(&Point{q.X, y2_neg, p.A, p.B}) {
		return NewPoint(nil, nil, p.A, p.B)
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

	return NewPoint(x3, y3, p.A, p.B)
}

func (p *Point) calculateSlope(q *Point) (*finitefield.FieldElement, error) {
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
	return p.Equal(q) && p.Y.Value.Cmp(big.NewInt(0)) == 0
}

func (p *Point) calculateX3(q *Point, slope *finitefield.FieldElement) (*finitefield.FieldElement, error) {
	slopeSquared, err := slope.Squared()
	if err != nil {
		return nil, err
	}

	xTotal, err := p.X.Add(q.X)
	if err != nil {
		return nil, err
	}

	x3, err := slopeSquared.Subtract(xTotal)
	if err != nil {
		return nil, err
	}

	return x3, nil
}

func (p *Point) calculateY3(q *Point, x3 *finitefield.FieldElement, slope *finitefield.FieldElement) (*finitefield.FieldElement, error) {
	dx13, err := p.X.Subtract(x3)
	if err != nil {
		return nil, err
	}

	slopedx13, err := slope.Multiply(dx13)
	if err != nil {
		return nil, err
	}

	y3, err := slopedx13.Subtract(p.Y)
	if err != nil {
		return nil, err
	}

	return y3, nil
}

// Calculates dx and dy needed to compute the slope.
func (p *Point) calculatedxdy(q *Point) (*finitefield.FieldElement, *finitefield.FieldElement, error) {
	if p.Equal(q) {
		// In this case we need to compute the differential
		three, err := finitefield.NewFieldElement(big.NewInt(3), p.X.Prime)
		if err != nil {
			return nil, nil, err
		}
		dy, err := p.X.Squared()
		if err != nil {
			return nil, nil, err
		}
		dy, err = dy.Multiply(three)
		if err != nil {
			return nil, nil, err
		}
		dy, err = dy.Add(p.A)
		if err != nil {
			return nil, nil, err
		}
		dx, err := p.Y.Add(p.Y)
		if err != nil {
			return nil, nil, err
		}
		return dx, dy, nil
	}
	dy, err := q.Y.Subtract(p.Y)
	if err != nil {
		return nil, nil, err
	}

	dx, err := q.X.Subtract(p.X)
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
	result, err := NewPoint(nil, nil, p.A, p.B)
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
