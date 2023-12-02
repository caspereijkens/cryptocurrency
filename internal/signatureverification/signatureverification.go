// There are many cryptographic curves and they have different security/convenience trade-offs.
// The one that Bitcoin uses is secp256k1. It is a relatively simple curve and p is very close to 2^256.
// So most numbers under 2^256 are in the prime field.
// Any point on the curve has x and y coordinates that are expressible in 256 bits each.
// n is also very close to 2^256, so any scalar multiple can also be expressed in 256 bits.
// 2^256 is a huge number, but can still be stored in 32 bytes, so the private key can be stored easily.

package signatureverification

import (
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func NewSignature(r, s *big.Int) *Signature {
	return &Signature{R: new(big.Int).Set(r), S: new(big.Int).Set(s)}
}

func (sig *Signature) String() string {
	return fmt.Sprintf("Signature(%x,%x)", sig.R, sig.S)
}

func (p256 *S256Point) Verify(z *big.Int, sig *Signature) bool {
	// Calculate s_inv (modular inverse of s)
	sInv := new(big.Int).ModInverse(sig.S, N)
	if sInv == nil {
		return false
	}

	// Calculate u = z/s
	u := new(big.Int).Mod(new(big.Int).Mul(z, sInv), N)

	// Calculate v = r/s
	v := new(big.Int).Mod(new(big.Int).Mul(sig.R, sInv), N)

	// Calculate u*G
	uG, err := G.ScalarMultiplication(u)
	if err != nil {
		return false
	}

	// Calculate v*P
	vPoint, err := p256.ScalarMultiplication(v)
	if err != nil {
		return false
	}

	// Calculate uG + vP
	sumPoint, err := uG.Add(&vPoint.Point)
	if err != nil {
		return false
	}

	// yG + vP = R = kG
	// Check if the x-coordinate of the result matches target r
	if sumPoint.X.Value.Cmp(sig.R) != 0 {
		return false
	}

	return true
}
