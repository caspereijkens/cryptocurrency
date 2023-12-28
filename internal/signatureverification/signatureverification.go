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
	"io"
	"math/big"

	"github.com/caspereijkens/cryptocurrency/internal/utils"
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

func (sig *Signature) Serialize() []byte {
	rSerialized := utils.SerializeInt(sig.R)
	sSerialized := utils.SerializeInt(sig.S)

	result := append([]byte{0x02, byte(len(rSerialized))}, rSerialized...)
	result = append(result, []byte{0x02, byte(len(sSerialized))}...)
	result = append(result, sSerialized...)

	return append([]byte{0x30, byte(len(result))}, result...)
}

func ParseDER(data []byte) (*Signature, error) {
	reader := bytes.NewReader(data)

	compound, err := reader.ReadByte()
	if err != nil || compound != 0x30 {
		return nil, fmt.Errorf("bad signature")
	}

	length, err := reader.ReadByte()
	if err != nil || length+2 != byte(len(data)) {
		return nil, fmt.Errorf("incorrect signature length")
	}

	r, err := parseBigInt(reader)
	if err != nil {
		return nil, err
	}

	s, err := parseBigInt(reader)
	if err != nil {
		return nil, err
	}

	if length != 6+byte(r.BitLen()/8+s.BitLen()/8) {
		return nil, fmt.Errorf("Signature too long")
	}

	return NewSignature(r, s), nil
}

func parseBigInt(reader *bytes.Reader) (*big.Int, error) {
	marker, err := reader.ReadByte()
	if err != nil || marker != 0x02 {
		return nil, fmt.Errorf("bad Signature")
	}

	valLength, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("bad Signature")
	}

	valBytes := make([]byte, valLength)
	_, err = io.ReadFull(reader, valBytes)
	if err != nil {
		return nil, fmt.Errorf("bad Signature")
	}

	intVal := new(big.Int).SetBytes(valBytes)
	return intVal, nil
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

// The Standards for Efficient Cryptography are rules for writing down ECDSA public keys.
// There are two ways to serialize elliptic curve points: compressed and uncompressed.
//
// Uncompressed
// Starts with x04, then Px[bytes], and Py[bytes].
//
// Compressed
// In elliptic curve cryptography, for each x on the curve, there are at most two possible y values because of the curve's equation.
// This is also true for finite fields.
// If a point (x, y) satisfies the curve's equation y^2 = x^3 + ax + b, then (x, -y) will work too.
// Also, in a finite field, -y % p = p-y % p. This means if (x, y) satisfies the equation, then (x, p-y) also works.
// Thus for each x, there are only two possible y points: y or p-y.
// Since p is a prime number bigger than 2 and odd, y and p-y will always be one even and one odd.
// We use this fact in the compressed SEC format. Instead of writing the whole y value, we just say if it's even or odd, and give the x value.
// So, the compressed SEC format is shorter because it turns the y value into just one byte that tells us if it's even or odd.
func (p256 *S256Point) Serialize(compressed bool) []byte {
	if compressed {
		prefix := byte(0x02)
		if new(big.Int).Mod(p256.Y.Value, big.NewInt(2)).Cmp(big.NewInt(0)) != 0 {
			prefix = 0x03
		}
		xBytes := p256.X.Value.FillBytes(make([]byte, 32))
		return append([]byte{prefix}, xBytes...)
	} else {
		prefix := []byte{byte(0x04)}
		xBytes := p256.X.Value.FillBytes(make([]byte, 32))
		yBytes := p256.Y.Value.FillBytes(make([]byte, 32))
		return append(append(prefix, xBytes...), yBytes...)
	}
}

func (p256 *S256Point) Hash160(compressed bool) []byte {
	return utils.Hash160(p256.Serialize(compressed))
}

func (p256 *S256Point) Address(compressed, testnet bool) string {
	h160 := p256.Hash160(compressed)
	if testnet {
		prefix := []byte{byte(0x6f)}
		return utils.EncodeBase58Checksum(append(prefix, h160...))
	}
	prefix := []byte{byte(0x00)}
	return utils.EncodeBase58Checksum(append(prefix, h160...))
}

func ParseSEC(sec []byte) (*S256Point, error) {
	var yField *S256FieldElement

	if len(sec) < 33 {
		return nil, fmt.Errorf("invalid SEC format")
	}

	if sec[0] == 4 {
		// Uncompressed SEC
		if len(sec) < 65 {
			return nil, fmt.Errorf("invalid uncompressed SEC format")
		}
		xField, err := NewS256FieldElement(new(big.Int).SetBytes(sec[1:33]))
		if err != nil {
			return nil, err
		}
		yField, err = NewS256FieldElement(new(big.Int).SetBytes(sec[33:65]))
		if err != nil {
			return nil, err
		}
		return NewS256Point(xField, yField)
	}

	// Compressed SEC
	x, err := NewS256FieldElement(new(big.Int).SetBytes(sec[1:]))
	if err != nil {
		return nil, err
	}

	xCubed, err := x.Exponentiate(big.NewInt(3))
	if err != nil {
		return nil, err
	}

	ySquared, err := xCubed.Add(&B.FieldElement)
	if err != nil {
		return nil, err
	}

	yEven, yOdd, err := ySquared.GetEvenOddSquareRoots()
	if err != nil {
		return nil, err
	}

	isEven := sec[0] == 2
	if isEven {
		yField, err = NewS256FieldElement(yEven)
	} else {
		yField, err = NewS256FieldElement(yOdd)
	}

	if err != nil {
		return nil, err
	}

	return NewS256Point(x, yField)
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
	k = utils.HmacSHA256(k, append(append(v, 0x00), append(secretBytes, zBytes...)...))
	v = utils.HmacSHA256(k, v)
	k = utils.HmacSHA256(k, append(append(v, 0x01), append(secretBytes, zBytes...)...))
	v = utils.HmacSHA256(k, v)

	candidate := new(big.Int)
	for {
		v = utils.HmacSHA256(k, v)
		candidate.SetBytes(v)

		if candidate.Cmp(big.NewInt(1)) >= 0 && candidate.Cmp(N) < 0 {
			return candidate
		}

		k = utils.HmacSHA256(k, append(v, 0x00))
		v = utils.HmacSHA256(k, v)
	}
}

func (e *PrivateKey) Serialize(compressed bool, testnet bool) string {
	secretBytes := e.Secret.FillBytes(make([]byte, 32))

	if compressed {
		secretBytes = append(secretBytes, byte(0x01))
	}

	var prefix []byte
	if testnet {
		prefix = []byte{0xef}
	} else {
		prefix = []byte{0x80}
	}

	payload := append(prefix, secretBytes...)

	return utils.EncodeBase58Checksum(payload)
}
