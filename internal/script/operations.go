package script

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/caspereijkens/cryptocurrency/internal/signatureverification"
	"github.com/caspereijkens/cryptocurrency/internal/utils"
)

type Stack [][]byte

func encodeNum(num int) []byte {
	if num == 0 {
		return []byte{}
	}

	absNum := num
	if absNum < 0 {
		absNum = -absNum
	}

	var result []byte
	for absNum > 0 {
		result = append(result, byte(absNum&0xff))
		absNum >>= 8
	}

	if result[len(result)-1]&0x80 != 0 {
		if num < 0 {
			result = append(result, 0x80)
		} else {
			result = append(result, 0)
		}
	} else if num < 0 {
		result[len(result)-1] |= 0x80
	}

	return result
}

func decodeNum(element []byte) int {
	if len(element) == 0 {
		return 0
	}

	var bigEndian []byte
	for i := len(element) - 1; i >= 0; i-- {
		bigEndian = append(bigEndian, element[i])
	}

	var negative bool
	var result int

	if bigEndian[0]&0x80 != 0 {
		negative = true
		result = int(bigEndian[0] & 0x7f)
	} else {
		negative = false
		result = int(bigEndian[0])
	}

	for _, c := range bigEndian[1:] {
		result <<= 8
		result += int(c)
	}

	if negative {
		return -result
	}

	return result
}

func op0(stack *Stack) (bool, error) {
	stack.push(encodeNum(0))
	return true, nil
}

func op1Negate(stack *Stack) (bool, error) {
	stack.push(encodeNum(-1))
	return true, nil
}

func op1(stack *Stack) (bool, error) {
	stack.push(encodeNum(1))
	return true, nil
}

func op2(stack *Stack) (bool, error) {
	stack.push(encodeNum(2))
	return true, nil
}

func op3(stack *Stack) (bool, error) {
	stack.push(encodeNum(3))
	return true, nil
}

func op4(stack *Stack) (bool, error) {
	stack.push(encodeNum(4))
	return true, nil
}

func op5(stack *Stack) (bool, error) {
	stack.push(encodeNum(5))
	return true, nil
}

func op6(stack *Stack) (bool, error) {
	stack.push(encodeNum(6))
	return true, nil
}

func op7(stack *Stack) (bool, error) {
	stack.push(encodeNum(7))
	return true, nil
}

func op8(stack *Stack) (bool, error) {
	stack.push(encodeNum(8))
	return true, nil
}

func op9(stack *Stack) (bool, error) {
	stack.push(encodeNum(9))
	return true, nil
}

func op10(stack *Stack) (bool, error) {
	stack.push(encodeNum(10))
	return true, nil
}

func op11(stack *Stack) (bool, error) {
	stack.push(encodeNum(11))
	return true, nil
}

func op12(stack *Stack) (bool, error) {
	stack.push(encodeNum(12))
	return true, nil
}

func op13(stack *Stack) (bool, error) {
	stack.push(encodeNum(13))
	return true, nil
}

func op14(stack *Stack) (bool, error) {
	stack.push(encodeNum(14))
	return true, nil
}

func op15(stack *Stack) (bool, error) {
	stack.push(encodeNum(15))
	return true, nil
}

func op16(stack *Stack) (bool, error) {
	stack.push(encodeNum(16))
	return true, nil
}

func opNop(stack *Stack) (bool, error) {
	return true, nil
}

func opIf(stack, items *Stack) (bool, error) {
	if len(*stack) < 1 {
		return false, fmt.Errorf("stack is empty")
	}

	// go through and re-make the items array based on the top stack element
	trueItems, falseItems := new(Stack), new(Stack)
	var found bool
	currentArray := trueItems
	numEndifsNeeded := 1

	for len(*items) > 0 {
		item, err := items.pop(0)
		if err != nil {
			return false, err
		}

		if bytes.Equal(item, encodeNum(99)) || bytes.Equal(item, encodeNum(100)) {
			// nested if, we have to go another endif
			numEndifsNeeded++
			*currentArray = append(*currentArray, item)
		} else if numEndifsNeeded == 1 && bytes.Equal(item, encodeNum(103)) {
			currentArray = falseItems
		} else if bytes.Equal(item, encodeNum(104)) {
			if numEndifsNeeded == 1 {
				found = true
				break
			} else {
				numEndifsNeeded--
				*currentArray = append(*currentArray, item)
			}
		} else {
			*currentArray = append(*currentArray, item)
		}
	}

	if !found {
		return false, nil
	}

	element, _ := stack.pop(-1)
	if bytes.Equal(element, encodeNum(0)) {
		*items = append(*falseItems, *items...)
	} else {
		*items = append(*trueItems, *items...)
	}

	return true, nil
}

