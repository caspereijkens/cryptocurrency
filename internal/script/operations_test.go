package script

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEncodeDecodeNum(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{42, 42},
		{-42, -42},
		{127, 127},
		{-127, -127},
		{128, 128},
		{-128, -128},
		{300, 300},
		{-300, -300},
	}

	for _, test := range tests {
		encoded := encodeNum(test.input)
		decoded := decodeNum(encoded)

		if decoded != test.expected {
			t.Errorf("Failed for input %d. Expected %d, got %d", test.input, test.expected, decoded)
		}
	}
}

func TestIntegerOperations(t *testing.T) {
	var stack Stack

	// Define the operations in a map for easy iteration
	operations := []func(*Stack) bool{op1, op2, op3, op4, op5, op6, op7, op8, op9, op10, op11, op12, op13, op14, op15, op16}

	// Perform dynamic tests for each operation
	for i, op := range operations {
		expected := decodeNum(encodeNum(i + 1)) // For op1, it should be 1, for op2, it should be 2, and so on
		t.Run(fmt.Sprintf("op%d", i), func(t *testing.T) {
			performOperation(op, &stack, expected, t)
		})
	}
}

func TestOp1Negate(t *testing.T) {
	stack := new(Stack)
	op1Negate(stack)

	expected := encodeNum(-1)
	if !bytes.Equal((*stack)[0], expected) {
		t.Errorf("Failed for op_1negate. Expected %v, got %v", expected, (*stack)[0])
	}
}

func TestOperations(t *testing.T) {
	// Test all operations together
	stack := new(Stack)
	op0(stack)
	op1Negate(stack)
	op1(stack)

	expected0 := []byte{}
	if !bytes.Equal((*stack)[0], expected0) {
		t.Errorf("Failed for op_0. Expected %v, got %v", expected0, (*stack)[0])
	}

	expected1Negate := encodeNum(-1)
	if !bytes.Equal((*stack)[1], expected1Negate) {
		t.Errorf("Failed for op_1negate. Expected %v, got %v", expected1Negate, (*stack)[1])
	}

	expected1 := encodeNum(1)
	if !bytes.Equal((*stack)[2], expected1) {
		t.Errorf("Failed for op_1. Expected %v, got %v", expected1, (*stack)[2])
	}
}

func TestOpNop(t *testing.T) {
	var stack Stack

	// Call the opNop function
	opNop(&stack)

	// Check that the stack remains unchanged
	if len(stack) != 0 {
		t.Errorf("opNop should not modify the stack. Expected length 0, got %d", len(stack))
	}
}

func TestOpVerify(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack := opVerify(&emptyStack)
	if resultEmptyStack {
		t.Errorf("opVerify failed for empty stack. Expected false, got true")
	}

	// Test when the top element of the stack is 0
	stackWithZero := Stack{encodeNum(0)}
	resultWithZero := opVerify(&stackWithZero)
	if resultWithZero {
		t.Errorf("opVerify failed for stack with top element 0. Expected false, got true")
	}

	// Test when the top element of the stack is non-zero
	stackWithNonZero := Stack{encodeNum(42)}
	resultWithNonZero := opVerify(&stackWithNonZero)
	if !resultWithNonZero {
		t.Errorf("opVerify failed for stack with top element 42. Expected true, got false")
	}
}

func TestOpReturn(t *testing.T) {
	stack := Stack{encodeNum(42)} // Sample stack with one element

	// Call opReturn and check the result
	result := opReturn(&stack)

	// opReturn should always return false
	if result {
		t.Errorf("opReturn failed. Expected false, got true")
	}
}

func TestOpToAltStack(t *testing.T) {
	stack := Stack{encodeNum(42)} // Sample stack with one element
	altStack := Stack{}           // Empty alternative stack

	// Call opToAltStack and check the result
	result := opToAltStack(&stack, &altStack)

	// The top element of stack should be moved to altStack
	if !result || len(stack) != 0 || len(altStack) != 1 || decodeNum(altStack[0]) != 42 {
		t.Errorf("opToAltStack failed. Unexpected state after the operation")
	}
}

