package models

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	exchangerateservice "github.com/ibkr-cost-basis-calc/services"
	"github.com/ibkr-cost-basis-calc/utils"
)

type Trade struct {
	Category          string    `csv:"Trades"`
	Header            string    `csv:"Header"`
	DataDiscriminator string    `csv:"DataDiscriminator"`
	AssetCategory     string    `csv:"Asset Category"`
	Currency          string    `csv:"Currency"`
	Symbol            string    `csv:"Symbol"`
	DateTime          time.Time `csv:"Date/Time"`
	Quantity          int64     `csv:"Quantity"`
	TransactionPrice  float64   `csv:"T. Price"`
	Proceeds          float64   `csv:"Proceeds"`
	CommFee           float64   `csv:"Comm/Fee"`
	Basis             float64   `csv:"Basis"`
	RealizedPL        float64   `csv:"Realized P/L"`
	Code              string    `csv:"Code"`
	Extra             CalculatedData
}

type CalculatedData struct {
	Action      string  // "BUY" or "SELL"
	JPYRate     float64 // Exchange rate for the transaction
	ProceedsJPY float64 // Proceeds converted to JPY
	BasisJPY    float64 // Basis converted to JPY
	PLJPY       float64 // Profit/Loss in JPY
}

var rateService, _ = exchangerateservice.NewExchangeRateService()

const (
	CSVCategoryIndex          = 0
	CSVHeaderIndex            = 1
	CSVDataDiscriminatorIndex = 2
	CSVAssetCategoryIndex     = 3
	CSVCurrencyIndex          = 4
	CSVSymbolIndex            = 5
	CSVDateTimeIndex          = 6 // Transaction date string; parse as needed.
	CSVQuantityIndex          = 7
	CSVTransactionPriceIndex  = 8
	CSVProceedsIndex          = 9
	CSVCommFeeIndex           = 10
	CSVBasisIndex             = 11
	CSVRealizedPLIndex        = 12
	CSVCodeIndex              = 13
)

func NewTrade(data []string) (*Trade, error) {
	// Sanity check
	if data[CSVCategoryIndex] != "Trades" {
		return nil, fmt.Errorf("Invalid data. Expected 'Trades', got %s", data[CSVCategoryIndex])
	}
	// Find how many fields Trade should have and validate against data
	tradeType := reflect.TypeOf(Trade{})
	numFields := tradeType.NumField()
	if len(data) != numFields-1 {
		return nil, fmt.Errorf("Trades needs %v properties. Got %v with values %v", numFields, len(data[CSVCategoryIndex]), data)
	}

	transactionDate := parseDate(data[CSVDateTimeIndex])

	trade := &Trade{
		Category:          data[CSVCategoryIndex],
		Header:            data[CSVHeaderIndex],
		DataDiscriminator: data[CSVDataDiscriminatorIndex],
		AssetCategory:     data[CSVAssetCategoryIndex],
		Currency:          data[CSVCurrencyIndex],
		Symbol:            data[CSVSymbolIndex],
		DateTime:          transactionDate,
		Quantity:          parseInt(data[CSVQuantityIndex]),
		TransactionPrice:  parseFloat(data[CSVTransactionPriceIndex]),
		Proceeds:          parseFloat(data[CSVProceedsIndex]),
		CommFee:           parseFloat(data[CSVCommFeeIndex]),
		Basis:             parseFloat(data[CSVBasisIndex]),
		RealizedPL:        parseFloat(data[CSVRealizedPLIndex]),
		Code:              data[CSVCodeIndex],
	}

	// If currency isn't loaded, load it
	if !rateService.HasCurrency(trade.Currency) {
		rateService.AddCurrency(trade.Currency)
	}

	// rate, err := rateService.GetRate(trade.Currency, transactionDate)
	// if err != nil {
	// 	log.Fatalf("Error getting exchange rate: %v", err)
	// }

	// extra := CalculatedData{
	// 	Action:      parseAction(data[9]),
	// 	JPYRate:     rate,
	// 	ProceedsJPY: getJPYProceeds(parseFloat(data[9]), rate),
	// }

	// trade.Extra = extra

	return trade, nil
}

func parseAction(proceeds string) string {
	// Negative proceeds means I bought something, see:
	// 	Trades,Data,Order,Stocks,USD,SOFI,"2024-01-12, 16:20:00",200,9.5,-1900,0,1900,0,A;O
	//  Trades,Data,Order,Stocks,USD,SOFI,"2024-11-15, 16:20:00",-100,12.5,1250,-0.05135,-922.681914,345.211868,A;C
	//  Trades,SubTotal,,Stocks,USD,SOFI,,100,,-650,-0.05135,977.318086,345.211868,
	if strings.HasPrefix(proceeds, "-") {
		return "BUY"
	}
	return "SELL"
}

func getJPYProceeds(proceeds float64, fxRate float64) float64 {
	return proceeds * fxRate
}

// parseDate parses a string into a time.Time object
// But also truncates it to fit YYYY-MM-DD
func parseDate(dateStr string) time.Time {
	dateTime, err := utils.ParseDate(dateStr[:10])
	if len(dateStr) < 10 {
		fmt.Printf("Error parsing date--need YYYY-MM-DD at the very least. Got: %v\n", dateStr)
	}
	if err != nil {
		fmt.Printf("Error parsing date: %v\n", err)
		panic(err)
	}

	return dateTime
}

func parseInt(value string) int64 {
	if value == "" {
		panic("Empty value for int")
	}
	intValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		fmt.Printf("Error parsing int: %v\n", err)
		panic(err)
	}
	return intValue
}

func parseFloat(value string) float64 {
	if value == "" {
		panic("Empty value for float")
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Printf("Error parsing float: %v\n", err)
		panic(err)
	}
	return floatValue
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
		fmt.Printf("Trade: %+v %5d with basis %10.2f\n", trade.Symbol, trade.Quantity, trade.Basis)
	}
	return 0
}