func opNotIf(stack, items *Stack) (bool, error) {
	if len(*stack) < 1 {
		return false, fmt.Errorf("stack is empty")
	}

	// go through and re-make the items array based on the top stack element
	trueItems, falseItems := new(Stack), new(Stack)
	var found bool
	currentArray := trueItems
	numEndifsNeeded := 1

	for len(*items) > 0 {
		item, err := items.pop(0)
		if err != nil {
			return false, err
		}

		if bytes.Equal(item, encodeNum(99)) || bytes.Equal(item, encodeNum(100)) {
			// nested if, we have to go another endif
			numEndifsNeeded++
			*currentArray = append(*currentArray, item)
		} else if numEndifsNeeded == 1 && bytes.Equal(item, encodeNum(103)) {
			currentArray = falseItems
		} else if bytes.Equal(item, encodeNum(104)) {
			if numEndifsNeeded == 1 {
				found = true
				break
			} else {
				numEndifsNeeded--
				*currentArray = append(*currentArray, item)
			}
		} else {
			*currentArray = append(*currentArray, item)
		}
	}

	if !found {
		return false, nil
	}

	element, _ := stack.pop(-1)
	if bytes.Equal(element, encodeNum(0)) {
		*items = append(*trueItems, *items...)
	} else {
		*items = append(*falseItems, *items...)
	}

	return true, nil
}

func opVerify(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	return (decodeNum(element) != 0), nil
}

func opReturn(stack *Stack) (bool, error) {
	return false, nil
}

func opToAltStack(stack, altStack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	*altStack = append(*altStack, element)

	return true, nil
}

func opFromAltStack(stack, altStack *Stack) (bool, error) {
	element, err := altStack.pop(-1)

	if err != nil {
		return false, err
	}

	stack.push(element)

	return true, nil
}

func op2Drop(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	*stack = (*stack)[:len(*stack)-2]
	return true, nil
}

func op2Dup(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	*stack = append(*stack, (*stack)[len(*stack)-2:]...)
	return true, nil
}

func op3Dup(stack *Stack) (bool, error) {
	if len(*stack) < 3 {
		return false, fmt.Errorf("not enough elements in stack: %d < 3", len(*stack))
	}

	*stack = append(*stack, (*stack)[len(*stack)-3:]...)
	return true, nil
}

func op2Over(stack *Stack) (bool, error) {
	if len(*stack) < 4 {
		return false, fmt.Errorf("not enough elements in stack: %d < 4", len(*stack))
	}

	*stack = append(*stack, (*stack)[len(*stack)-4:len(*stack)-2]...)
	return true, nil
}

func op2Rot(stack *Stack) (bool, error) {
	if len(*stack) < 6 {
		return false, fmt.Errorf("not enough elements in stack: %d < 6", len(*stack))
	}

	*stack = append(*stack, (*stack)[len(*stack)-6:len(*stack)-4]...)
	return true, nil
}

func op2Swap(stack *Stack) (bool, error) {
	if len(*stack) < 4 {
		return false, fmt.Errorf("not enough elements in stack: %d < 4", len(*stack))
	}

	lastFour := (*stack)[len(*stack)-4:]
	(*stack)[len(*stack)-4] = lastFour[2]
	(*stack)[len(*stack)-3] = lastFour[3]
	(*stack)[len(*stack)-2] = lastFour[0]
	(*stack)[len(*stack)-1] = lastFour[1]

	return true, nil
}

