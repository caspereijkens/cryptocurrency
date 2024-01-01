package script

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

type Script [][]byte

// NewScript creates a new Script from a byte slice.
// OP_PUSHDATA1/2 can be used to group data in a single []byte.
func NewScript(data []byte) (Script, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty script data")
	}

	length, varintSize := binary.Uvarint(data)
	if length <= 0 {
		return nil, fmt.Errorf("failed to read script length")
	}

	data = data[varintSize:]

	script := make(Script, 0)
	count := 0

	for count < int(length) {
		currentByte := data[count]
		count++

		switch {
		case currentByte >= 1 && currentByte <= 75:
			// For a number between 1 and 75 inclusive, the next n bytes are an element.
			n := int(currentByte)
			script = append(script, data[count:count+n])
			count += n
		case currentByte == 76:
			// 76 is OP_PUSHDATA1, so the next byte tells us how many bytes to read.
			dataLength := int(data[count])
			count++
			script = append(script, data[count:count+dataLength])
			count += dataLength
		case currentByte == 77:
			// 77 is OP_PUSHDATA2, so the next two bytes tell us how many bytes to read.
			dataLength := binary.LittleEndian.Uint16(data[count : count+2])
			count += 2
			script = append(script, data[count:count+int(dataLength)])
			count += int(dataLength)
		default:
			script = append(script, []byte{currentByte})
		}
	}

	if count != len(data) {
		return nil, fmt.Errorf("parsing script failed")
	}

	return script, nil
}

func (s *Script) Parse(data []byte) error {
	script, err := NewScript(data)
	if err != nil {
		return err
	}
	*s = script
	return nil
}

func (s Script) rawSerialize() ([]byte, error) {
	var result []byte

	for _, cmd := range s {
		length := len(cmd)
		switch {
		case isInteger(cmd):
			result = append(result, cmd...)
		case length < 75:
			result = append(result, byte(length))
			result = append(result, cmd...)
		case length > 75 && length < 0x100:
			result = append(result, 76)
			result = append(result, byte(length))
			result = append(result, cmd...)
		case length >= 0x100 && length <= 520:
			result = append(result, 77)
			binary.LittleEndian.PutUint16(result[len(result):], uint16(length))
			result = append(result, cmd...)
		default:
			return nil, errors.New("too long a cmd")
		}
	}
	return result, nil
}

func isInteger(data []byte) bool {
	// Convert []byte to string
	strData := string(data)

	// Parse the string as an integer with base 10 and bit size 64
	_, err := strconv.ParseInt(strData, 10, 8)

	// If there is no error, it is a valid integer
	return err == nil
}

// serialize serializes the Script and adds the total length prefix.
func (s Script) Serialize() ([]byte, error) {
	rawResult, err := s.rawSerialize()
	if err != nil {
		return nil, err
	}

	// Get the varint bytes
	varint := make([]byte, binary.MaxVarintLen64)
	length := binary.PutUvarint(varint, uint64(len(rawResult)))

	// Append the varint and the serialized script
	result := append(varint[:length], rawResult...)

	return result, nil
}
