package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caspereijkens/cryptocurrency/internal/script"
	"github.com/caspereijkens/cryptocurrency/internal/signatureverification"
	"github.com/caspereijkens/cryptocurrency/internal/transaction"
	"github.com/caspereijkens/cryptocurrency/internal/utils"
)

func main() {
	// Define command-line flags
	var inFlags, outFlags []string
	var secret string

	// Parse command-line arguments
	flag.Var((*stringSlice)(&inFlags), "in", "Input file(s)")
	flag.Var((*stringSlice)(&outFlags), "out", "Output file(s)")

	// Parse the command-line
	flag.Parse()

	txIns := parseTxIns(inFlags)
	txOuts := parseTxOuts(outFlags)

	tx := transaction.NewTx(uint32(1), txIns, txOuts, uint32(0), true)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("\nType the secret you want to sign this transaction with: ")

	if scanner.Scan() {
		secret = scanner.Text()
	}

	privateKey, err := signatureverification.NewPrivateKey(utils.Hash256ToBigInt(secret))
	if err != nil {
		panic("couldn't create private key with this")
	}

	tx.SignInput(uint32(0), privateKey)

	fmt.Println("The following transaction was SIGNED:")
	fmt.Println(tx.String())

	txBytes, err := tx.Serialize()
	if err != nil {
		panic("couldn't serialize this transaction")
	}

	fmt.Printf("The transaction is:\n\n%s\n\n", hex.EncodeToString(txBytes))

	fmt.Println("You can broadcast the transaction at https://blockstream.info/testnet/tx/push")
}

// Custom type to handle multiple string values for a flag
type stringSlice []string

func (ss *stringSlice) String() string {
	return fmt.Sprint(*ss)
}

func (ss *stringSlice) Set(value string) error {
	*ss = append(*ss, value)
	return nil
}

// Function to parse -in flags and create TxIn instances
func parseTxIns(ins []string) []*transaction.TxIn {
	var txIns []*transaction.TxIn

	for _, in := range ins {
		parts := strings.Split(in, ":")
		if len(parts) != 2 {
			fmt.Println("Invalid -in argument:", in)
			continue
		}

		txID, err := hex.DecodeString(parts[0])
		if err != nil {
			fmt.Println("Invalid hex encoding in in -in argument:", in)
			continue
		}
		index, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			fmt.Println("Invalid index in -in argument:", in)
			continue
		}

		txIn := transaction.NewTxIn(txID, uint32(index), script.Script{}, uint32(0xffffffff))
		txIns = append(txIns, txIn)
	}

	return txIns
}

// Function to parse -out flags and create TxOut instances
func parseTxOuts(outs []string) []*transaction.TxOut {
	var txOuts []*transaction.TxOut

	for _, out := range outs {
		parts := strings.Split(out, ":")
		if len(parts) != 2 {
			fmt.Printf("Invalid -out argument: %s\nUsage:-out <amount>:<address>\n", out)
			continue
		}

		amount, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			fmt.Println("Invalid amount in -out argument:", out)
			continue
		}

		addressH160, _ := utils.DecodeBase58(parts[1])
		scriptPubkey := script.CreateP2pkhScript(addressH160)

		txOut := transaction.NewTxOut(amount, scriptPubkey)
		txOuts = append(txOuts, txOut)
	}

	return txOuts
}