func opIfDup(stack *Stack) (bool, error) {
	if len(*stack) < 1 {
		return false, fmt.Errorf("stack is empty")
	}

	element := (*stack)[len(*stack)-1]

	if decodeNum(element) != 0 {
		stack.push(element)
	}

	return true, nil
}

func opDepth(stack *Stack) (bool, error) {
	stack.push(encodeNum(len(*stack)))
	return true, nil
}

func opDrop(stack *Stack) (bool, error) {
	_, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	return true, nil
}

func opDup(stack *Stack) (bool, error) {
	if len(*stack) < 1 {
		return false, fmt.Errorf("stack is empty")
	}

	element := (*stack)[len(*stack)-1]

	stack.push(element)

	return true, nil
}

func opNip(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	*stack = append((*stack)[:len(*stack)-2], (*stack)[len(*stack)-1])
	return true, nil
}

func opOver(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	stack.push((*stack)[len(*stack)-2])

	return true, nil
}

func opPick(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	n := decodeNum(element)

	if len(*stack) < n+1 {
		return false, fmt.Errorf("not enough elements in stack: %d < %d", len(*stack), n+1)
	}

	stack.push((*stack)[len(*stack)-n-1])

	return true, nil
}

func opRoll(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	n := decodeNum(element)

	if len(*stack) < n+1 {
		return false, fmt.Errorf("not enough elements in stack: %d < %d", len(*stack), n+1)
	}

	if n > 0 {
		rolled := (*stack)[len(*stack)-n-1]
		*stack = append((*stack)[:len(*stack)-n-1], (*stack)[len(*stack)-n:]...)
		stack.push(rolled)
	}

	return true, nil
}

func opRot(stack *Stack) (bool, error) {
	element, err := stack.pop(-3)

	if err != nil {
		return false, err
	}

	stack.push(element)
	return true, nil
}

func opSwap(stack *Stack) (bool, error) {
	element, err := stack.pop(-2)

	if err != nil {
		return false, err
	}

	stack.push(element)
	return true, nil
}

func opTuck(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	err := stack.insert(-2, (*stack)[len(*stack)-1])

	if err != nil {
		return false, err
	}

	return true, nil
}

// pushes the size of the last item on the stack
func opSize(stack *Stack) (bool, error) {
	if len(*stack) < 1 {
		return false, fmt.Errorf("stack is empty")
	}

	element := (*stack)[len(*stack)-1]
	stack.push(encodeNum(len(element)))
	return true, nil
}

func opEqual(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(element1, element2) {
		stack.push(encodeNum(0))
		return true, nil
	}

	stack.push(encodeNum(1))
	return true, nil
}

func opEqualVerify(stack *Stack) (bool, error) {
	resultEqual, err := opEqual(stack)

	if err != nil || !resultEqual {
		return false, err
	}

	return opVerify(stack)
}

func op1Add(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	stack.push(encodeNum(decodeNum(element) + 1))
	return true, nil
}

func op1Sub(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	stack.push(encodeNum(decodeNum(element) - 1))
	return true, nil
}

func opNegate(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	stack.push(encodeNum(-decodeNum(element)))
	return true, nil
}

func opAbs(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	if decodeNum(element) < 0 {
		stack.push(encodeNum(-decodeNum(element)))
		return true, nil
	}

	stack.push(encodeNum(decodeNum(element)))
	return true, nil
}

func opNot(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	var notElement int

	if decodeNum(element) == 0 {
		notElement = 1
	}

	stack.push(encodeNum(notElement))
	return true, nil
}

func op0NotEqual(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	var notElement int

	if decodeNum(element) != 0 {
		notElement = 1
	}

	stack.push(encodeNum(notElement))
	return true, nil
}

func opAdd(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	stack.push(encodeNum(decodeNum(element1) + decodeNum(element2)))
	return true, nil
}

func opSub(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	stack.push(encodeNum(decodeNum(element2) - decodeNum(element1)))
	return true, nil
}

func opMul(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	stack.push(encodeNum(decodeNum(element2) * decodeNum(element1)))
	return true, nil
}