func TestOpFromAltStack(t *testing.T) {
	stack := Stack{}                 // Empty stack
	altStack := Stack{encodeNum(42)} // Sample alternative stack with one element

	// Call opFromAltStack and check the result
	result := opFromAltStack(&stack, &altStack)

	// The top element of altStack should be moved to stack
	if !result || len(stack) != 1 || len(altStack) != 0 || decodeNum(stack[0]) != 42 {
		t.Errorf("opFromAltStack failed. Unexpected state after the operation")
	}
}

func TestOp2Drop(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(4)}

	result := op2Drop(&stack)

	if !result || len(stack) != 2 || decodeNum(stack[0]) != 1 || decodeNum(stack[1]) != 2 {
		t.Errorf("op2Drop failed. Unexpected state after the operation")
	}
}

func TestOp2Dup(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2)}

	result := op2Dup(&stack)

	if !result || len(stack) != 4 || decodeNum(stack[2]) != 1 || decodeNum(stack[3]) != 2 {
		t.Errorf("op2Dup failed. Unexpected state after the operation")
	}
}

func TestOp3Dup(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}

	result := op3Dup(&stack)

	if !result || len(stack) != 6 || decodeNum(stack[3]) != 1 || decodeNum(stack[4]) != 2 || decodeNum(stack[5]) != 3 {
		t.Errorf("op3Dup failed. Unexpected state after the operation")
	}
}

func TestOp2Over(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(4)}

	result := op2Over(&stack)

	if !result || len(stack) != 6 || decodeNum(stack[4]) != 1 || decodeNum(stack[5]) != 2 {
		t.Errorf("op2Over failed. Unexpected state after the operation")
	}
}

func TestOp2Rot(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(4), encodeNum(5), encodeNum(6)}

	result := op2Rot(&stack)

	if !result || len(stack) != 8 || decodeNum(stack[6]) != 1 || decodeNum(stack[7]) != 2 {
		t.Errorf("op2Rot failed. Unexpected state after the operation")
	}
}

func TestOp2Swap(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(4), encodeNum(5), encodeNum(6)}

	result := op2Swap(&stack)

	if !result || len(stack) != 6 || decodeNum(stack[2]) != 5 || decodeNum(stack[3]) != 6 {
		t.Errorf("op2Swap failed. Unexpected state after the operation")
	}
}

func TestOpIfDup(t *testing.T) {
	stack := Stack{encodeNum(0)}

	result := opIfDup(&stack)

	if !result || len(stack) == 2 {
		t.Errorf("opIfDup failed. Unexpected state after the operation")
	}

	stack = Stack{encodeNum(42)}

	result = opIfDup(&stack)

	if !result || len(stack) != 2 || decodeNum(stack[1]) != 42 {
		t.Errorf("opIfDup failed. Unexpected state after the operation")
	}
}

func TestOpDepth(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}

	result := opDepth(&stack)

	if !result || len(stack) != 4 || decodeNum(stack[3]) != 3 {
		t.Errorf("opDepth failed. Unexpected state after the operation")
	}
}

func TestOpDrop(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}

	result := opDrop(&stack)

	if !result || len(stack) != 2 || decodeNum(stack[1]) != 2 {
		t.Errorf("opDrop failed. Unexpected state after the operation")
	}
}

