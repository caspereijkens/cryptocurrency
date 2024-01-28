package transaction

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/caspereijkens/cryptocurrency/internal/script"
	"github.com/caspereijkens/cryptocurrency/internal/signatureverification"
	"github.com/caspereijkens/cryptocurrency/internal/utils"
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

func TestTxSigHash(t *testing.T) {
	testnet = false
	id := "452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03"
	tx, err := NewTxFetcher().Fetch(id, testnet, fresh)
	if err != nil {
		t.Fatalf("Failed to fetch transaction: %v", err)
	}

	want, _ := new(big.Int).SetString("0x27e0c5994dec7824e56dec6b2fcb342eb7cdb0d0957c2fce9882f715e85d81a6", 0)

	result, err := tx.SigHash(0, nil)
	if err != nil {
		t.Fatalf("Error calling SigHash: %v", err)
	}

	if result.Cmp(want) != 0 {
		t.Errorf("SigHash result mismatch, got: %s, want: %s", result.Text(16), want.Text(16))
	}
}

func TestTxVerifyP2PKH(t *testing.T) {
	testnet = false
	// Test case 1
	tx1, err := NewTxFetcher().Fetch("452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03", testnet, fresh)
	if err != nil {
		t.Fatalf("Error fetching transaction: %v", err)
	}
	if !tx1.Verify() {
		t.Errorf("Verification failed for transaction 1")
	}

	// Test case 2
	testnet = true
	tx2, err := NewTxFetcher().Fetch("5418099cc755cb9dd3ebc6cf1a7888ad53a1a3beb5a025bce89eb1bf7f1650a2", testnet, fresh)
	if err != nil {
		t.Fatalf("Error fetching transaction: %v", err)
	}
	if !tx2.Verify() {
		t.Errorf("Verification failed for transaction 2")
	}
}