func opBoolAnd(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	if decodeNum(element1) != 0 && decodeNum(element2) != 0 {
		stack.push(encodeNum(1))
		return true, nil
	}

	stack.push(encodeNum(0))
	return true, nil
}

func opBoolOr(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	if decodeNum(element1) != 0 || decodeNum(element2) != 0 {
		stack.push(encodeNum(1))
		return true, nil
	}

	stack.push(encodeNum(0))
	return true, nil
}

func opNumEqual(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	if decodeNum(element1) != decodeNum(element2) {
		stack.push(encodeNum(0))
		return true, nil
	}

	stack.push(encodeNum(1))
	return true, nil
}

func opNumEqualVerify(stack *Stack) (bool, error) {
	resultNumEqual, err := opNumEqual(stack)

	if err != nil || !resultNumEqual {
		return false, err
	}

	return opVerify(stack)
}

func opNumNotEqual(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	if decodeNum(element1) == decodeNum(element2) {
		stack.push(encodeNum(0))
		return true, nil
	}

	stack.push(encodeNum(1))
	return true, nil
}

func opLessThan(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	if decodeNum(element2) >= decodeNum(element1) {
		stack.push(encodeNum(0))
		return true, nil
	}

	stack.push(encodeNum(1))
	return true, nil
}

func opGreaterThan(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	if decodeNum(element2) <= decodeNum(element1) {
		stack.push(encodeNum(0))
		return true, nil
	}

	stack.push(encodeNum(1))
	return true, nil
}

func opLessThanOrEqual(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	if decodeNum(element2) > decodeNum(element1) {
		stack.push(encodeNum(0))
		return true, nil
	}

	stack.push(encodeNum(1))
	return true, nil
}

func opGreaterThanOrEqual(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	if decodeNum(element2) < decodeNum(element1) {
		stack.push(encodeNum(0))
		return true, nil
	}

	stack.push(encodeNum(1))
	return true, nil
}

func opMin(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	stack.push(encodeNum(min(decodeNum(element1), decodeNum(element2))))
	return true, nil
}

func opMax(stack *Stack) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	element1, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element2, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	stack.push(encodeNum(max(decodeNum(element1), decodeNum(element2))))
	return true, nil
}

func opWithin(stack *Stack) (bool, error) {
	if len(*stack) < 3 {
		return false, fmt.Errorf("not enough elements in stack: %d < 3", len(*stack))
	}

	maximum, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	minimum, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	element, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	var within int

	if decodeNum(element) >= decodeNum(minimum) && decodeNum(element) < decodeNum(maximum) {
		within = 1
	}

	stack.push(encodeNum(within))
	return true, nil
}

func opRipemd160(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	stack.push(utils.Ripemd160Hash(element))
	return true, nil
}

func opSha1(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	stack.push(utils.Sha1Hash(element))
	return true, nil
}

func opSha256(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	stack.push(utils.Sha256Hash(element))
	return true, nil
}

func opHash160(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	stack.push(utils.Hash160(element))
	return true, nil
}

func opHash256(stack *Stack) (bool, error) {
	element, err := stack.pop(-1)

	if err != nil {
		return false, err
	}

	stack.push(utils.Hash256(element))
	return true, nil
}

func opCheckSig(stack *Stack, z *big.Int) (bool, error) {
	if len(*stack) < 2 {
		return false, fmt.Errorf("not enough elements in stack: %d < 2", len(*stack))
	}

	secPubkey, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	derSignatureBytes, err := stack.pop(-1)
	if err != nil {
		return false, err
	}

	// take off the last byte of the signature as that"s the hash type
	derSignature, err := signatureverification.ParseDER(derSignatureBytes[:len(derSignatureBytes)-1])
	if err != nil {
		return false, err
	}

	point, err := signatureverification.ParseSEC(secPubkey)
	if err != nil {
		return false, err
	}

	if !point.Verify(z, derSignature) {
		op0(stack)
		return false, fmt.Errorf("signature validation failed")
	}

	op1(stack)
	return true, nil
}

func opCheckSigVerify(stack *Stack, z *big.Int) (bool, error) {
	resultChecksig, err := opCheckSig(stack, z)

	if err != nil || !resultChecksig {
		return false, err
	}

	return opVerify(stack)
}

