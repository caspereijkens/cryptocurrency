package script

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestNewScript(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected Script
		wantErr  bool
	}{
		{
			name:     "Empty data",
			input:    []byte{},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Valid script with just chars",
			input:    []byte{0x04, 't', 'e', 's', 't'},
			expected: Script{[]byte{'t'}, []byte{'e'}, []byte{'s'}, []byte{'t'}},
			wantErr:  false,
		},
		{
			name:     "Valid script with OP_PUSHDATA1",
			input:    []byte{0x04, 0x4C, 0x04, 't', 'e', 's', 't'},
			expected: Script{[]byte{'t', 'e', 's', 't'}},
			wantErr:  false,
		},
		{
			name:     "Valid script with OP_PUSHDATA2",
			input:    []byte{0x04, 0x4D, 0x02, 0x00, 'a', 'b'},
			expected: Script{[]byte{'a', 'b'}},
			wantErr:  false,
		},
		{
			name:     "Valid script with OP_PUSHDATA2 (case showcasing pushdata2)",
			input:    []byte{0x06, 0x4D, 0x02, 0x00, 'c', 'd', 'e'},
			expected: Script{[]byte{'c', 'd'}, []byte{'e'}},
			wantErr:  false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := NewScript(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(script, tt.expected) {
				t.Errorf("NewScript() got = %v, want %v", script, tt.expected)
			}
		})
	}
}

func TestScriptParsing(t *testing.T) {
	scriptPubKeyHex := "6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	scriptPubKey, _ := hex.DecodeString(scriptPubKeyHex)

	script, err := NewScript(scriptPubKey)
	if err != nil {
		t.Fatalf("NewScript() error: %v", err)
	}

	wantCmd1, _ := hex.DecodeString("304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a71601")
	if !bytes.Equal(script[0], wantCmd1) {
		t.Errorf("Script.Parse() cmds[0] = %x, want %x", script[0], wantCmd1)
	}

	wantCmd2, _ := hex.DecodeString("035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937")
	if !bytes.Equal(script[1], wantCmd2) {
		t.Errorf("Script.Parse() cmds[1] = %x, want %x", script[1], wantCmd2)
	}
}

func TestSerialize(t *testing.T) {
	want := "6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	wantBytes, _ := hex.DecodeString(want)

	var script Script
	err := script.Parse(wantBytes)
	if err != nil {
		t.Errorf("Failed to parse script: %v", err)
		return
	}

	serialized, err := script.Serialize()
	if err != nil {
		t.Errorf("Failed to serialize script: %v", err)
		return
	}

	if !bytes.Equal(serialized, wantBytes) {
		t.Errorf("Serialized result does not match. Got: %x, Want: %s", serialized, want)
	}
}
