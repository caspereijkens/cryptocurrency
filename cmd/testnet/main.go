package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/caspereijkens/cryptocurrency/internal/signatureverification"
	"github.com/caspereijkens/cryptocurrency/internal/util"
)

func main() {
	var data string

	// Create a new scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Type a long secret that only you know: ")

	// Use scanner to read the entire line, including spaces
	if scanner.Scan() {
		data = scanner.Text()
	}
	fmt.Print("\n")
	hash256 := util.Hash256([]byte(data))

	// Convert the second hash bytes to a big.Int
	bigInt := new(big.Int)
	bigInt.SetBytes(hash256)

	privKey, err := signatureverification.NewPrivateKey(bigInt)
	if err != nil {
		panic(err)
	}

	address := privKey.Point.Address(true, true)

	fmt.Println("The testnet address that is connected to this secret is:")
	fmt.Println(address)

	fmt.Print("\n")
	fmt.Println("now go to https://coinfaucet.eu/en/btc-testnet/ and enter this address. Press 'Get bitcoins!'")
}
