package script

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"

	"github.com/caspereijkens/cryptocurrency/internal/utils"
	"golang.org/x/crypto/ripemd160"
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
	operations := []func(*Stack) (bool, error){op1, op2, op3, op4, op5, op6, op7, op8, op9, op10, op11, op12, op13, op14, op15, op16}

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

func TestOpIf(t *testing.T) {
	// Test case 1: Regular if statement without else
	stack1 := &Stack{encodeNum(1)}
	items1 := &Stack{encodeNum(1), encodeNum(104)}

	result1, err1 := opIf(stack1, items1)
	assertOpIfResult(t, result1, err1, stack1, &Stack{}, items1, &Stack{encodeNum(1)}, true)

	// Test case 2: Regular if statement (false condition)
	stack2 := &Stack{encodeNum(0)}
	items2 := &Stack{encodeNum(5), encodeNum(104)}

	result2, err2 := opIf(stack2, items2)
	assertOpIfResult(t, result2, err2, stack2, &Stack{}, items2, &Stack{}, true)

	// Test case 3: Empty stack
	stack3 := &Stack{}
	items3 := &Stack{encodeNum(1), encodeNum(104)}

	_, err3 := opIf(stack3, items3)
	assertOpIfError(t, err3, "stack is empty")

	// Test case 4: Nested if/else statement
	stack4 := &Stack{encodeNum(1)}
	items4 := &Stack{encodeNum(99), encodeNum(1), encodeNum(103), encodeNum(2), encodeNum(104), encodeNum(104)}

	result4, err4 := opIf(stack4, items4)
	assertOpIfResult(t, result4, err4, stack4, &Stack{}, items4, &Stack{encodeNum(99), encodeNum(1), encodeNum(103), encodeNum(2), encodeNum(104)}, true)

	// Test case 5 if/else statement
	stack5 := &Stack{encodeNum(0)}
	items5 := &Stack{encodeNum(1), encodeNum(103), encodeNum(2), encodeNum(104)}

	result5, err5 := opIf(stack5, items5)
	assertOpIfResult(t, result5, err5, stack5, &Stack{}, items5, &Stack{encodeNum(2)}, true)
}

func TestOpNotIf(t *testing.T) {
	// Test case 1: Regular if statement without else
	stack1 := &Stack{encodeNum(0)}
	items1 := &Stack{encodeNum(1), encodeNum(104)}

	result1, err1 := opNotIf(stack1, items1)
	assertOpNotIfResult(t, result1, err1, stack1, &Stack{}, items1, &Stack{encodeNum(1)}, true)

	// Test case 2: Regular if statement (false condition)
	stack2 := &Stack{encodeNum(1)}
	items2 := &Stack{encodeNum(5), encodeNum(104)}

	result2, err2 := opNotIf(stack2, items2)
	assertOpNotIfResult(t, result2, err2, stack2, &Stack{}, items2, &Stack{}, true)

	// Test case 3: Empty stack
	stack3 := &Stack{}
	items3 := &Stack{encodeNum(1), encodeNum(104)}

	_, err3 := opNotIf(stack3, items3)
	assertOpNotIfError(t, err3, "stack is empty")

	// Test case 4: Nested if/else statement
	stack4 := &Stack{encodeNum(0)}
	items4 := &Stack{encodeNum(99), encodeNum(1), encodeNum(103), encodeNum(2), encodeNum(104), encodeNum(104)}

	result4, err4 := opNotIf(stack4, items4)
	assertOpNotIfResult(t, result4, err4, stack4, &Stack{}, items4, &Stack{encodeNum(99), encodeNum(1), encodeNum(103), encodeNum(2), encodeNum(104)}, true)

	// Test case 5 if/else statement
	stack5 := &Stack{encodeNum(1)}
	items5 := &Stack{encodeNum(1), encodeNum(103), encodeNum(2), encodeNum(104)}

	result5, err5 := opNotIf(stack5, items5)
	assertOpNotIfResult(t, result5, err5, stack5, &Stack{}, items5, &Stack{encodeNum(2)}, true)

}

func assertOpIfResult(t *testing.T, result bool, err error, stack, expectedStack, items, expectedItems *Stack, success bool) {
	t.Helper()

	if success {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !equalStacks(stack, expectedStack) || !equalStacks(items, expectedItems) || !result {
			t.Errorf("opIf test failed")
		}
	} else {
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	}
}

func assertOpNotIfResult(t *testing.T, result bool, err error, stack, expectedStack, items, expectedItems *Stack, success bool) {
	t.Helper()

	if success {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !equalStacks(stack, expectedStack) || !equalStacks(items, expectedItems) || !result {
			t.Errorf("opNotIf test failed")
		}
	} else {
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	}
}

func assertOpIfError(t *testing.T, err error, expectedError string) {
	t.Helper()

	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %v, got: %v", expectedError, err)
	}
}

func assertOpNotIfError(t *testing.T, err error, expectedError string) {
	t.Helper()

	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %v, got: %v", expectedError, err)
	}
}

func equalStacks(s1, s2 *Stack) bool {
	if len(*s1) != len(*s2) {
		return false
	}
	for i := range *s1 {
		if !bytes.Equal((*s1)[i], (*s2)[i]) {
			return false
		}
	}
	return true
}

func TestOpVerify(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opVerify(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opVerify failed for empty stack. Expected false, got true")
	}

	// Test when the top element of the stack is 0
	stackWithZero := Stack{encodeNum(0)}
	resultWithZero, err := opVerify(&stackWithZero)
	if resultWithZero || err != nil {
		t.Errorf("opVerify failed for stack with top element 0. Expected false, got true")
	}

	// Test when the top element of the stack is non-zero
	stackWithNonZero := Stack{encodeNum(42)}
	resultWithNonZero, err := opVerify(&stackWithNonZero)
	if !resultWithNonZero || err != nil {
		t.Errorf("opVerify failed for stack with top element 42. Expected true, got false")
	}
}

func TestOpReturn(t *testing.T) {
	stack := Stack{encodeNum(42)} // Sample stack with one element

	// Call opReturn and check the result
	result, err := opReturn(&stack)

	// opReturn should always return false
	if result || err != nil {
		t.Errorf("opReturn failed. Expected false, got true")
	}
}