// TODO: opCheckmultisig, opCheckmultisigVerify

func opCheckLockTimeVerify(stack *Stack, locktime, sequence int) (bool, error) {
	if sequence == 0xffffffff {
		return false, fmt.Errorf("invalid sequence value")
	}

	if len(*stack) < 1 {
		return false, fmt.Errorf("stack is empty")
	}

	element := decodeNum((*stack)[len(*stack)-1])
	if element < 0 {
		return false, fmt.Errorf("negative element in stack")
	}

	if element < 500000000 && locktime > 500000000 {
		return false, fmt.Errorf("locktime exceeds 500000000 for element less than 500000000")
	}

	if locktime < element {
		return false, fmt.Errorf("locktime is less than element in stack")
	}

	return true, nil
}

func opCheckSequenceVerify(stack *Stack, version, sequence int) (bool, error) {
	if sequence&(1<<31) == (1 << 31) {
		return false, fmt.Errorf("invalid sequence value")
	}

	if len(*stack) < 1 {
		return false, fmt.Errorf("stack is empty")
	}

	element := decodeNum((*stack)[len(*stack)-1])
	if element < 0 {
		return false, fmt.Errorf("negative element in stack")
	}

	if element&(1<<31) == (1 << 31) {
		if version < 2 {
			return false, fmt.Errorf("version is less than 2 for sequence with sign bit set")
		}

		if element&(1<<22) != sequence&(1<<22) {
			return false, fmt.Errorf("mismatch in bits 22-31 between element and sequence")
		}

		if element&0xffff > sequence&0xffff {
			return false, fmt.Errorf("sequence value is less than element in stack")
		}
	}

	return true, nil
}

func (s *Stack) push(value []byte) {
	*s = append(*s, value)
}

func (stack *Stack) pop(index int) ([]byte, error) {
	if len(*stack) < 1 {
		return nil, fmt.Errorf("stack is empty")
	}

	if index < 0 {
		index = len(*stack) + index
	}

	if index < 0 || index >= len(*stack) {
		return nil, fmt.Errorf("index out of bounds")
	}

	element := (*stack)[index]
	*stack = append((*stack)[:index], (*stack)[index+1:]...)

	return element, nil
}

func (stack *Stack) insert(index int, element []byte) error {
	if index < 0 {
		index = len(*stack) + index + 1
	}

	if index < 0 || index > len(*stack) {
		return fmt.Errorf("index out of bounds")
	}

	stack.push(nil) // Ensure enough capacity for the new element
	copy((*stack)[index+1:], (*stack)[index:])
	(*stack)[index] = element

	return nil
}

// OpCodesFunctions is a map of opcode values to their corresponding functions
var OpCodesFunctions = map[int]interface{}{
	0:   op0,
	79:  op1Negate,
	81:  op1,
	82:  op2,
	83:  op3,
	84:  op4,
	85:  op5,
	86:  op6,
	87:  op7,
	88:  op8,
	89:  op9,
	90:  op10,
	91:  op11,
	92:  op12,
	93:  op13,
	94:  op14,
	95:  op15,
	96:  op16,
	97:  opNop,
	99:  opIf,
	100: opNotIf,
	105: opVerify,
	106: opReturn,
	107: opToAltStack,
	108: opFromAltStack,
	109: op2Drop,
	110: op2Dup,
	111: op3Dup,
	112: op2Over,
	113: op2Rot,
	114: op2Swap,
	115: opIfDup,
	116: opDepth,
	117: opDrop,
	118: opDup,
	119: opNip,
	120: opOver,
	121: opPick,
	122: opRoll,
	123: opRot,
	124: opSwap,
	125: opTuck,
	130: opSize,
	135: opEqual,
	136: opEqualVerify,
	139: op1Add,
	140: op1Sub,
	143: opNegate,
	144: opAbs,
	145: opNot,
	146: op0NotEqual,
	147: opAdd,
	148: opSub,
	149: opMul,
	154: opBoolAnd,
	155: opBoolOr,
	156: opNumEqual,
	157: opNumEqualVerify,
	158: opNumNotEqual,
	159: opLessThan,
	160: opGreaterThan,
	161: opLessThanOrEqual,
	162: opGreaterThanOrEqual,
	163: opMin,
	164: opMax,
	165: opWithin,
	166: opRipemd160,
	167: opSha1,
	168: opSha256,
	169: opHash160,
	170: opHash256,
	172: opCheckSig,
	173: opCheckSigVerify,
	// 174: opCheckMultiSig,
	// 175: opCheckMultiSigVerify,
	176: opNop,
	177: opCheckLockTimeVerify,
	178: opCheckSequenceVerify,
	179: opNop,
	180: opNop,
	181: opNop,
	182: opNop,
	183: opNop,
	184: opNop,
	185: opNop,
}

