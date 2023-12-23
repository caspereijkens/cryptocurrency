package script

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

func op0(stack *Stack) bool {
	*stack = append(*stack, encodeNum(0))
	return true
}

func op1Negate(stack *Stack) bool {
	*stack = append(*stack, encodeNum(-1))
	return true
}

func op1(stack *Stack) bool {
	*stack = append(*stack, encodeNum(1))
	return true
}

func op2(stack *Stack) bool {
	*stack = append(*stack, encodeNum(2))
	return true
}

func op3(stack *Stack) bool {
	*stack = append(*stack, encodeNum(3))
	return true
}

func op4(stack *Stack) bool {
	*stack = append(*stack, encodeNum(4))
	return true
}

func op5(stack *Stack) bool {
	*stack = append(*stack, encodeNum(5))
	return true
}

func op6(stack *Stack) bool {
	*stack = append(*stack, encodeNum(6))
	return true
}

func op7(stack *Stack) bool {
	*stack = append(*stack, encodeNum(7))
	return true
}

func op8(stack *Stack) bool {
	*stack = append(*stack, encodeNum(8))
	return true
}

func op9(stack *Stack) bool {
	*stack = append(*stack, encodeNum(9))
	return true
}

func op10(stack *Stack) bool {
	*stack = append(*stack, encodeNum(10))
	return true
}

func op11(stack *Stack) bool {
	*stack = append(*stack, encodeNum(11))
	return true
}

func op12(stack *Stack) bool {
	*stack = append(*stack, encodeNum(12))
	return true
}

func op13(stack *Stack) bool {
	*stack = append(*stack, encodeNum(13))
	return true
}

func op14(stack *Stack) bool {
	*stack = append(*stack, encodeNum(14))
	return true
}

func op15(stack *Stack) bool {
	*stack = append(*stack, encodeNum(15))
	return true
}

func op16(stack *Stack) bool {
	*stack = append(*stack, encodeNum(16))
	return true
}

func opNop(stack *Stack) bool {
	return true
}

// TODO opIf and opNotIf

func opVerify(stack *Stack) bool {
	if len(*stack) < 1 {
		return false
	}

	element := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	return decodeNum(element) != 0
}

func opReturn(stack *Stack) bool {
	return false
}

func opToAltStack(stack, altStack *Stack) bool {
	if len(*stack) < 1 {
		return false
	}

	*altStack = append(*altStack, (*stack)[len(*stack)-1])
	*stack = (*stack)[:len(*stack)-1]

	return true
}

func opFromAltStack(stack, altStack *Stack) bool {
	if len(*altStack) < 1 {
		return false
	}

	*stack = append(*stack, (*altStack)[len(*altStack)-1])
	*altStack = (*altStack)[:len(*altStack)-1]

	return true
}

func op2Drop(stack *Stack) bool {
	if len(*stack) < 2 {
		return false
	}

	*stack = (*stack)[:len(*stack)-2]
	return true
}

func op2Dup(stack *Stack) bool {
	if len(*stack) < 2 {
		return false
	}

	*stack = append(*stack, (*stack)[len(*stack)-2:]...)
	return true
}

func op3Dup(stack *Stack) bool {
	if len(*stack) < 3 {
		return false
	}

	*stack = append(*stack, (*stack)[len(*stack)-3:]...)
	return true
}

func op2Over(stack *Stack) bool {
	if len(*stack) < 4 {
		return false
	}

	*stack = append(*stack, (*stack)[len(*stack)-4:len(*stack)-2]...)
	return true
}

func op2Rot(stack *Stack) bool {
	if len(*stack) < 6 {
		return false
	}

	*stack = append(*stack, (*stack)[len(*stack)-6:len(*stack)-4]...)
	return true
}

func op2Swap(stack *Stack) bool {
	if len(*stack) < 4 {
		return false
	}

	lastFour := (*stack)[len(*stack)-4:]
	(*stack)[len(*stack)-4] = lastFour[2]
	(*stack)[len(*stack)-3] = lastFour[3]
	(*stack)[len(*stack)-2] = lastFour[0]
	(*stack)[len(*stack)-1] = lastFour[1]

	return true
}

func opIfDup(stack *Stack) bool {
	if len(*stack) < 1 {
		return false
	}

	lastElement := (*stack)[len(*stack)-1]
	if decodeNum(lastElement) != 0 {
		*stack = append(*stack, lastElement)
	}

	return true
}

func opDepth(stack *Stack) bool {
	*stack = append(*stack, encodeNum(len(*stack)))
	return true
}

func opDrop(stack *Stack) bool {
	if len(*stack) < 1 {
		return false
	}

	*stack = (*stack)[:len(*stack)-1]
	return true
}

func opDup(stack *Stack) bool {
	if len(*stack) < 1 {
		return false
	}

	lastElement := (*stack)[len(*stack)-1]
	*stack = append(*stack, lastElement)

	return true
}

func opNip(stack *Stack) bool {
	if len(*stack) < 2 {
		return false
	}

	*stack = append((*stack)[:len(*stack)-2], (*stack)[len(*stack)-1])
	return true
}

func opOver(stack *Stack) bool {
	if len(*stack) < 2 {
		return false
	}

	*stack = append(*stack, (*stack)[len(*stack)-2])

	return true
}

func opPick(stack *Stack) bool {
	if len(*stack) < 1 {
		return false
	}

	n := decodeNum((*stack)[len(*stack)-1])
	*stack = (*stack)[:len(*stack)-1]

	if len(*stack) < n+1 {
		return false
	}

	*stack = append(*stack, (*stack)[len(*stack)-n-1])

	return true
}

func opRoll(stack *Stack) bool {
	if len(*stack) < 1 {
		return false
	}

	n := decodeNum((*stack)[len(*stack)-1])
	*stack = (*stack)[:len(*stack)-1]

	if len(*stack) < n+1 {
		return false
	}

	if n > 0 {
		rolled := (*stack)[len(*stack)-n-1]
		*stack = append((*stack)[:len(*stack)-n-1], (*stack)[len(*stack)-n:]...)
		*stack = append(*stack, rolled)
	}

	return true
}
