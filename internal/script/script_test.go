package script

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"
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
			input:    []byte{0x06, 0x4C, 0x04, 't', 'e', 's', 't'},
			expected: Script{[]byte{'t', 'e', 's', 't'}},
			wantErr:  false,
		},
		{
			name:     "Valid script with OP_PUSHDATA2",
			input:    []byte{0x05, 0x4D, 0x02, 0x00, 'a', 'b'},
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
			script, err := NewScript(bufio.NewReader(bytes.NewBuffer(tt.input)))

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
	scriptPubKeyBytes, _ := hex.DecodeString(scriptPubKeyHex)
	scriptPubKey := bufio.NewReader(bytes.NewBuffer(scriptPubKeyBytes))

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

	scriptPubKeyHex = "6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a"
	scriptPubKeyBytes, _ = hex.DecodeString(scriptPubKeyHex)
	scriptPubKey = bufio.NewReader(bytes.NewBuffer(scriptPubKeyBytes))

	script, err = NewScript(scriptPubKey)
	if err != nil {
		t.Fatalf("NewScript() error: %v", err)
	}

	wantCmd1, _ = hex.DecodeString("3045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01")
	if !bytes.Equal(script[0], wantCmd1) {
		t.Errorf("Script.Parse() cmds[0] = %x, want %x", script[0], wantCmd1)
	}

	wantCmd2, _ = hex.DecodeString("0349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")
	if !bytes.Equal(script[1], wantCmd2) {
		t.Errorf("Script.Parse() cmds[1] = %x, want %x", script[1], wantCmd2)
	}
}

func TestSerialize(t *testing.T) {
	want := "6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	wantBytes, _ := hex.DecodeString(want)

	var script Script
	err := script.Parse(bufio.NewReader(bytes.NewBuffer(wantBytes)))
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

// Now a bunch of tests where I try the standard scripts from the book.

func TestPayToPubKeyExample(t *testing.T) {
	z, _ := new(big.Int).SetString("0x7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 0)
	sec, _ := hex.DecodeString("04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	sig, _ := hex.DecodeString("3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601")
	pubkeyScript := Script{sec, []byte{0xac}}
	sigScript := Script{sig}
	combinedScript := sigScript.Add(pubkeyScript)
	if ok := combinedScript.Evaluate(z); !ok {
		t.Errorf("Combined script does not match. Evalutation resulted in False")
	}

	// Just to prove that it really works, on purpose faulty signature:
	falseSig, _ := hex.DecodeString("3045022000eaa69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601")
	falseSigScript := Script{falseSig}
	combinedScript = falseSigScript.Add(pubkeyScript)
	if ok := combinedScript.Evaluate(z); ok {
		t.Errorf("Combined script should have failed. Evalutation resulted in True")
	}

}

func TestSomeArbitraryPrograms(t *testing.T) {
	// 4 + 5 = 9
	pubkeyScript1 := Script{[]byte{0x55}, []byte{0x93}, []byte{0x59}, []byte{0x87}}
	sigScript1 := Script{[]byte{0x54}}
	combinedScript1 := sigScript1.Add(pubkeyScript1)
	if ok := combinedScript1.Evaluate(nil); !ok {
		t.Errorf("Combined script does not match. Evalutation resulted in False")
	}

	// 2 + 2^2 = 6
	pubkeyScript2 := Script{[]byte{0x76}, []byte{0x76}, []byte{0x95}, []byte{0x93}, []byte{0x56}, []byte{0x87}}
	sigScript2 := Script{[]byte{0x52}}
	combinedScript2 := sigScript2.Add(pubkeyScript2)
	if ok := combinedScript2.Evaluate(nil); !ok {
		t.Errorf("Combined script does not match. Evalutation resulted in False")
	}
}

// This test proves a SHA-1 hash collision found in February 2017 https://security.googleblog.com/2017/02/announcing-first-sha1-collision.html
func TestSha1HashPinata(t *testing.T) {
	filePath1 := "resources/shattered-1.pdf"

	// Open the file1
	file1, err := os.Open(filePath1)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file1.Close()

	// Read the first 320 bytes
	buffer1 := make([]byte, 320)
	_, err = io.ReadFull(file1, buffer1)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		fmt.Println("Error reading file:", err)
		return
	}

	filePath2 := "resources/shattered-2.pdf"

	// Open the file2
	file2, err := os.Open(filePath2)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file2.Close()

	// Read the first 320 bytes
	buffer2 := make([]byte, 320)
	_, err = io.ReadFull(file2, buffer2)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		fmt.Println("Error reading file:", err)
		return
	}

	// Create a Script by feeding the two byte-slices to the hash Pinata
	pubkeyScriptHashPinataBytes, _ := hex.DecodeString("086e879169a77ca787")
	pubkeyScriptHashPinataBuf := bufio.NewReader(bytes.NewBuffer(pubkeyScriptHashPinataBytes))
	pubkeyScriptHashPinata, _ := NewScript(pubkeyScriptHashPinataBuf)
	sigScriptHashPinata := Script{buffer1, buffer2}
	combinedScriptHashPinata := sigScriptHashPinata.Add(pubkeyScriptHashPinata)
	if ok := combinedScriptHashPinata.Evaluate(nil); !ok {
		t.Errorf("Combined script does not match. Evalutation resulted in False")
	}
}
