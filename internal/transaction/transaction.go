package transaction

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/caspereijkens/cryptocurrency/internal/script"
	"github.com/caspereijkens/cryptocurrency/internal/utils"
)

type Tx struct {
	Version  uint32
	TxIns    []*TxIn
	TxOuts   []*TxOut
	Locktime uint32
	Testnet  bool
}

func NewTx(version uint32, txIns []*TxIn, txOuts []*TxOut, locktime uint32, testnet bool) *Tx {
	return &Tx{
		Version:  version,
		TxIns:    txIns,
		TxOuts:   txOuts,
		Locktime: locktime,
		Testnet:  testnet,
	}
}

func (tx *Tx) String() string {
	txInsStr := ""
	for _, txIn := range tx.TxIns {
		txInsStr += fmt.Sprintf("%s\n", txIn.String())
	}
	txOutsStr := ""
	for _, txOut := range tx.TxOuts {
		txOutsStr += fmt.Sprintf("%s\n", txOut.String())
	}
	id, err := tx.Id()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("tx: %s\nversion: %d\ntx_ins:\n%s\n"+
		"tx_outs:\n%s\nlocktime: %d", id, tx.Version, txInsStr, txOutsStr, tx.Locktime)
}

func (tx *Tx) Id() (string, error) {
	hash256, err := tx.Hash()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash256), nil
}

func (tx *Tx) Hash() ([]byte, error) {
	s, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	hash256 := utils.Hash256(s)
	slices.Reverse(hash256)
	return hash256, nil
}

func ParseTx(reader *bufio.Reader, testnet bool) (*Tx, error) {
	// version is an integer in 4 bytes, little-endian
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil {
		return nil, err
	}

	numInputs, err := utils.ReadVarint(reader)
	if err != nil {
		return nil, err
	}

	inputs := make([]*TxIn, 0, numInputs)
	for i := 0; i < int(numInputs); i++ {
		txIn, err := ParseTxIn(reader)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, txIn)
	}

	numOutputs, err := utils.ReadVarint(reader)
	if err != nil {
		return nil, err
	}

	// parse num_outputs number of TransactionOutputs
	outputs := make([]*TxOut, 0, numOutputs)
	for i := 0; i < int(numOutputs); i++ {
		txOut, err := ParseTxOut(reader)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, txOut)
	}

	// locktime is an integer in 4 bytes, little-endian
	var locktime uint32
	if err := binary.Read(reader, binary.LittleEndian, &locktime); err != nil {
		return nil, err
	}

	return NewTx(version, inputs, outputs, locktime, testnet), nil
}

func (tx *Tx) Serialize() ([]byte, error) {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, tx.Version)

	numInputs, err := utils.EncodeVarint(uint64(len(tx.TxIns)))
	if err != nil {
		return nil, err
	}

	result = append(result, numInputs...)

	for _, txIn := range tx.TxIns {
		serializedTxIn, err := txIn.Serialize()
		if err != nil {
			return nil, err
		}
		result = append(result, serializedTxIn...)
	}

	numOutputs, err := utils.EncodeVarint(uint64(len(tx.TxOuts)))
	if err != nil {
		return nil, err
	}

	result = append(result, numOutputs...)

	for _, txOut := range tx.TxOuts {
		serializedTxOut, err := txOut.Serialize()
		if err != nil {
			return nil, err
		}
		result = append(result, serializedTxOut...)
	}

	locktimeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktimeBytes, tx.Locktime)
	result = append(result, locktimeBytes...)

	return result, nil
}

func (tx *Tx) Fee() (uint64, error) {
	// initialize input sum and output sum
	var inputSum, outputSum uint64

	// use TransactionInput.Value() to sum up the input amounts
	for _, txIn := range tx.TxIns {
		value, err := txIn.Value(tx.Testnet)
		if err != nil {
			return 0, err
		}
		inputSum += value
	}

	// use TransactionOutput.Amount to sum up the output amounts
	for _, txOut := range tx.TxOuts {
		outputSum += txOut.Amount
	}

	fee := inputSum - outputSum
	return fee, nil
}

