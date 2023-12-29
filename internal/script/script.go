package script

import (
	"encoding/binary"
	"fmt"
)

type Script [][]byte

// NewScript creates a new Script from a byte slice.
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

		if currentByte >= 1 && currentByte <= 75 {
			// For a number between 1 and 75 inclusive, we know the next n bytes are an element.
			n := int(currentByte)
			script = append(script, data[count:count+n])
			count += n
		} else if currentByte == 76 {
			// 76 is OP_PUSHDATA1, so the next byte tels us how many bytes to read.
			dataLength := int(data[count])
			count++
			script = append(script, data[count:count+dataLength])
			count += dataLength
		} else if currentByte == 77 {
			// 77 is OP_PUSHDATA2, so the next two bytes tell us how many bytes to read.
			dataLength := binary.LittleEndian.Uint16(data[count : count+2])
			count += 2
			script = append(script, data[count:count+int(dataLength)])
			count += int(dataLength)
		} else {
			script = append(script, []byte{currentByte})
		}
	}

	if count != len(data) {
		return nil, fmt.Errorf("parsing script failed")
	}

	return script, nil
}