func TestOpToAltStack(t *testing.T) {
	stack := Stack{encodeNum(42)} // Sample stack with one element
	altStack := Stack{}           // Empty alternative stack

	// Call opToAltStack and check the result
	result, err := opToAltStack(&stack, &altStack)

	// The top element of stack should be moved to altStack
	if !result || err != nil || len(stack) != 0 || len(altStack) != 1 || decodeNum(altStack[0]) != 42 {
		t.Errorf("opToAltStack failed. Unexpected state after the operation")
	}
}

func TestOpFromAltStack(t *testing.T) {
	stack := Stack{}                 // Empty stack
	altStack := Stack{encodeNum(42)} // Sample alternative stack with one element

	// Call opFromAltStack and check the result
	result, err := opFromAltStack(&stack, &altStack)

	// The top element of altStack should be moved to stack
	if !result || err != nil || len(stack) != 1 || len(altStack) != 0 || decodeNum(stack[0]) != 42 {
		t.Errorf("opFromAltStack failed. Unexpected state after the operation")
	}
}

func TestOp2Drop(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(4)}

	result, err := op2Drop(&stack)

	if !result || err != nil || len(stack) != 2 || decodeNum(stack[0]) != 1 || decodeNum(stack[1]) != 2 {
		t.Errorf("op2Drop failed. Unexpected state after the operation")
	}
}

func TestOp2Dup(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2)}

	result, err := op2Dup(&stack)

	if !result || err != nil || len(stack) != 4 || decodeNum(stack[2]) != 1 || decodeNum(stack[3]) != 2 {
		t.Errorf("op2Dup failed. Unexpected state after the operation")
	}
}

func TestOp3Dup(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}

	result, err := op3Dup(&stack)

	if !result || err != nil || len(stack) != 6 || decodeNum(stack[3]) != 1 || decodeNum(stack[4]) != 2 || decodeNum(stack[5]) != 3 {
		t.Errorf("op3Dup failed. Unexpected state after the operation")
	}
}

func TestOp2Over(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(4)}

	result, err := op2Over(&stack)

	if !result || err != nil || len(stack) != 6 || decodeNum(stack[4]) != 1 || decodeNum(stack[5]) != 2 {
		t.Errorf("op2Over failed. Unexpected state after the operation")
	}
}

func TestOp2Rot(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(4), encodeNum(5), encodeNum(6)}

	result, err := op2Rot(&stack)

	if !result || err != nil || len(stack) != 8 || decodeNum(stack[6]) != 1 || decodeNum(stack[7]) != 2 {
		t.Errorf("op2Rot failed. Unexpected state after the operation")
	}
}

func TestOp2Swap(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(4), encodeNum(5), encodeNum(6)}

	result, err := op2Swap(&stack)

	if !result || err != nil || len(stack) != 6 || decodeNum(stack[2]) != 5 || decodeNum(stack[3]) != 6 {
		t.Errorf("op2Swap failed. Unexpected state after the operation")
	}
}

func TestOpIfDup(t *testing.T) {
	stack := Stack{encodeNum(0)}

	result, err := opIfDup(&stack)

	if !result || err != nil || len(stack) == 2 {
		t.Errorf("opIfDup failed. Unexpected state after the operation")
	}

	stack = Stack{encodeNum(42)}

	result, err = opIfDup(&stack)

	if !result || err != nil || len(stack) != 2 || decodeNum(stack[1]) != 42 {
		t.Errorf("opIfDup failed. Unexpected state after the operation")
	}
}

func TestOpDepth(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}

	result, err := opDepth(&stack)

	if !result || err != nil || len(stack) != 4 || decodeNum(stack[3]) != 3 {
		t.Errorf("opDepth failed. Unexpected state after the operation")
	}
}

func TestOpDrop(t *testing.T) {
	stack := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}

	result, err := opDrop(&stack)

	if !result || err != nil || len(stack) != 2 || decodeNum(stack[1]) != 2 {
		t.Errorf("opDrop failed. Unexpected state after the operation")
	}
}