// TxIn represents a transaction input
type TxIn struct {
	PrevTx    []byte
	PrevIndex uint32
	ScriptSig script.Script
	Sequence  uint32
}

// NewTxIn creates a new TxIn instance
func NewTxIn(prevTx []byte, prevIndex uint32, scriptSig script.Script, sequence uint32) *TxIn {
	return &TxIn{
		PrevTx:    prevTx,
		PrevIndex: prevIndex,
		ScriptSig: scriptSig,
		Sequence:  sequence,
	}
}

// String returns a string representation of TxIn
func (txIn *TxIn) String() string {
	return fmt.Sprintf("%s:%d", hex.EncodeToString(txIn.PrevTx), txIn.PrevIndex)
}

// ParseTxIn parses a byte stream and returns a TxIn object
// Possible IP: seems like because of historical reasons, the prevTxId was reversed: https://learnmeabitcoin.com/technical/txid
func ParseTxIn(reader *bufio.Reader) (*TxIn, error) {
	// prev_tx is 32 bytes, little endian
	prevTX := make([]byte, 32)
	if _, err := io.ReadFull(reader, prevTX); err != nil {
		return nil, err
	}
	slices.Reverse(prevTX)
	// prevIndex is an integer in 4 bytes, little endian
	var prevIndex uint32
	if err := binary.Read(reader, binary.LittleEndian, &prevIndex); err != nil {
		return nil, err
	}
	// use script.ParseScript to get the ScriptSig
	scriptSig, err := script.ParseScript(reader)
	if err != nil {
		return nil, err
	}
	// sequence is an integer in 4 bytes, little-endian
	var sequence uint32
	if err := binary.Read(reader, binary.LittleEndian, &sequence); err != nil {
		return nil, err
	}
	// return an instance of the class
	return NewTxIn(prevTX, prevIndex, scriptSig, sequence), nil
}

// Serialize returns the byte serialization of the transaction input
func (txIn *TxIn) Serialize() ([]byte, error) {
	var result []byte

	// serialize prev_tx, little endian
	prevTxLittleEndian := make([]byte, 32)
	copy(prevTxLittleEndian, txIn.PrevTx)
	slices.Reverse(prevTxLittleEndian)
	result = append(result, prevTxLittleEndian...)

	// serialize prev_index, 4 bytes, little endian
	prevIndexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(prevIndexBytes, txIn.PrevIndex)
	result = append(result, prevIndexBytes...)

	// serialize the ScriptSig
	scriptSig, err := txIn.ScriptSig.Serialize()
	if err != nil {
		return nil, err
	}
	result = append(result, scriptSig...)

	// serialize sequence, 4 bytes, little endian
	sequenceBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sequenceBytes, txIn.Sequence)
	result = append(result, sequenceBytes...)

	return result, nil
}

func (txIn *TxIn) FetchTx(testnet bool) (*Tx, error) {
	return NewTxFetcher().Fetch(hex.EncodeToString(txIn.PrevTx), testnet, false)
}

func (txIn *TxIn) Value(testnet bool) (uint64, error) {
	tx, err := txIn.FetchTx(testnet)
	if err != nil {
		return 0, err
	}

	numOutputs := uint32(len(tx.TxOuts))
	if txIn.PrevIndex >= numOutputs {
		return 0, fmt.Errorf("previous index %d out of range for transaction outputs", txIn.PrevIndex)
	}

	return tx.TxOuts[txIn.PrevIndex].Amount, nil
}

// TransactionInput represents a transaction input
type TxOut struct {
	Amount       uint64
	ScriptPubkey script.Script
}

// NewTransactionInput creates a new TxIn instance
func NewTxOut(amount uint64, scriptPubkey script.Script) *TxOut {
	return &TxOut{
		Amount:       amount,
		ScriptPubkey: scriptPubkey,
	}
}

// String returns a string representation of TxIn
func (txOut *TxOut) String() string {
	return fmt.Sprintf("%s:%s", utils.FormatWithUnderscore(int(txOut.Amount)), txOut.ScriptPubkey.String())
}

