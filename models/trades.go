package models

import (
	"fmt"
	"time"
)

type Trade struct {
	Category          string     `csv:"Trades"`
	Header            string     `csv:"Header"`
	DataDiscriminator string     `csv:"DataDiscriminator"`
	AssetCategory     string     `csv:"Asset Category"`
	Currency          string     `csv:"Currency"`
	Symbol            string     `csv:"Symbol"`
	DateTime          *time.Time `csv:"Date/Time"`
	Quantity          float64    `csv:"Quantity"`
	TransactionPrice  *float64   `csv:"T. Price"`
	ClosePrice        *float64   `csv:"C. Price"`
	Proceeds          float64    `csv:"Proceeds"`
	CommFee           float64    `csv:"Comm/Fee"`
	Basis             float64    `csv:"Basis"`
	RealizedPL        float64    `csv:"Realized P/L"`
	MTMPL             float64    `csv:"MTM P/L"`
	Code              string     `csv:"Code"`
}

// var expectedHeaders = []string{
//     "Trades",
//     "Header",
//     "DataDiscriminator",
//     "Asset Category",
//     "Currency",
//     "Symbol",
//     "Date/Time",
//     "Quantity",
//     "T. Price",
//     "C. Price",
//     "Proceeds",
//     "Comm/Fee",
//     "Basis",
//     "Realized P/L",
//     "MTM P/L",
//     "Code",
// }

// func New(data [][]string) ([]Trade, error) {
// 	var trades []Trade
// 	for _, line := range data {
// 		// Check that the row contains 'Trades' as the first column
// 		if line[0] != "Trades" {
// 			return nil, errors.New(fmt.Sprintf("Invalid data %s", line[0]))
// 		}
// 		// This row definitely contains 'trades', try to create object
// 		trade := Trade{}
// 		for idx, _ := range line {
// 			buildTrade(&trade, idx)
// 		}
// 	}
// }

func buildTrade(t *Trade, idx int) {
	// ?
}

// Processes trades in USD and finds the PL for each trade, converted to JPY
func ProcessTrades(data [][]string) int {
	fmt.Println("> Processing Trades")
	return 0
}
