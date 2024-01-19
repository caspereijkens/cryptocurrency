package transaction

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
)

const (
	cacheFile = "./resources/tx.cache"
)

var (
	txFetcher = NewTxFetcher()
	fresh     = false
	testnet   = true
)

func init() {
	txFetcher.LoadCache(cacheFile)
}

func TestParseVersion(t *testing.T) {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	stream := bytes.NewReader(rawTx)
	tx, _ := ParseTx(bufio.NewReader(stream), false)
	if tx.Version != 1 {
		t.Errorf("Expected version 1, got %d", tx.Version)
	}
}

func TestParseInputs(t *testing.T) {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	stream := bytes.NewReader(rawTx)
	tx, _ := ParseTx(bufio.NewReader(stream), false)
	if len(tx.TxIns) != 1 {
		t.Errorf("Expected 1 input, got %d", len(tx.TxIns))
	}
	want, _ := hex.DecodeString("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81")
	if !bytes.Equal(tx.TxIns[0].PrevTx, want) {
		t.Errorf("Expected PrevTx %x, got %x", want, tx.TxIns[0].PrevTx)
	}
	if tx.TxIns[0].PrevIndex != 0 {
		t.Errorf("Expected PrevIndex 0, got %d", tx.TxIns[0].PrevIndex)
	}
	want, _ = hex.DecodeString("6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")
	have, err := tx.TxIns[0].ScriptSig.Serialize()
	if err != nil {
		t.Errorf("Error serializing first transaction input: %v", err)
	}
	if !bytes.Equal(have, want) {
		t.Errorf("Expected ScriptSig %x, got %x", want, have)
	}
	if tx.TxIns[0].Sequence != 0xfffffffe {
		t.Errorf("Expected Sequence 0xfffffffe, got %d", tx.TxIns[0].Sequence)
	}
}

func TestParseOutputs(t *testing.T) {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	stream := bytes.NewReader(rawTx)
	tx, _ := ParseTx(bufio.NewReader(stream), false)
	if len(tx.TxOuts) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(tx.TxOuts))
	}
	if tx.TxOuts[0].Amount != 32454049 {
		t.Errorf("Expected Amount 32454049, got %d", tx.TxOuts[0].Amount)
	}
	want, _ := hex.DecodeString("1976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac")
	have, err := tx.TxOuts[0].ScriptPubkey.Serialize()
	if err != nil {
		t.Errorf("Error serializing first transaction input: %v", err)
	}
	if !bytes.Equal(have, want) {
		t.Errorf("Expected ScriptPubkey %x, got %x", want, have)
	}
	if tx.TxOuts[1].Amount != 10011545 {
		t.Errorf("Expected Amount 10011545, got %d", tx.TxOuts[1].Amount)
	}
	want, _ = hex.DecodeString("1976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac")
	have, err = tx.TxOuts[1].ScriptPubkey.Serialize()
	if err != nil {
		t.Errorf("Error serializing first transaction input: %v", err)
	}
	if !bytes.Equal(have, want) {
		t.Errorf("Expected ScriptPubkey %x, got %x", want, have)
	}
}

func TestParseLocktime(t *testing.T) {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	stream := bytes.NewReader(rawTx)
	tx, _ := ParseTx(bufio.NewReader(stream), false)
	if tx.Locktime != 410393 {
		t.Errorf("Expected Locktime 410393, got %d", tx.Locktime)
	}
}

func TestTxId(t *testing.T) {
	expectedId := "0d6fe5213c0b3291f208cba8bfb59b7476dffacc4e5cb66f6eb20a080843a299"
	tx, err := txFetcher.Fetch(expectedId, testnet, fresh)
	if err != nil {
		t.Errorf("Error loading tx: %v", err)
	}

	id, err := tx.Id()
	if err != nil {
		t.Errorf("Error generating id of tx: %v", fresh)
	}

	if id != expectedId {
		t.Errorf("Expected id and generated id do not match.\nwant: %s\ngot: %s", expectedId, id)
	}
}

func TestTxFee(t *testing.T) {
	id := "184d3393cea44574a7b521575878a5485fc3c18e4920808235c8f58264c1dc48"
	tx, err := txFetcher.Fetch(id, testnet, fresh)
	if err != nil {
		t.Errorf("Error loading tx: %v", err)
	}

	fee, err := tx.Fee()
	if err != nil {
		t.Errorf("Error calculating fee: %v", err)
	}
	expectedFee := uint64(534528)
	if fee != expectedFee {
		t.Errorf("Error calculating fee:\nwant: %d\nhave: %d", fee, expectedFee)
	}
}

func TestSigHash(t *testing.T) {
	testnet = false
	id := "452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03"
	tx, err := NewTxFetcher().Fetch(id, testnet, fresh)
	if err != nil {
		t.Fatalf("Failed to fetch transaction: %v", err)
	}

	want, _ := new(big.Int).SetString("0x27e0c5994dec7824e56dec6b2fcb342eb7cdb0d0957c2fce9882f715e85d81a6", 0)

	result, err := tx.SigHash(0)
	if err != nil {
		t.Fatalf("Error calling SigHash: %v", err)
	}

	if result.Cmp(want) != 0 {
		t.Errorf("SigHash result mismatch, got: %s, want: %s", result.Text(16), want.Text(16))
	}
}

func TestTxInValue(t *testing.T) {
	expectedValue := uint64(250000000)
	testnet = false
	id := "42f7d0545ef45bd3b9cfee6b170cf6314a3bd8b3f09b610eeb436d92993ad440"
	tx, err := txFetcher.Fetch(id, testnet, fresh)
	if err != nil {
		t.Errorf("Error loading tx: %v", err)
	}

	txIn0 := tx.TxIns[0]
	value, err := txIn0.Value(testnet)
	if err != nil {
		t.Errorf("Error calculating value: %v", err)
	}

	if value != expectedValue {
		t.Errorf("Value of input is wrong.\nExpected:%d\nGot:%d", expectedValue, value)
	}
}