func TestOpDup(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opDup(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opDup failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42)}
	resultOneElement, err := opDup(&stackWithOneElement)
	if !resultOneElement || err != nil || len(stackWithOneElement) != 2 || decodeNum(stackWithOneElement[1]) != 42 {
		t.Errorf("opDup failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}
	resultMultipleElements, err := opDup(&stackWithMultipleElements)
	if !resultMultipleElements || err != nil || len(stackWithMultipleElements) != 4 || decodeNum(stackWithMultipleElements[3]) != 3 {
		t.Errorf("opDup failed for stack with multiple elements. Unexpected state after the operation")
	}
}

func TestOpNip(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opNip(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opNip failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42)}
	resultOneElement, err := opNip(&stackWithOneElement)
	if resultOneElement || err == nil || len(stackWithOneElement) != 1 || decodeNum(stackWithOneElement[0]) != 42 {
		t.Errorf("opNip failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}
	resultMultipleElements, err := opNip(&stackWithMultipleElements)
	if !resultMultipleElements || err != nil || len(stackWithMultipleElements) != 2 || decodeNum(stackWithMultipleElements[1]) != 3 {
		t.Errorf("opNip failed for stack with multiple elements. Unexpected state after the operation")
	}
}

func TestOpOver(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opOver(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opOver failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42)}
	resultOneElement, err := opOver(&stackWithOneElement)
	if resultOneElement || err == nil || len(stackWithOneElement) != 1 || decodeNum(stackWithOneElement[0]) != 42 {
		t.Errorf("opOver failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}
	resultMultipleElements, err := opOver(&stackWithMultipleElements)
	if !resultMultipleElements || err != nil || len(stackWithMultipleElements) != 4 || decodeNum(stackWithMultipleElements[3]) != 2 {
		t.Errorf("opOver failed for stack with multiple elements. Unexpected state after the operation")
	}
}

func TestOpPick(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opPick(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opPick failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42), encodeNum(0)}
	resultOneElement, err := opPick(&stackWithOneElement)
	if !resultOneElement || err != nil || len(stackWithOneElement) != 2 || decodeNum(stackWithOneElement[1]) != 42 {
		t.Errorf("opPick failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(1)}
	resultMultipleElements, err := opPick(&stackWithMultipleElements)
	if !resultMultipleElements || err != nil || len(stackWithMultipleElements) != 4 || decodeNum(stackWithMultipleElements[3]) != 2 {
		t.Errorf("opPick failed for stack with multiple elements. Unexpected state after the operation")
	}

	// Test when the stack does not have enough elements for pick
	stackNotEnoughElements := Stack{encodeNum(1)}
	resultNotEnoughElements, err := opPick(&stackNotEnoughElements)
	if resultNotEnoughElements || err == nil || len(stackNotEnoughElements) != 0 {
		t.Errorf("opPick failed for stack with not enough elements. Unexpected state after the operation")
	}
}

func TestOpRoll(t *testing.T) {
	// Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opRoll(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opRoll failed for empty stack. Expected false, got true")
	}

	// Test when the stack has one element
	stackWithOneElement := Stack{encodeNum(42), encodeNum(0)}
	resultOneElement, err := opRoll(&stackWithOneElement)
	if !resultOneElement || err != nil || len(stackWithOneElement) != 1 || decodeNum(stackWithOneElement[0]) != 42 {
		t.Errorf("opRoll failed for stack with one element. Unexpected state after the operation")
	}

	// Test when the stack has multiple elements
	stackWithMultipleElements := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(2)}
	resultMultipleElements, err := opRoll(&stackWithMultipleElements)
	if !resultMultipleElements || err != nil || len(stackWithMultipleElements) != 3 || decodeNum(stackWithMultipleElements[2]) != 1 {
		t.Errorf("opRoll failed for stack with multiple elements. Unexpected state after the operation")
	}

	// Test when the stack does not have enough elements for roll
	stackNotEnoughElements := Stack{encodeNum(1)}
	resultNotEnoughElements, err := opRoll(&stackNotEnoughElements)
	if resultNotEnoughElements || err == nil || len(stackNotEnoughElements) != 0 {
		t.Errorf("opRoll failed for stack with not enough elements. Unexpected state after the operation")
	}

	// Test roll with n out of bounds
	stackWithZeroN := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(99)}
	resultZeroN, err := opRoll(&stackWithZeroN)
	if resultZeroN || err == nil || len(stackWithZeroN) != 3 {
		t.Errorf("opRoll failed for n=0. Unexpected state after the operation")
	}
}

func TestOpRot(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opRot(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opRot failed for empty stack. Expected false, got true")
	}

	// Test case 2: Test when the stack has less than 3 elements
	stackLessThan3 := Stack{encodeNum(1), encodeNum(2)}
	resultLessThan3, err := opRot(&stackLessThan3)
	if resultLessThan3 || err == nil || len(stackLessThan3) != 2 {
		t.Errorf("opRot failed for stack with less than 3 elements. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has 3 or more elements
	stack3OrMore := Stack{encodeNum(1), encodeNum(2), encodeNum(3), encodeNum(4)}
	result3OrMore, err := opRot(&stack3OrMore)
	if !result3OrMore || err != nil || len(stack3OrMore) != 4 || decodeNum(stack3OrMore[3]) != 2 {
		t.Errorf("opRot failed for stack with 3 or more elements. Unexpected state after the operation")
	}
}

func TestOpSwap(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opSwap(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opSwap failed for empty stack. Expected false, got true")
	}

	// Test case 2: Test when the stack has less than 2 elements
	stackLessThan2 := Stack{encodeNum(1)}
	resultLessThan2, err := opSwap(&stackLessThan2)
	if resultLessThan2 || err == nil || len(stackLessThan2) != 1 {
		t.Errorf("opSwap failed for stack with less than 2 elements. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has 2 or more elements
	stack2OrMore := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}
	result2OrMore, err := opSwap(&stack2OrMore)
	if !result2OrMore || err != nil || len(stack2OrMore) != 3 || decodeNum(stack2OrMore[2]) != 2 {
		t.Errorf("opSwap failed for stack with 2 or more elements. Unexpected state after the operation")
	}
}

func TestOpTuck(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opTuck(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opTuck failed for empty stack. Expected false, got true")
	}

	// Test case 2: Test when the stack has less than 1 element
	stackLessThan2 := Stack{encodeNum(1)}
	resultLessThan1, err := opTuck(&stackLessThan2)
	if resultLessThan1 || err == nil || len(stackLessThan2) != 1 {
		t.Errorf("opTuck failed for stack with less than 1 element. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has 1 or more elements
	stack2OrMore := Stack{encodeNum(1), encodeNum(2), encodeNum(3)}
	result1OrMore, err := opTuck(&stack2OrMore)
	if !result1OrMore || err != nil || len(stack2OrMore) != 4 || decodeNum(stack2OrMore[3]) != 3 {
		t.Errorf("opTuck failed for stack with 1 or more elements. Unexpected state after the operation")
	}
}

func TestOpSize(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opSize(&emptyStack)
	if resultEmptyStack || err == nil || err.Error() != "stack is empty" {
		t.Errorf("opSize failed for empty stack. Expected false, error 'stack is empty'; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element
	stackWithElement := Stack{[]byte{1, 2, 3}}
	resultWithElement, err := opSize(&stackWithElement)
	if !resultWithElement || err != nil || len(stackWithElement) != 2 || decodeNum(stackWithElement[len(stackWithElement)-1]) != 3 {
		t.Errorf("opSize failed for stack with at least 1 element. Unexpected state after the operation")
	}
}

func TestOpEqual(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opEqual(&emptyStack)
	if resultEmptyStack || err == nil || err.Error() != "not enough elements in stack: 0 < 2" {
		t.Errorf("opEqual failed for empty stack. Expected false, error 'not enough elements in stack: 0 < 2'; got true, %v", err)
	}

	// Test case 2: Test when the stack has less than 2 elements
	stackLessThan2 := Stack{[]byte{1}}
	resultLessThan2, err := opEqual(&stackLessThan2)
	if resultLessThan2 || err == nil || err.Error() != "not enough elements in stack: 1 < 2" {
		t.Errorf("opEqual failed for stack with less than 2 elements. Expected false, error 'not enough elements in stack: 1 < 2'; got true, %v", err)
	}

	// Test case 3: Test when the stack has 2 or more elements, and they are equal
	stackEqual := Stack{[]byte{1, 2, 3}, []byte{1, 2, 3}}
	resultEqual, err := opEqual(&stackEqual)
	if !resultEqual || err != nil || len(stackEqual) != 1 || decodeNum(stackEqual[len(stackEqual)-1]) != 1 {
		t.Errorf("opEqual failed for stack with equal elements. Unexpected state after the operation")
	}

	// Test case 4: Test when the stack has 2 or more elements, and they are not equal
	stackNotEqual := Stack{[]byte{1, 2, 3}, []byte{4, 5, 6}}
	resultNotEqual, err := opEqual(&stackNotEqual)
	if !resultNotEqual || err != nil || len(stackNotEqual) != 1 || decodeNum(stackNotEqual[len(stackEqual)-1]) != 0 {
		t.Errorf("opEqual failed for stack with non-equal elements. Unexpected state after the operation")
	}
}

func TestOpEqualVerify(t *testing.T) {
	// Test case 1: Test when opEqual and opVerify both succeed
	stackEqualVerify := Stack{[]byte{1, 2, 3}, []byte{1, 2, 3}}
	resultEqualVerify, err := opEqualVerify(&stackEqualVerify)
	if !resultEqualVerify || err != nil || len(stackEqualVerify) != 0 {
		t.Errorf("opEqualVerify failed for stack with equal elements. Unexpected state after the operation")
	}

	// Test case 2: Test when opEqual fails
	stackNotEqualVerify := Stack{[]byte{1, 2, 3}, []byte{4, 5, 6}}
	resultNotEqualVerify, err := opEqualVerify(&stackNotEqualVerify)
	if resultNotEqualVerify || err != nil {
		t.Errorf("opEqualVerify failed for stack with non-equal elements. Expected false, error nil; got true, %v", err)
	}

	// Test case 3: Test when opVerify fails
	stackEqualNoVerify := Stack{}
	resultEqualNoVerify, err := opVerify(&stackEqualNoVerify)
	if resultEqualNoVerify || err == nil || err.Error() != "stack is empty" {
		t.Errorf("opEqualVerify failed for stack with equal elements. Expected false, error 'not enough elements in stack: 2 < 1'; got true, %v", err)
	}
}

func TestOp1Add(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := op1Add(&emptyStack)
	if resultEmptyStack || err == nil || err.Error() != "stack is empty" {
		t.Errorf("op1Add failed for empty stack. Expected false, error 'not enough elements in stack: 0 < 1'; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element
	stackWithElement := Stack{[]byte{42}}
	resultWithElement, err := op1Add(&stackWithElement)
	if !resultWithElement || err != nil || len(stackWithElement) != 1 || decodeNum(stackWithElement[len(stackWithElement)-1]) != 43 {
		t.Errorf("op1Add failed for stack with at least 1 element. Unexpected state after the operation")
	}
}

func TestOp1Sub(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := op1Add(&emptyStack)
	if resultEmptyStack || err == nil || err.Error() != "stack is empty" {
		t.Errorf("op1Add failed for empty stack. Expected false, error 'not enough elements in stack: 0 < 1'; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element
	stackWithElement := Stack{[]byte{42}}
	resultWithElement, err := op1Sub(&stackWithElement)
	if !resultWithElement || err != nil || len(stackWithElement) != 1 || decodeNum(stackWithElement[len(stackWithElement)-1]) != 41 {
		t.Errorf("op1Add failed for stack with at least 1 element. Unexpected state after the operation")
	}
}

func TestOpNegate(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opNegate(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opNegate failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element
	stackWithElement := Stack{encodeNum(42)}
	resultWithElement, err := opNegate(&stackWithElement)
	if !resultWithElement || err != nil || len(stackWithElement) != 1 || decodeNum(stackWithElement[len(stackWithElement)-1]) != -42 {
		t.Errorf("opNegate failed for stack with at least 1 element. Unexpected state after the operation")
	}
}

func TestOpAbs(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opAbs(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opAbs failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element, and it is positive
	stackPositive := Stack{encodeNum(42)}
	resultPositive, err := opAbs(&stackPositive)
	if !resultPositive || err != nil || len(stackPositive) != 1 || !bytes.Equal(stackPositive[len(stackPositive)-1], encodeNum(42)) {
		t.Errorf("opAbs failed for stack with positive element. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 1 element, and it is negative
	stackNegative := Stack{encodeNum(-42)}
	resultNegative, err := opAbs(&stackNegative)
	if !resultNegative || err != nil || len(stackNegative) != 1 || !bytes.Equal(stackNegative[len(stackNegative)-1], encodeNum(42)) {
		t.Errorf("opAbs failed for stack with negative element. Unexpected state after the operation")
	}
}

func TestOpNot(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opNot(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opNot failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element, and it is 0
	stackZero := Stack{encodeNum(0)}
	resultZero, err := opNot(&stackZero)
	if !resultZero || err != nil || len(stackZero) != 1 || !bytes.Equal(stackZero[len(stackZero)-1], encodeNum(1)) {
		t.Errorf("opNot failed for stack with element 0. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 1 element, and it is non-zero
	stackNonZero := Stack{encodeNum(42)}
	resultNonZero, err := opNot(&stackNonZero)
	if !resultNonZero || err != nil || len(stackNonZero) != 1 || !bytes.Equal(stackNonZero[len(stackNonZero)-1], encodeNum(0)) {
		t.Errorf("opNot failed for stack with non-zero element. Unexpected state after the operation")
	}
}

func TestOp0NotEqual(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := op0NotEqual(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("op0NotEqual failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element, and it is 0
	stackZero := Stack{encodeNum(0)}
	resultZero, err := op0NotEqual(&stackZero)
	if !resultZero || err != nil || len(stackZero) != 1 || !bytes.Equal(stackZero[len(stackZero)-1], encodeNum(0)) {
		t.Errorf("op0NotEqual failed for stack with element 0. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 1 element, and it is non-zero
	stackNonZero := Stack{encodeNum(42)}
	resultNonZero, err := op0NotEqual(&stackNonZero)
	if !resultNonZero || err != nil || len(stackNonZero) != 1 || !bytes.Equal(stackNonZero[len(stackNonZero)-1], encodeNum(1)) {
		t.Errorf("op0NotEqual failed for stack with non-zero element. Unexpected state after the operation")
	}
}

func TestOpAdd(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opAdd(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opAdd failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements
	stackWithElements := Stack{encodeNum(42), encodeNum(13)}
	resultWithElements, err := opAdd(&stackWithElements)
	if !resultWithElements || err != nil || len(stackWithElements) != 1 || !bytes.Equal(stackWithElements[len(stackWithElements)-1], encodeNum(55)) {
		t.Errorf("opAdd failed for stack with at least 2 elements. Unexpected state after the operation")
	}
}

func TestOpSub(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opSub(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opSub failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements
	stackWithElements := Stack{encodeNum(42), encodeNum(13)}
	resultWithElements, err := opSub(&stackWithElements)
	if !resultWithElements || err != nil || len(stackWithElements) != 1 || !bytes.Equal(stackWithElements[len(stackWithElements)-1], encodeNum(29)) {
		t.Errorf("opSub failed for stack with at least 2 elements. Unexpected state after the operation")
	}
}

func TestOpMul(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opMul(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opMul failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements
	stackWithElements := Stack{encodeNum(6), encodeNum(7)}
	resultWithElements, err := opMul(&stackWithElements)
	if !resultWithElements || err != nil || len(stackWithElements) != 1 || !bytes.Equal(stackWithElements[len(stackWithElements)-1], encodeNum(42)) {
		t.Errorf("opMul failed for stack with at least 2 elements. Unexpected state after the operation")
	}
}

func TestOpBoolAnd(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opBoolAnd(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opBoolAnd failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements, both non-zero
	stackNonZero := Stack{encodeNum(1), encodeNum(42)}
	resultNonZero, err := opBoolAnd(&stackNonZero)
	if !resultNonZero || err != nil || len(stackNonZero) != 1 || !bytes.Equal(stackNonZero[len(stackNonZero)-1], encodeNum(1)) {
		t.Errorf("opBoolAnd failed for stack with non-zero elements. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, one zero
	stackWithZero := Stack{encodeNum(0), encodeNum(42)}
	resultWithZero, err := opBoolAnd(&stackWithZero)
	if !resultWithZero || err != nil || len(stackWithZero) != 1 || !bytes.Equal(stackWithZero[len(stackWithZero)-1], encodeNum(0)) {
		t.Errorf("opBoolAnd failed for stack with one zero element. Unexpected state after the operation")
	}
}

func TestOpBoolOr(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opBoolOr(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opBoolOr failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements, both non-zero
	stackNonZero := Stack{encodeNum(1), encodeNum(42)}
	resultNonZero, err := opBoolOr(&stackNonZero)
	if !resultNonZero || err != nil || len(stackNonZero) != 1 || !bytes.Equal(stackNonZero[len(stackNonZero)-1], encodeNum(1)) {
		t.Errorf("opBoolOr failed for stack with non-zero elements. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, one zero
	stackWithZero := Stack{encodeNum(0), encodeNum(42)}
	resultWithZero, err := opBoolOr(&stackWithZero)
	if !resultWithZero || err != nil || len(stackWithZero) != 1 || !bytes.Equal(stackWithZero[len(stackWithZero)-1], encodeNum(1)) {
		t.Errorf("opBoolOr failed for stack with one zero element. Unexpected state after the operation")
	}
}

func TestOpNumEqual(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opNumEqual(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opNumEqual failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements
	stackEqual := Stack{encodeNum(42), encodeNum(42)}
	resultEqual, err := opNumEqual(&stackEqual)
	if !resultEqual || err != nil || len(stackEqual) != 1 || !bytes.Equal(stackEqual[len(stackEqual)-1], encodeNum(1)) {
		t.Errorf("opNumEqual failed for stack with equal elements. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, not equal
	stackNotEqual := Stack{encodeNum(42), encodeNum(13)}
	resultNotEqual, err := opNumEqual(&stackNotEqual)
	if !resultNotEqual || err != nil || len(stackNotEqual) != 1 || !bytes.Equal(stackNotEqual[len(stackNotEqual)-1], encodeNum(0)) {
		t.Errorf("opNumEqual failed for stack with non-equal elements. Unexpected state after the operation")
	}
}

func TestOpNumEqualVerify(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opNumEqualVerify(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opNumEqualVerify failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements, equal
	stackEqual := Stack{encodeNum(42), encodeNum(42)}
	resultEqual, err := opNumEqualVerify(&stackEqual)
	if !resultEqual || err != nil || len(stackEqual) != 0 {
		t.Errorf("opNumEqualVerify failed for stack with equal elements. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, not equal
	stackNotEqual := Stack{encodeNum(42), encodeNum(13)}
	resultNotEqual, err := opNumEqualVerify(&stackNotEqual)
	if resultNotEqual || err != nil || len(stackNotEqual) != 0 {
		t.Errorf("opNumEqualVerify failed for stack with non-equal elements. Unexpected state after the operation")
	}
}

func TestOpNumNotEqual(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opNumNotEqual(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opNumNotEqual failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements
	stackEqual := Stack{encodeNum(42), encodeNum(42)}
	resultEqual, err := opNumNotEqual(&stackEqual)
	if !resultEqual || err != nil || len(stackEqual) != 1 || !bytes.Equal(stackEqual[len(stackEqual)-1], encodeNum(0)) {
		t.Errorf("opNumNotEqual failed for stack with equal elements. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, not equal
	stackNotEqual := Stack{encodeNum(42), encodeNum(13)}
	resultNotEqual, err := opNumNotEqual(&stackNotEqual)
	if !resultNotEqual || err != nil || len(stackNotEqual) != 1 || !bytes.Equal(stackNotEqual[len(stackNotEqual)-1], encodeNum(1)) {
		t.Errorf("opNumNotEqual failed for stack with non-equal elements. Unexpected state after the operation")
	}
}

func TestOpLessThan(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opLessThan(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opLessThan failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements, element2 < element1
	stackLessThan := Stack{encodeNum(13), encodeNum(42)}
	resultLessThan, err := opLessThan(&stackLessThan)
	if !resultLessThan || err != nil || len(stackLessThan) != 1 || !bytes.Equal(stackLessThan[len(stackLessThan)-1], encodeNum(1)) {
		t.Errorf("opLessThan failed for stack with element2 < element1. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, element2 >= element1
	stackNotLessThan := Stack{encodeNum(42), encodeNum(13)}
	resultNotLessThan, err := opLessThan(&stackNotLessThan)
	if !resultNotLessThan || err != nil || len(stackNotLessThan) != 1 || !bytes.Equal(stackNotLessThan[len(stackNotLessThan)-1], encodeNum(0)) {
		t.Errorf("opLessThan failed for stack with element2 >= element1. Unexpected state after the operation")
	}
}

func TestOpGreaterThan(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opGreaterThan(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opGreaterThan failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements, element2 > element1
	stackGreaterThan := Stack{encodeNum(42), encodeNum(13)}
	resultGreaterThan, err := opGreaterThan(&stackGreaterThan)
	if !resultGreaterThan || err != nil || len(stackGreaterThan) != 1 || !bytes.Equal(stackGreaterThan[len(stackGreaterThan)-1], encodeNum(1)) {
		t.Errorf("opGreaterThan failed for stack with element2 > element1. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, element2 <= element1
	stackNotGreaterThan := Stack{encodeNum(13), encodeNum(42)}
	resultNotGreaterThan, err := opGreaterThan(&stackNotGreaterThan)
	if !resultNotGreaterThan || err != nil || len(stackNotGreaterThan) != 1 || !bytes.Equal(stackNotGreaterThan[len(stackNotGreaterThan)-1], encodeNum(0)) {
		t.Errorf("opGreaterThan failed for stack with element2 <= element1. Unexpected state after the operation")
	}
}

func TestOpLessThanOrEqual(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opLessThanOrEqual(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opLessThanOrEqual failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements, element2 <= element1
	stackLessThanOrEqual := Stack{encodeNum(13), encodeNum(42)}
	resultLessThanOrEqual, err := opLessThanOrEqual(&stackLessThanOrEqual)
	if !resultLessThanOrEqual || err != nil || len(stackLessThanOrEqual) != 1 || !bytes.Equal(stackLessThanOrEqual[len(stackLessThanOrEqual)-1], encodeNum(1)) {
		t.Errorf("opLessThanOrEqual failed for stack with element2 <= element1. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, element2 > element1
	stackNotLessThanOrEqual := Stack{encodeNum(42), encodeNum(13)}
	resultNotLessThanOrEqual, err := opLessThanOrEqual(&stackNotLessThanOrEqual)
	if !resultNotLessThanOrEqual || err != nil || len(stackNotLessThanOrEqual) != 1 || !bytes.Equal(stackNotLessThanOrEqual[len(stackNotLessThanOrEqual)-1], encodeNum(0)) {
		t.Errorf("opLessThanOrEqual failed for stack with element2 > element1. Unexpected state after the operation")
	}
}

func TestOpGreaterThanOrEqual(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opGreaterThanOrEqual(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opGreaterThanOrEqual failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements, element2 >= element1
	stackGreaterThanOrEqual := Stack{encodeNum(42), encodeNum(13)}
	resultGreaterThanOrEqual, err := opGreaterThanOrEqual(&stackGreaterThanOrEqual)
	if !resultGreaterThanOrEqual || err != nil || len(stackGreaterThanOrEqual) != 1 || !bytes.Equal(stackGreaterThanOrEqual[len(stackGreaterThanOrEqual)-1], encodeNum(1)) {
		t.Errorf("opGreaterThanOrEqual failed for stack with element2 >= element1. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, element2 < element1
	stackNotGreaterThanOrEqual := Stack{encodeNum(13), encodeNum(42)}
	resultNotGreaterThanOrEqual, err := opGreaterThanOrEqual(&stackNotGreaterThanOrEqual)
	if !resultNotGreaterThanOrEqual || err != nil || len(stackNotGreaterThanOrEqual) != 1 || !bytes.Equal(stackNotGreaterThanOrEqual[len(stackNotGreaterThanOrEqual)-1], encodeNum(0)) {
		t.Errorf("opGreaterThanOrEqual failed for stack with element2 < element1. Unexpected state after the operation")
	}
}

func TestOpMin(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opMin(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opMin failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements, element1 < element2
	stackMin := Stack{encodeNum(13), encodeNum(42)}
	resultMin, err := opMin(&stackMin)
	if !resultMin || err != nil || len(stackMin) != 1 || !bytes.Equal(stackMin[len(stackMin)-1], encodeNum(13)) {
		t.Errorf("opMin failed for stack with element1 < element2. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, element1 >= element2
	stackNotMin := Stack{encodeNum(42), encodeNum(13)}
	resultNotMin, err := opMin(&stackNotMin)
	if !resultNotMin || err != nil || len(stackNotMin) != 1 || !bytes.Equal(stackNotMin[len(stackNotMin)-1], encodeNum(13)) {
		t.Errorf("opMin failed for stack with element1 >= element2. Unexpected state after the operation")
	}
}

func TestOpMax(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opMax(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opMax failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 2 elements, element1 > element2
	stackMax := Stack{encodeNum(42), encodeNum(13)}
	resultMax, err := opMax(&stackMax)
	if !resultMax || err != nil || len(stackMax) != 1 || !bytes.Equal(stackMax[len(stackMax)-1], encodeNum(42)) {
		t.Errorf("opMax failed for stack with element1 > element2. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 2 elements, element1 <= element2
	stackNotMax := Stack{encodeNum(13), encodeNum(42)}
	resultNotMax, err := opMax(&stackNotMax)
	if !resultNotMax || err != nil || len(stackNotMax) != 1 || !bytes.Equal(stackNotMax[len(stackNotMax)-1], encodeNum(42)) {
		t.Errorf("opMax failed for stack with element1 <= element2. Unexpected state after the operation")
	}
}

func TestOpWithin(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opWithin(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opWithin failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 3 elements, element inside range
	stackWithin := Stack{encodeNum(15), encodeNum(10), encodeNum(20)}
	resultWithin, err := opWithin(&stackWithin)
	if !resultWithin || err != nil || len(stackWithin) != 1 || !bytes.Equal(stackWithin[len(stackWithin)-1], encodeNum(1)) {
		t.Errorf("opWithin failed for stack with element inside range. Unexpected state after the operation")
	}

	// Test case 3: Test when the stack has at least 3 elements, element outside range
	stackNotWithin := Stack{encodeNum(25), encodeNum(10), encodeNum(20)}
	resultNotWithin, err := opWithin(&stackNotWithin)
	if !resultNotWithin || err != nil || len(stackNotWithin) != 1 || !bytes.Equal(stackNotWithin[len(stackNotWithin)-1], encodeNum(0)) {
		t.Errorf("opWithin failed for stack with element outside range. Unexpected state after the operation")
	}
}

func TestOpRipemd160(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opRipemd160(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opRipemd160 failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element
	hash := ripemd160.New()
	hash.Write([]byte("hello"))
	expected := hash.Sum(nil)
	stackRipemd160 := Stack{[]byte("hello")}
	resultRipemd160, err := opRipemd160(&stackRipemd160)
	if !resultRipemd160 || err != nil || len(stackRipemd160) != 1 || !bytes.Equal(stackRipemd160[len(stackRipemd160)-1], expected) {
		t.Errorf("opRipemd160 failed for stack with at least 1 element. Unexpected state after the operation")
	}

}

func TestOpSha1(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opSha1(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opSha1 failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element
	hash := sha1.New()
	hash.Write([]byte("hello"))
	expected := hash.Sum(nil)
	stackSha1 := Stack{[]byte("hello")}
	resultSha1, err := opSha1(&stackSha1)
	if !resultSha1 || err != nil || len(stackSha1) != 1 || !bytes.Equal(stackSha1[len(stackSha1)-1], expected) {
		t.Errorf("opSha1 failed for stack with at least 1 element. Unexpected state after the operation")
	}

}

func TestOpSha256(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opSha256(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opSha256 failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element
	hash := sha256.New()
	hash.Write([]byte("hello"))
	expected := hash.Sum(nil)
	stackSha256 := Stack{[]byte("hello")}
	resultSha256, err := opSha256(&stackSha256)
	if !resultSha256 || err != nil || len(stackSha256) != 1 || !bytes.Equal(stackSha256[len(stackSha256)-1], expected) {
		t.Errorf("opSha256 failed for stack with at least 1 element. Unexpected state after the operation")
	}
}

func TestOpHash160(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opHash160(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opHash160 failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element
	data := []byte("hello")
	hash160 := utils.Hash160(data)
	stackHash160 := Stack{data}
	resultHash160, err := opHash160(&stackHash160)
	if !resultHash160 || err != nil || len(stackHash160) != 1 || !bytes.Equal(stackHash160[len(stackHash160)-1], hash160) {
		t.Errorf("opHash160 failed for stack with at least 1 element. Unexpected state after the operation")
	}
}

func TestOpHash256(t *testing.T) {
	// Test case 1: Test when the stack is empty
	emptyStack := Stack{}
	resultEmptyStack, err := opHash256(&emptyStack)
	if resultEmptyStack || err == nil {
		t.Errorf("opHash256 failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: Test when the stack has at least 1 element
	data := []byte("hello")
	hash256 := utils.Hash256(data)
	stackHash256 := Stack{data}
	resultHash256, err := opHash256(&stackHash256)
	if !resultHash256 || err != nil || len(stackHash256) != 1 || !bytes.Equal(stackHash256[len(stackHash256)-1], hash256) {
		t.Errorf("opHash256 failed for stack with at least 1 element. Unexpected state after the operation")
	}
}

func TestOpChecksig(t *testing.T) {
	z, _ := new(big.Int).SetString("0x7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 0)
	// Test case 1: Test when the stack is empty

	emptyStack := Stack{}
	resultEmptyStack, err := opCheckSig(&emptyStack, z)
	if resultEmptyStack || err == nil {
		t.Errorf("opChecksig failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: proper Signature
	sec, _ := new(big.Int).SetString("0x04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34", 0)
	sig, _ := new(big.Int).SetString("0x3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601", 0)
	signedStack := Stack{sig.Bytes(), sec.Bytes()}

	resultSignedStack, err := opCheckSig(&signedStack, z)
	if !resultSignedStack || err != nil || !bytes.Equal(signedStack[len(signedStack)-1], encodeNum(1)) {
		t.Errorf("opChecksig failed for stack with correct Digital Signature. Unexpected state after the operation")
	}
}

func TestOpCheckMultisig(t *testing.T) {
	// doing 2-of-3 multisig
	z, _ := new(big.Int).SetString("0xe71bfa115715d6fd33796948126f40a8cdd39f187e4afb03896795189fe1423c", 0)
	sig1, _ := new(big.Int).SetString("0x3045022100dc92655fe37036f47756db8102e0d7d5e28b3beb83a8fef4f5dc0559bddfb94e02205a36d4e4e6c7fcd16658c50783e00c341609977aed3ad00937bf4ee942a8993701", 0)
	sig2, _ := new(big.Int).SetString("0x3045022100da6bee3c93766232079a01639d07fa869598749729ae323eab8eef53577d611b02207bef15429dcadce2121ea07f233115c6f09034c0be68db99980b9a6c5e75402201", 0)
	sec1, _ := new(big.Int).SetString("0x022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb70", 0)
	sec2, _ := new(big.Int).SetString("0x03b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb71", 0)
	sec3, _ := new(big.Int).SetString("0x04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34", 0)

	stack := Stack{
		encodeNum(0),
		sig1.Bytes(),
		sig2.Bytes(),
		encodeNum(2),
		sec1.Bytes(),
		sec2.Bytes(),
		sec3.Bytes(),
		encodeNum(3),
	}

	result, err := opCheckMultiSig(&stack, z)
	if !result || err != nil {
		t.Errorf("opCheckMultisig failed. Expected true, nil; got %v, %v", result, err)
	}

	if decodeNum(stack[0]) != 1 {
		t.Errorf("Decoded number mismatch. Expected 1, got %v", decodeNum(stack[0]))
	}
}

func TestOpChecksigVerify(t *testing.T) {
	z, _ := new(big.Int).SetString("0x7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 0)
	// Test case 1: Test when the stack is empty

	emptyStack := Stack{}
	resultEmptyStack, err := opCheckSigVerify(&emptyStack, z)
	if resultEmptyStack || err == nil {
		t.Errorf("opChecksigVerify failed for empty stack. Expected false, nil; got true, %v", err)
	}

	// Test case 2: proper Signature
	sec, _ := new(big.Int).SetString("0x04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34", 0)
	sig, _ := new(big.Int).SetString("0x3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601", 0)
	signedStack := Stack{sig.Bytes(), sec.Bytes()}

	resultSignedStack, err := opCheckSigVerify(&signedStack, z)
	if !resultSignedStack || err != nil {
		t.Errorf("opChecksigVerify failed for stack with correct Digital Signature. Unexpected state after the operation")
	}
}

func TestOpCheckLockTimeVerify(t *testing.T) {
	// Test case 1: Invalid sequence value
	stack := Stack{encodeNum(100)}
	locktime := 123
	sequence := 0xffffffff
	result, err := opCheckLockTimeVerify(&stack, locktime, sequence)
	if result || err == nil || err.Error() != "invalid sequence value" {
		t.Errorf("opCheckLockTimeVerify failed for invalid sequence value. Expected false, 'invalid sequence value'; got true, %v", err)
	}

	// Test case 2: Insufficient elements in stack
	emptyStack := Stack{}
	sequence = 0xfffffffe
	result, err = opCheckLockTimeVerify(&emptyStack, locktime, sequence)
	if result || err == nil || err.Error() != "stack is empty" {
		t.Errorf("opCheckLockTimeVerify failed for insufficient elements in stack. Expected false, 'insufficient elements in stack'; got true, %v", err)
	}

	// Test case 3: Negative element in stack
	negativeElementStack := Stack{encodeNum(-100)}
	result, err = opCheckLockTimeVerify(&negativeElementStack, locktime, sequence)
	if result || err == nil || err.Error() != "negative element in stack" {
		t.Errorf("opCheckLockTimeVerify failed for negative element in stack. Expected false, 'negative element in stack'; got true, %v", err)
	}

	// Test case 4: Locktime exceeds 500000000 for element less than 500000000
	stack = Stack{encodeNum(400000000)}
	locktime = 600000000
	result, err = opCheckLockTimeVerify(&stack, locktime, sequence)
	if result || err == nil || err.Error() != "locktime exceeds 500000000 for element less than 500000000" {
		t.Errorf("opCheckLockTimeVerify failed for locktime exceeds 500000000 for element less than 500000000. Expected false, 'locktime exceeds 500000000 for element less than 500000000'; got true, %v", err)
	}

	// Test case 5: Locktime is less than element in stack
	stack = Stack{encodeNum(600000000)}
	locktime = 400000000
	result, err = opCheckLockTimeVerify(&stack, locktime, sequence)
	if result || err == nil || err.Error() != "locktime is less than element in stack" {
		t.Errorf("opCheckLockTimeVerify failed for locktime is less than element in stack. Expected false, 'locktime is less than element in stack'; got true, %v", err)
	}

	// Test case 6: Valid case
	stack = Stack{encodeNum(122)}
	locktime = 123
	result, err = opCheckLockTimeVerify(&stack, locktime, sequence)
	if !result || err != nil {
		t.Errorf("opCheckLockTimeVerify failed for valid case. Expected true, nil; got %v, %v", result, err)
	}
}

func TestOpCheckSequenceVerify(t *testing.T) {
	// Test case 1: Invalid sequence value
	stack := Stack{encodeNum(100)}
	version := 1
	sequence := 0xffffffff
	result, err := opCheckSequenceVerify(&stack, version, sequence)
	if result || err == nil || err.Error() != "invalid sequence value" {
		t.Errorf("opCheckSequenceVerify failed for invalid sequence value. Expected false, 'invalid sequence value'; got true, %v", err)
	}

	// Test case 2: Insufficient elements in stack
	emptyStack := Stack{}
	sequence = 0x7FFFFFFF
	result, err = opCheckSequenceVerify(&emptyStack, version, sequence)
	if result || err == nil || err.Error() != "stack is empty" {
		t.Errorf("opCheckSequenceVerify failed for insufficient elements in stack. Expected false, 'insufficient elements in stack'; got true, %v", err)
	}

	// Test case 3: Negative element in stack
	negativeElementStack := Stack{encodeNum(-100)}
	result, err = opCheckSequenceVerify(&negativeElementStack, version, sequence)
	if result || err == nil || err.Error() != "negative element in stack" {
		t.Errorf("opCheckSequenceVerify failed for negative element in stack. Expected false, 'negative element in stack'; got true, %v", err)
	}

	// Test case 4: Version is less than 2 for sequence with sign bit set
	stack = Stack{encodeNum(0xC0000000)}
	version = 1
	sequence = 0x40000000
	result, err = opCheckSequenceVerify(&stack, version, sequence)
	if result || err == nil || err.Error() != "version is less than 2 for sequence with sign bit set" {
		t.Errorf("opCheckSequenceVerify failed for version is less than 2 for sequence with sign bit set. Expected false, 'version is less than 2 for sequence with sign bit set'; got true, %v", err)
	}

	// Test case 5: Invalid sequence value
	version = 2
	sequence = 0xffffffff
	result, err = opCheckSequenceVerify(&stack, version, sequence)
	if result || err == nil || err.Error() != "invalid sequence value" {
		t.Errorf("opCheckSequenceVerify failed for invalid sequence value. Expected false, 'invalid sequence value'; got true, %v", err)
	}

	// Test case 6: Mismatch in bits 22-31 between element and sequence
	stack = Stack{encodeNum(0b11111111111111111110000000000000)}
	version = 2
	sequence = 0b00000000000000000001111111111111
	result, err = opCheckSequenceVerify(&stack, version, sequence)
	if result || err == nil || err.Error() != "mismatch in bits 22-31 between element and sequence" {
		t.Errorf("opCheckSequenceVerify failed for mismatch in bits 22-31 between element and sequence. Expected false, 'mismatch in bits 22-31 between element and sequence'; got true, %v", err)
	}

	// Test case 7: Sequence value is less than element in stack
	stack = Stack{encodeNum(0b11111111111111111111111111111111)}
	version = 2
	sequence = 0b01111111111111110000000000000000
	result, err = opCheckSequenceVerify(&stack, version, sequence)
	if result || err == nil || err.Error() != "sequence value is less than element in stack" {
		t.Errorf("opCheckSequenceVerify failed for sequence value is less than element in stack. Expected false, 'sequence value is less than element in stack'; got true, %v", err)
	}

	// Test case 8: Valid case
	stack = Stack{encodeNum(0)}
	version = 2
	sequence = 0x0020ffff
	result, err = opCheckSequenceVerify(&stack, version, sequence)
	if !result || err != nil {
		t.Errorf("opCheckSequenceVerify failed for valid case. Expected true, nil; got %v, %v", result, err)
	}
}

func performOperation(op func(*Stack) (bool, error), stack *Stack, expected int, t *testing.T) {
	op(stack)
	result := decodeNum((*stack)[len(*stack)-1])

	if result != expected {
		t.Errorf("Failed for %s. Expected %d, got %d", getOpName(op), expected, result)
	}
}

func getOpName(op interface{}) string {
	return fmt.Sprintf("%p", op)
}
