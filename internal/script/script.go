package script

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"reflect"

	"github.com/caspereijkens/cryptocurrency/internal/utils"
)

type Script [][]byte

// ParseScript creates a new Script from a byte slice.
// OP_PUSHDATA1/2 can be used to group data in a []byte.
func ParseScript(reader *bufio.Reader) (Script, error) {
	length, err := utils.ReadVarint(reader)

	if err != nil {
		return nil, fmt.Errorf("no uvarint could be read from reader: %v", err)
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read script data: %v", err)
	}

	script := make(Script, 0)
	count := 0

	for count < int(length) {
		currentByte := buf[count]
		count++

		switch {
		case currentByte >= 1 && currentByte <= 75:
			// For a number between 1 and 75 inclusive, the next n bytes are an element.
			n := int(currentByte)
			script = append(script, buf[count:count+n])
			count += n
		case currentByte == 76:
			// 76 is OP_PUSHDATA1, so the next byte tells us how many bytes to read.
			bufLength := int(buf[count])
			count++
			script = append(script, buf[count:count+bufLength])
			count += bufLength
		case currentByte == 77:
			// 77 is OP_PUSHDATA2, so the next two bytes tell us how many bytes to read.
			bufLength := binary.LittleEndian.Uint16(buf[count : count+2])
			count += 2
			script = append(script, buf[count:count+int(bufLength)])
			count += int(bufLength)
		default:
			script = append(script, []byte{currentByte})
		}
	}

	if count != len(buf) {
		return nil, fmt.Errorf("parsing script failed")
	}

	return script, nil
}

func (s *Script) String() string {
	var result []string
	for _, cmd := range *s {
		if len(cmd) == 1 {
			opCode := int(cmd[0])
			name, ok := opCodeNames[opCode]
			if !ok {
				result = append(result, fmt.Sprintf("OP_[%d]", opCode))
				continue
			}
			result = append(result, name)
			continue
		}
		result = append(result, fmt.Sprintf("%x", cmd))
	}
	return " " + fmt.Sprintf("%v", result)
}

func (s *Script) Add(otherScript Script) Script {
	return append(*s, otherScript...)
}

func (s *Script) Parse(reader *bufio.Reader) error {
	script, err := ParseScript(reader)
	if err != nil {
		return err
	}
	*s = script
	return nil
}

func (s *Script) rawSerialize() ([]byte, error) {
	var result []byte

	for _, cmd := range *s {
		length := len(cmd)
		switch {
		case len(cmd) == 1:
			// if the command is an integer, we know it's an op code
			result = append(result, cmd...)
		case length < 75:
			// if the length is between 1 and 75, we encode the length as a single byte
			result = append(result, byte(length))
			result = append(result, cmd...)
		case length > 75 && length < 0x100:
			// For any element with length 76 to 255, we put OP_PUSHDATA1 first, then encode the length as a single byte, followed by the element.
			result = append(result, 76)
			result = append(result, byte(length))
			result = append(result, cmd...)
		case length >= 0x100 && length <= 520:
			// For any element with length 256 to 520, we put OP_PUSHDATA2 first, then encode the length as two bytes, followed by the element.
			result = append(result, 77)
			binary.LittleEndian.PutUint16(result[len(result):], uint16(length))
			result = append(result, cmd...)
		default:
			return nil, fmt.Errorf("too long a cmd")
		}
	}
	return result, nil
}

// serialize serializes the Script and adds the total length prefix.
func (s *Script) Serialize() ([]byte, error) {
	rawResult, err := s.rawSerialize()
	if err != nil {
		return nil, err
	}

	// Get the varint bytes
	varint, err := utils.EncodeVarint(uint64(len(rawResult)))
	if err != nil {
		return nil, err
	}

	// Append the varint and the serialized script
	result := append(varint, rawResult...)

	return result, nil
}

func (s *Script) Evaluate(z *big.Int) bool {
	cmds := make(Script, len(*s))
	copy(cmds, *s)

	var stack Stack
	var altStack Stack

	for len(cmds) > 0 {
		cmd := cmds[0]
		cmds = cmds[1:]

		if len(cmd) == 1 {
			opCode := int(cmd[0])

			operation := OpCodeFunctions[opCode]
			opName := opCodeNames[opCode]

			switch opCode {
			case 99, 100:
				ok, err := callOperation(operation, &stack, cmds)
				if !ok || err != nil {
					fmt.Printf("bad op: '%s', error: %v\n", opName, err)
					return false
				}
			case 107, 108:
				ok, err := callOperation(operation, &stack, &altStack)
				if !ok || err != nil {
					fmt.Printf("bad op: '%s', error: %v\n", opName, err)
					return false
				}
			case 172, 173, 174, 175:
				ok, err := callOperation(operation, &stack, z)
				if !ok || err != nil {
					fmt.Printf("bad op: '%s', error: %v\n", opName, err)
					return false
				}
			default:
				ok, err := callOperation(operation, &stack)
				if !ok || err != nil {
					fmt.Printf("bad op: '%s', error: %v\n", opName, err)
					return false
				}
			}
		} else {
			stack.push(cmd)
		}
	}

	if len(stack) == 0 || string((stack)[len(stack)-1]) == "" {
		return false
	}

	return true
}

func callOperation(fn interface{}, args ...interface{}) (bool, error) {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return false, fmt.Errorf("not a function")
	}

	// Prepare the arguments
	var input []reflect.Value
	for _, arg := range args {
		input = append(input, reflect.ValueOf(arg))
	}

	// Call the function
	result := v.Call(input)

	// Extract the return values
	if len(result) != 2 {
		// Assuming the first return value is bool and the second is error
		return false, fmt.Errorf("function did not return expected values")
	}

	if result[1].Interface() != nil {
		return result[0].Bool(), result[1].Interface().(error)
	}

	return result[0].Bool(), nil
}

func (s *Script) TranslateToOps() []string {
	ops := make([]string, len(*s))
	for i, cmd := range *s {
		ops[i] = opCodeNames[int(cmd[0])]
	}
	return ops
}

// Takes a hash160 and returns the p2pkh ScriptPubKey
func CreateP2pkhScript(h160 []byte) Script {
	return Script{[]byte{0x76}, []byte{0xa9}, h160, []byte{0x88}, []byte{0xac}}
}
