package models

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/ibkr-cost-basis-calc/utils"
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
	Proceeds          float64    `csv:"Proceeds"`
	CommFee           float64    `csv:"Comm/Fee"`
	Basis             float64    `csv:"Basis"`
	RealizedPL        float64    `csv:"Realized P/L"`
	Code              string     `csv:"Code"`
}

func NewTrade(data []string) (*Trade, error) {
	trade := &Trade{}

	if data[0] != "Trades" {
		return nil, fmt.Errorf("Invalid data. First column should be 'Trades', got %s", data[0])
	}

	tradeType := reflect.TypeOf(Trade{})
	numFields := tradeType.NumField()
	if len(data) != numFields {
		return nil, fmt.Errorf("Trades needs %v properties. Got %v with values %v", numFields, len(data[0]), data)
	}

	// Map each element to its corresponding field
	trade.Category = data[0]          // "Trades"
	trade.Header = data[1]            // "Header"
	trade.DataDiscriminator = data[2] // "DataDiscriminator"
	trade.AssetCategory = data[3]     // "Asset Category"
	trade.Currency = data[4]          // "Currency"
	trade.Symbol = data[5]            // "Symbol"

	// Handle DateTime (pointer to time.Time)
	if data[6] != "" {
		dateTime, err := utils.ParseDate(data[6])
		if err != nil {
			return nil, fmt.Errorf("failed to parse DateTime %s: %v", data[6], err)
		}

		trade.DateTime = &dateTime
	}

	// Handle Quantity
	quantity, err := strconv.ParseFloat(data[7], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Quantity %s: %v", data[7], err)
	}
	trade.Quantity = quantity

	// Handle TransactionPrice (pointer to float64)
	if data[8] != "" {
		tPrice, err := strconv.ParseFloat(data[8], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TransactionPrice %s: %v", data[8], err)
		}
		trade.TransactionPrice = &tPrice
	}

	// Handle Proceeds
	proceeds, err := strconv.ParseFloat(data[9], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Proceeds %s: %v", data[9], err)
	}
	trade.Proceeds = proceeds

	// Handle CommFee
	commFee, err := strconv.ParseFloat(data[10], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CommFee %s: %v", data[10], err)
	}
	trade.CommFee = commFee

	// Handle Basis
	basis, err := strconv.ParseFloat(data[11], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Basis %s: %v", data[11], err)
	}
	trade.Basis = basis

	// Handle RealizedPL
	realizedPL, err := strconv.ParseFloat(data[12], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RealizedPL %s: %v", data[12], err)
	}
	trade.RealizedPL = realizedPL

	// Handle Code
	trade.Code = data[13]

	return trade, nil
}

// Processes trades in USD and finds the PL for each trade, converted to JPY
func ProcessTrades(data [][]string) int {
	fmt.Println("> Processing Trades")
	trades := make([]*Trade, 0)

	for _, line := range data {
		// Only process trades that are orders
		// It's possible that it's a Subtotal row or a Header row, in which case skip
		if line[1] == "Header" || line[1] == "SubTotal" {
			continue
		}
		if line[1] == "Data" && line[2] == "Order" {
			trade, err := NewTrade(line)
			if err != nil {
				fmt.Printf("Error creating trade: %v\n", err)
				return 0
			}
			trades = append(trades, trade)
		} else {
			fmt.Printf("Not a trade: %v\n", line)
			panic("Trade is not of the right format")
		}

	}
	fmt.Printf("Number of trades processed: %d\n", len(trades))

	fmt.Println("Trade details:")
	for _, trade := range trades {
		fmt.Printf("Trade: %+v %v with basis %v\n", trade.Symbol, trade.Quantity, trade.Basis)
	}
	return 0
}
