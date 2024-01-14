package main

import (
	"flag"
	"fmt"

	"github.com/caspereijkens/cryptocurrency/internal/transaction"
)

func main() {
	// Define a boolean flag
	var isTestnet bool
	var fresh = true
	flag.BoolVar(&isTestnet, "testnet", false, "enable testnet mode")

	// Parse the command-line arguments
	flag.Parse()

	// Retrieve the non-flag command-line arguments
	args := flag.Args()

	// Check if at least one argument is provided
	if len(args) == 0 {
		fmt.Println("Please provide a transaction ID.")
		return
	}

	// Extract the transaction ID
	transactionID := args[0]

	tx, err := transaction.NewTxFetcher().Fetch(transactionID, isTestnet, fresh)
	if err != nil {
		fmt.Println("Transaction could not be found. Please provide a correct transaction ID.")
		return
	}

	fmt.Println(tx.String())
}