func TestOpDup(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack := opDup(&emptyStack)
	if resultEmptyStack {
		t.Errorf("opDup failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42)}
	resultOneElement := opDup(&stackWithOneElement)
	if !resultOneElement || len(stackWithOneElement) != 2 || decodeNum(stackWithOneElement[1]) != 42 {
		t.Errorf("opDup failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}
	resultMultipleElements := opDup(&stackWithMultipleElements)
	if !resultMultipleElements || len(stackWithMultipleElements) != 4 || decodeNum(stackWithMultipleElements[3]) != 3 {
		t.Errorf("opDup failed for stack with multiple elements. Unexpected state after the operation")
	}
}

func TestOpNip(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack := opNip(&emptyStack)
	if resultEmptyStack {
		t.Errorf("opNip failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42)}
	resultOneElement := opNip(&stackWithOneElement)
	if resultOneElement || len(stackWithOneElement) != 1 || decodeNum(stackWithOneElement[0]) != 42 {
		t.Errorf("opNip failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}
	resultMultipleElements := opNip(&stackWithMultipleElements)
	if !resultMultipleElements || len(stackWithMultipleElements) != 2 || decodeNum(stackWithMultipleElements[1]) != 3 {
		t.Errorf("opNip failed for stack with multiple elements. Unexpected state after the operation")
	}
}

func TestOpOver(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack := opOver(&emptyStack)
	if resultEmptyStack {
		t.Errorf("opOver failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42)}
	resultOneElement := opOver(&stackWithOneElement)
	if resultOneElement || len(stackWithOneElement) != 1 || decodeNum(stackWithOneElement[0]) != 42 {
		t.Errorf("opOver failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}
	resultMultipleElements := opOver(&stackWithMultipleElements)
	if !resultMultipleElements || len(stackWithMultipleElements) != 4 || decodeNum(stackWithMultipleElements[3]) != 2 {
		t.Errorf("opOver failed for stack with multiple elements. Unexpected state after the operation")
	}
}

func TestOpPick(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack := opPick(&emptyStack)
	if resultEmptyStack {
		t.Errorf("opPick failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42), encodeNum(0)}
	resultOneElement := opPick(&stackWithOneElement)
	if !resultOneElement || len(stackWithOneElement) != 2 || decodeNum(stackWithOneElement[1]) != 42 {
		t.Errorf("opPick failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(1)}
	resultMultipleElements := opPick(&stackWithMultipleElements)
	if !resultMultipleElements || len(stackWithMultipleElements) != 4 || decodeNum(stackWithMultipleElements[3]) != 2 {
		t.Errorf("opPick failed for stack with multiple elements. Unexpected state after the operation")
	}

	// Test when the stack does not have enough elements for pick
	stackNotEnoughElements := Stack{encodeNum(1)}
	resultNotEnoughElements := opPick(&stackNotEnoughElements)
	if resultNotEnoughElements || len(stackNotEnoughElements) != 0 {
		t.Errorf("opPick failed for stack with not enough elements. Unexpected state after the operation")
	}
}

func TestOpRoll(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack := opRoll(&emptyStack)
	if resultEmptyStack {
		t.Errorf("opRoll failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42), encodeNum(0)}
	resultOneElement := opRoll(&stackWithOneElement)
	if !resultOneElement || len(stackWithOneElement) != 1 || decodeNum(stackWithOneElement[0]) != 42 {
		t.Errorf("opRoll failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(2)}
	resultMultipleElements := opRoll(&stackWithMultipleElements)
	if !resultMultipleElements || len(stackWithMultipleElements) != 3 || decodeNum(stackWithMultipleElements[2]) != 1 {
		t.Errorf("opRoll failed for stack with multiple elements. Unexpected state after the operation")
	}

	// Test when the stack does not have enough elements for roll
	stackNotEnoughElements := Stack{encodeNum(1)}
	resultNotEnoughElements := opRoll(&stackNotEnoughElements)
	if resultNotEnoughElements || len(stackNotEnoughElements) != 0 {
		t.Errorf("opRoll failed for stack with not enough elements. Unexpected state after the operation")
	}

	// Test roll with n out of bounds
	stackWithZeroN := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(99)}
	resultZeroN := opRoll(&stackWithZeroN)
	if resultZeroN || len(stackWithZeroN) != 3 {
		t.Errorf("opRoll failed for n=0. Unexpected state after the operation")
	}
}

func performOperation(op func(*Stack) bool, stack *Stack, expected int, t *testing.T) {
	op(stack)
	result := decodeNum((*stack)[len(*stack)-1])

	if result != expected {
		t.Errorf("Failed for %s. Expected %d, got %d", getOpName(op), expected, result)
	}
}

func getOpName(op interface{}) string {
	return fmt.Sprintf("%p", op)
}