// ParseTxOut parses a byte stream and returns a TxOut object
func ParseTxOut(reader *bufio.Reader) (*TxOut, error) {
	var amount uint64
	if err := binary.Read(reader, binary.LittleEndian, &amount); err != nil {
		return nil, err
	}

	scriptPubkey, err := script.ParseScript(reader)
	if err != nil {
		return nil, err
	}

	return NewTxOut(amount, scriptPubkey), nil
}

// Serialize returns the byte serialization of the transaction output
func (txOut *TxOut) Serialize() ([]byte, error) {
	amountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(amountBytes, txOut.Amount)

	scriptPubkeyBytes, err := txOut.ScriptPubkey.Serialize()
	if err != nil {
		return nil, err
	}

	result := append(amountBytes, scriptPubkeyBytes...)

	return result, nil
}

type TxFetcher struct {
	Cache map[string]*Tx
}

func NewTxFetcher() *TxFetcher {
	return &TxFetcher{
		Cache: make(map[string]*Tx),
	}
}

func (tf *TxFetcher) GetURL(testnet bool) string {
	if testnet {
		return "https://blockstream.info/testnet/api"
	}
	return "https://blockstream.info/api"
}

func (tf *TxFetcher) Fetch(txID string, testnet, fresh bool) (*Tx, error) {
	if !fresh {
		if cachedTx, ok := tf.Cache[txID]; ok {
			cachedTx.Testnet = testnet
			return cachedTx, nil
		}
	}

	url := fmt.Sprintf("%s/tx/%s/hex", tf.GetURL(testnet), txID)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rawHex, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	raw, err := hex.DecodeString(string(rawHex))
	if err != nil {
		return nil, err
	}

	var tx *Tx
	if raw[4] == 0 {
		raw = append(raw[:4], raw[6:]...)
		tx, err = ParseTx(bufio.NewReader(bytes.NewBuffer(raw)), testnet)
		if err != nil {
			return nil, err
		}
		tx.Locktime = binary.LittleEndian.Uint32(raw[len(raw)-4:])
	} else {
		tx, err = ParseTx(bufio.NewReader(bytes.NewBuffer(raw)), testnet)
		if err != nil {
			return nil, err
		}
	}

	id, err := tx.Id()
	if err != nil {
		return nil, err
	}

	if id != txID {
		return nil, fmt.Errorf("not the same id: %s vs %s", id, txID)
	}

	tf.Cache[txID] = tx
	return tx, nil
}

func (tf *TxFetcher) LoadCache(filename string) error {
	diskCacheFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer diskCacheFile.Close()

	diskCache := make(map[string]string)
	err = json.NewDecoder(diskCacheFile).Decode(&diskCache)
	if err != nil {
		return err
	}

	for k, rawHex := range diskCache {
		raw, err := hex.DecodeString(rawHex)
		if err != nil {
			return err
		}

		var tx *Tx
		if raw[4] == 0 {
			raw = append(raw[:4], raw[6:]...)
			tx, err = ParseTx(bufio.NewReader(bytes.NewReader(raw)), false)
			if err != nil {
				return err
			}
			// TODO Why is this reassigning the Locktime?
			// tx.Locktime = binary.LittleEndian.Uint32(raw[len(raw)-4:])
		} else {
			tx, err = ParseTx(bufio.NewReader(bytes.NewReader(raw)), false)
			if err != nil {
				return err
			}
		}

		tf.Cache[k] = tx
	}

	return nil
}

func (tf *TxFetcher) DumpCache(filename string) error {
	diskCacheFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer diskCacheFile.Close()

	toDump := make(map[string]string)
	for k, tx := range tf.Cache {
		serializedTx, err := tx.Serialize()
		if err != nil {
			return err
		}
		toDump[k] = hex.EncodeToString(serializedTx)
	}

	err = json.NewEncoder(diskCacheFile).Encode(toDump)
	if err != nil {
		return err
	}

	return nil
}
