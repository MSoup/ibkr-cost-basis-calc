package main

import (
	"fmt"
	"os"

	"github.com/ibkr-cost-basis-calc/models"
	"github.com/ibkr-cost-basis-calc/utils"
)

func main() {
	const FILENAME = "2024_trades.csv"
	const colorYellow = "\033[0;33m"
	const colorReset = "\033[0m"

	file, err := os.Open(FILENAME)
	if err != nil {
		fmt.Println("Error: unable to open file")
		panic(err)
	}

	defer file.Close()

	// Reads a csv and returns a map[firstColumn][][]row
	m := utils.ReadCSV(file)

	for key := range m {
		fmt.Fprintf(os.Stdout, "> Looking at key: %s %s %s\n", colorYellow, key, colorReset)
		fmt.Println(m[key][0])
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
