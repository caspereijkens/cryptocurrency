// There are many cryptographic curves and they have different security/convenience trade-offs.
// The one that Bitcoin uses is secp256k1. It is a relatively simple curve and p is very close to 2^256.
// So most numbers under 2^256 are in the prime field.
// Any point on the curve has x and y coordinates that are expressible in 256 bits each.
// n is also very close to 2^256, so any scalar multiple can also be expressed in 256 bits.
// 2^256 is a huge number, but can still be stored in 32 bytes, so the private key can be stored easily.

package signatureverification

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/caspereijkens/cryptocurrency/internal/util"
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

// The verification procedure is as follows:
// 1. We are given (r,s) as the signature, z as the hash of the thing being signed, and P as the public key (or public point) of the signer;
// 2. We calcualte u = z/s, v = r/s;
// 3. We calculate uG + vP = R;
// 4. If R's x-coordinate equals r, the signature is valid;
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

type PrivateKey struct {
	Secret *big.Int
	Point  *S256Point
}

func NewPrivateKey(secret *big.Int) (*PrivateKey, error) {
	point, err := G.ScalarMultiplication(secret)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{secret, point}, nil
}

// The signing procedure is as follows:
// 1. We are given signature hash z and and know private key e such that eG = P;
// 2. Choose a random k;
// 3. Calculate R = kG. r is the x-coordinate of R;
// 4. Calculate s = (z + re)/k;
// 5. Signature is (r,s);
func (e *PrivateKey) Sign(z *big.Int) (*Signature, error) {

	if z == nil {
		return nil, fmt.Errorf("one or more signature inputs were invalid")
	}

	k := e.GetDeterministicK(z)

	// Calculate the target R
	R, err := G.ScalarMultiplication(k)

	if err != nil {
		return nil, err
	}

	// Calculate r, the x-value of target R
	r := R.X.Value

	// Calculate r * e
	re := new(big.Int).Mul(r, e.Secret)

	// Calculate re + z
	rePlusZ := new(big.Int).Add(re, z)

	// Calculate (re + z) * kInv
	kInv := new(big.Int).ModInverse(k, N)
	product := new(big.Int).Mul(rePlusZ, kInv)

	// Modulo with N to get the final result
	s := new(big.Int).Mod(product, N)

	// P, err := G.ScalarMultiplication(e)

	if err != nil {
		return nil, err
	}

	return NewSignature(r, s), nil
}

// Deterministic k generation standard that uses the secret and z to create a unique, deterministic k every time.
// Specification is in RFC 6979
// If our secret is e and we are reusing k to sign z1 and z2:
// kG = (r,y)
// s1 = (z1 + re)/k, s2 = (z2 + re)/k
// s1/s2 = (z1 + re) / (z2 + re)
// s1(z2 + re) = s2(z1 +re)
// s1z2 + s1re = s2z1 + s2re
// s1re - s2re = s2z1 - s1z2
// e = (s2z1 - s1z2) / (s1r - s2r)
func (e *PrivateKey) GetDeterministicK(z *big.Int) *big.Int {
	// Ensure z is within the correct range
	if z.Cmp(N) > 0 {
		z.Sub(z, N)
	}

	k := make([]byte, 32)
	v := bytes.Repeat([]byte{0x01}, 32)
	zBytes := z.FillBytes(make([]byte, 32))
	secretBytes := e.Secret.FillBytes(make([]byte, 32))

	// Updating k and v
	k = util.HmacSHA256(k, append(append(v, 0x00), append(secretBytes, zBytes...)...))
	v = util.HmacSHA256(k, v)
	k = util.HmacSHA256(k, append(append(v, 0x01), append(secretBytes, zBytes...)...))
	v = util.HmacSHA256(k, v)

	candidate := new(big.Int)
	for {
		v = util.HmacSHA256(k, v)
		candidate.SetBytes(v)

		if candidate.Cmp(big.NewInt(1)) >= 0 && candidate.Cmp(N) < 0 {
			return candidate
		}

		k = util.HmacSHA256(k, append(v, 0x00))
		v = util.HmacSHA256(k, v)
	}
}