var OP_CODE_NAMES = map[int]string{
	0:   "OP_0",
	76:  "OP_PUSHDATA1",
	77:  "OP_PUSHDATA2",
	78:  "OP_PUSHDATA4",
	79:  "OP_1NEGATE",
	81:  "OP_1",
	82:  "OP_2",
	83:  "OP_3",
	84:  "OP_4",
	85:  "OP_5",
	86:  "OP_6",
	87:  "OP_7",
	88:  "OP_8",
	89:  "OP_9",
	90:  "OP_10",
	91:  "OP_11",
	92:  "OP_12",
	93:  "OP_13",
	94:  "OP_14",
	95:  "OP_15",
	96:  "OP_16",
	97:  "OP_NOP",
	99:  "OP_IF",
	100: "OP_NOTIF",
	103: "OP_ELSE",
	104: "OP_ENDIF",
	105: "OP_VERIFY",
	106: "OP_RETURN",
	107: "OP_TOALTSTACK",
	108: "OP_FROMALTSTACK",
	109: "OP_2DROP",
	110: "OP_2DUP",
	111: "OP_3DUP",
	112: "OP_2OVER",
	113: "OP_2ROT",
	114: "OP_2SWAP",
	115: "OP_IFDUP",
	116: "OP_DEPTH",
	117: "OP_DROP",
	118: "OP_DUP",
	119: "OP_NIP",
	120: "OP_OVER",
	121: "OP_PICK",
	122: "OP_ROLL",
	123: "OP_ROT",
	124: "OP_SWAP",
	125: "OP_TUCK",
	130: "OP_SIZE",
	135: "OP_EQUAL",
	136: "OP_EQUALVERIFY",
	139: "OP_1ADD",
	140: "OP_1SUB",
	143: "OP_NEGATE",
	144: "OP_ABS",
	145: "OP_NOT",
	146: "OP_0NOTEQUAL",
	147: "OP_ADD",
	148: "OP_SUB",
	149: "OP_MUL",
	154: "OP_BOOLAND",
	155: "OP_BOOLOR",
	156: "OP_NUMEQUAL",
	157: "OP_NUMEQUALVERIFY",
	158: "OP_NUMNOTEQUAL",
	159: "OP_LESSTHAN",
	160: "OP_GREATERTHAN",
	161: "OP_LESSTHANOREQUAL",
	162: "OP_GREATERTHANOREQUAL",
	163: "OP_MIN",
	164: "OP_MAX",
	165: "OP_WITHIN",
	166: "OP_RIPEMD160",
	167: "OP_SHA1",
	168: "OP_SHA256",
	169: "OP_HASH160",
	170: "OP_HASH256",
	171: "OP_CODESEPARATOR",
	172: "OP_CHECKSIG",
	173: "OP_CHECKSIGVERIFY",
	174: "OP_CHECKMULTISIG",
	175: "OP_CHECKMULTISIGVERIFY",
	176: "OP_NOP1",
	177: "OP_CHECKLOCKTIMEVERIFY",
	178: "OP_CHECKSEQUENCEVERIFY",
	179: "OP_NOP4",
	180: "OP_NOP5",
	181: "OP_NOP6",
	182: "OP_NOP7",
	183: "OP_NOP8",
	184: "OP_NOP9",
	185: "OP_NOP10",
}
