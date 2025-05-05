package main

// This script reads your Interactive Brokers summary of trades, dividends, interest, and withholding tax for the year,
// finds the cost basis for each trade, and calculates the profit/loss for each trade
// Run with `go run . <filepath_to_trades_csv.csv>`

import (
	"fmt"
	"os"

	"github.com/ibkr-cost-basis-calc/models"
	"github.com/ibkr-cost-basis-calc/utils"
)

// Make the CLI output a little prettier
const colorYellow = "\033[0;33m"
const colorReset = "\033[0m"
const colorGreen = "\033[0;32m"

// Replace this with your actual file name
// The file should be in the same directory as this script
// or provide the relative path to the file
var FILENAME = "2024_trades.csv"

func main() {
	// Check if a filename was provided as a command line argument
	if len(os.Args) > 1 {
		FILENAME = os.Args[1]
	}

	file, err := os.Open(FILENAME)
	if err != nil {
		fmt.Println("Error: unable to open file")
		panic(err)
	}

	defer file.Close()

	// Reads a csv and returns a map[firstColumn][][]row
	m := utils.ReadCSV(file)

	for key := range m {
		fmt.Fprintf(os.Stdout, "> Processing key %s %s %s\n", colorYellow, key, colorReset)
		switch key {
		case "Trades":
			models.ProcessTrades(m[key])
		case "Dividends":
			models.ProcessDividends(m[key])
		case "Withholding Tax":
			models.ProcessWithholdingTax(m[key])
		case "Interest":
			models.ProcessInterest(m[key])
		}
	}

}