func TestVerifyP2SH(t *testing.T) {
	testnet = false
	// Test case
	tx, err := NewTxFetcher().Fetch("46df1a9484d0a81d03ce0ee543ab6e1a23ed06175c104a178268fad381216c2b", testnet, fresh)
	if err != nil {
		t.Fatalf("Error fetching transaction: %v", err)
	}
	if !tx.Verify() {
		t.Errorf("Verification failed for transaction")
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

func TestTxInPubkey(t *testing.T) {
	txHash := "d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81"
	index := uint32(0)
	wantHex := "1976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88ac"

	txInBytes, err := hex.DecodeString(txHash)
	if err != nil {
		t.Fatalf("Error fetching ScriptPubkey: %v", err)
	}
	txIn := NewTxIn(txInBytes, index, &script.Script{}, uint32(0xffffffff))
	scriptPubkey, err := txIn.ScriptPubkey(false)
	if err != nil {
		t.Fatalf("Error fetching ScriptPubkey: %v", err)
	}

	want, err := hex.DecodeString(wantHex)
	if err != nil {
		t.Fatalf("Error decoding expected ScriptPubkey: %v", err)
	}

	have, err := scriptPubkey.Serialize()
	if err != nil {
		t.Fatalf("Error decoding expected ScriptPubkey: %v", err)
	}
	if !bytes.Equal(have, want) {
		t.Errorf("ScriptPubkey mismatch. Got %x, want %x", have, want)
	}
}

func TestCreateAndSignTransaction(t *testing.T) {
	expectedHex := "010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d0000006b4830450221008ed46aa2cf12d6d81065bfabe903670165b538f65ee9a3385e6327d80c66d3b502203124f804410527497329ec4715e18558082d489b218677bd029e7fa306a72236012103935581e52c354cd2f484fe8ed83af7a3097005b2f9c60bff71d35bd795f54b67ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000"
	prevTx, _ := hex.DecodeString("0d6fe5213c0b3291f208cba8bfb59b7476dffacc4e5cb66f6eb20a080843a299")
	prevIndex := uint32(13)
	txIn := NewTxIn(prevTx, prevIndex, &script.Script{}, uint32(0xffffffff))
	changeAmount := uint64(0.33 * 100000000)
	changeH160, _ := utils.DecodeBase58("mzx5YhAH9kNHtcN481u6WkjeHjYtVeKVh2")
	changeScript := script.CreateP2pkhScript(changeH160)
	changeOutput := NewTxOut(changeAmount, changeScript)
	targetAmount := uint64(0.1 * 100000000)
	targetH160, _ := utils.DecodeBase58("mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf")
	targetScript := script.CreateP2pkhScript(targetH160)
	targetOutput := NewTxOut(targetAmount, targetScript)
	tx := NewTx(1, []*TxIn{txIn}, []*TxOut{changeOutput, targetOutput}, 0, true)
	inputIndex := uint32(0)
	z, err := tx.SigHash(inputIndex, nil)
	if err != nil {
		t.Fatalf("Failed to compute message 'z': %v", err)
	}
	privateKey, err := signatureverification.NewPrivateKey(big.NewInt(8675309))
	if err != nil {
		t.Fatalf("Failed to create new random private key: %v", err)
	}
	signature, err := privateKey.Sign(z)
	if err != nil {
		t.Fatalf("Failed to sign message z: %v", err)
	}
	der := signature.Serialize()
	sig := append(der, byte(SigHashAll))
	sec := privateKey.Point.Serialize(true)
	scriptSig := script.Script{sig, sec}
	tx.TxIns[inputIndex].ScriptSig = &scriptSig
	txBytes, err := tx.Serialize()
	if err != nil {
		t.Fatalf("Failed to Serialize transaction message z: %v", err)
	}
	receivedHex := hex.EncodeToString(txBytes)
	if receivedHex != expectedHex {
		t.Fatalf("Something did not go right in signing the transaction.\nExpected: %s\nReceived: %s", expectedHex, receivedHex)
	}
}

func TestSignInput(t *testing.T) {
	// Create a private key with a secret value of 8675309
	privateKey, _ := signatureverification.NewPrivateKey(big.NewInt(8675309))

	// Create a transaction stream from the hex string
	txHex := "010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d00000000ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000"
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	// Create a transaction object from the stream
	tx, err := ParseTx(bufio.NewReader(bytes.NewReader(txBytes)), true)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	// Test signing input at index 0
	if !tx.SignInput(0, privateKey) {
		t.Fatal("Failed to sign input")
	}

	// Expected serialized result after signing
	wantHex := "010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d0000006b4830450221008ed46aa2cf12d6d81065bfabe903670165b538f65ee9a3385e6327d80c66d3b502203124f804410527497329ec4715e18558082d489b218677bd029e7fa306a72236012103935581e52c354cd2f484fe8ed83af7a3097005b2f9c60bff71d35bd795f54b67ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000"

	// Compare the serialized result with the expected value
	gotHex, err := tx.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialze transaction to hex: %v", err)
	}
	if hex.EncodeToString(gotHex) != wantHex {
		t.Fatalf("Unexpected serialized result. Got: %s, Want: %s", gotHex, wantHex)
	}
}

func TestValidateP2SH(t *testing.T) {
	txHex := "0100000001868278ed6ddfb6c1ed3ad5f8181eb0c7a385aa0836f01d5e4789e6bd304d87221a000000db00483045022100dc92655fe37036f47756db8102e0d7d5e28b3beb83a8fef4f5dc0559bddfb94e02205a36d4e4e6c7fcd16658c50783e00c341609977aed3ad00937bf4ee942a8993701483045022100da6bee3c93766232079a01639d07fa869598749729ae323eab8eef53577d611b02207bef15429dcadce2121ea07f233115c6f09034c0be68db99980b9a6c5e75402201475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152aeffffffff04d3b11400000000001976a914904a49878c0adfc3aa05de7afad2cc15f483a56a88ac7f400900000000001976a914418327e3f3dda4cf5b9089325a4b95abdfa0334088ac722c0c00000000001976a914ba35042cfe9fc66fd35ac2224eebdafd1028ad2788acdc4ace020000000017a91474d691da1574e6b3c192ecfb52cc8984ee7b6c568700000000"
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := ParseTx(bufio.NewReader(bytes.NewReader(txBytes)), false)
	if err != nil {
		t.Fatalf("Failed to decode parse tx: %v", err)
	}

	redeemScriptBytes, err := hex.DecodeString("475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152ae")
	if err != nil {
		t.Fatalf("Failed to decode redeemscript hex: %v", err)
	}
	redeemScript, err := script.ParseScript(bufio.NewReader(bytes.NewReader(redeemScriptBytes)))
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	tx.TxIns[0].ScriptSig = redeemScript

	modifiedTx, err := tx.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize script: %v", err)
	}

	sigHashBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sigHashBytes, SigHashAll)
	modifiedTx = append(modifiedTx, sigHashBytes...)

	h256 := utils.Hash256(modifiedTx)

	z := new(big.Int).SetBytes(h256)

	// First signature
	sec, err := hex.DecodeString("022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb70")
	if err != nil {
		t.Fatalf("Failed to decode sec hex: %v", err)
	}

	der, err := hex.DecodeString("3045022100dc92655fe37036f47756db8102e0d7d5e28b3beb83a8fef4f5dc0559bddfb94e02205a36d4e4e6c7fcd16658c50783e00c341609977aed3ad00937bf4ee942a89937")
	if err != nil {
		t.Fatalf("Failed to decode der hex: %v", err)
	}

	point, err := signatureverification.ParseSEC(sec)
	if err != nil {
		t.Fatalf("Failed to parse sec: %v", err)
	}

	sig, err := signatureverification.ParseDER(der)
	if err != nil {
		t.Fatalf("Failed to parse der: %v", err)
	}

	if !point.Verify(z, sig) {
		t.Error("failed to verify firs signature")
	}

	// Second signature
	sec, err = hex.DecodeString("03b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb71")
	if err != nil {
		t.Fatalf("Failed to decode sec hex: %v", err)
	}

	der, err = hex.DecodeString("3045022100da6bee3c93766232079a01639d07fa869598749729ae323eab8eef53577d611b02207bef15429dcadce2121ea07f233115c6f09034c0be68db99980b9a6c5e754022")
	if err != nil {
		t.Fatalf("Failed to decode der hex: %v", err)
	}

	point, err = signatureverification.ParseSEC(sec)
	if err != nil {
		t.Fatalf("Failed to parse sec: %v", err)
	}

	sig, err = signatureverification.ParseDER(der)
	if err != nil {
		t.Fatalf("Failed to parse der: %v", err)
	}

	if !point.Verify(z, sig) {
		t.Error("failed to verify second signature")
	}

}
