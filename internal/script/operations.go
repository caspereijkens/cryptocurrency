package script

import (
	"fmt"
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
	*stack = append(*stack, encodeNum(0))
	return true, nil
}

func op1Negate(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(-1))
	return true, nil
}

func op1(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(1))
	return true, nil
}

func op2(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(2))
	return true, nil
}

func op3(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(3))
	return true, nil
}

func op4(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(4))
	return true, nil
}

func op5(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(5))
	return true, nil
}

func op6(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(6))
	return true, nil
}

func op7(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(7))
	return true, nil
}

func op8(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(8))
	return true, nil
}

func op9(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(9))
	return true, nil
}

func op10(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(10))
	return true, nil
}

func op11(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(11))
	return true, nil
}

func op12(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(12))
	return true, nil
}

func op13(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(13))
	return true, nil
}

func op14(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(14))
	return true, nil
}

func op15(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(15))
	return true, nil
}

func op16(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(16))
	return true, nil
}

func opNop(stack *Stack) (bool, error) {
	return true, nil
}

// TODO opIf and opNotIf

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

	*stack = append(*stack, element)

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
		*stack = append(*stack, element)
	}

	return true, nil
}

func opDepth(stack *Stack) (bool, error) {
	*stack = append(*stack, encodeNum(len(*stack)))
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

	*stack = append(*stack, element)

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

	*stack = append(*stack, (*stack)[len(*stack)-2])

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

	*stack = append(*stack, (*stack)[len(*stack)-n-1])

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
		*stack = append(*stack, rolled)
	}

	return true, nil
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
