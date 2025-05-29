package exchangerateservice

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ExchangeRateService struct {
	rates map[string]map[string]float64
}

// Rates will look like
// {
// 	"USD": {
// 		"2024-10-25": 152.30,
// 		"2024-10-24": 152.25,
// 	},
// 	"CAD": {
// 		"2024-10-25": 112.30,
// 		"2024-10-24": 112.25,
// 	}
// }

func NewExchangeRateService() (*ExchangeRateService, error) {
	return &ExchangeRateService{
		rates: map[string]map[string]float64{},
	}, nil
}

func (s *ExchangeRateService) HasCurrency(currencySymbol string) bool {
	return s.rates[currencySymbol] != nil
}

func (s *ExchangeRateService) AddCurrency(currencySymbol string) {
	if s.rates[currencySymbol] == nil {
		// Performance optimization knowing that there will never be more than 365 trading days in a year
		// I could choose 255-270ish since weekends will never trade
		s.rates[currencySymbol] = make(map[string]float64, 365)
	}

	currencyPair := strings.ToLower(currencySymbol) + "jpy"
	directory := fmt.Sprintf("./data/%s/", currencyPair)
	s.LoadExchangeRates(directory, currencySymbol)
}

// NewExchangeRateService creates a new exchange rate service from a CSV file
func (s *ExchangeRateService) LoadExchangeRates(fileDirectory, currencySymbol string) error {
	fmt.Println("> Reading CSV file from path:", fileDirectory, "for currency:", currencySymbol)

	csvPath, err := os.ReadDir(fileDirectory)

	if err != nil {
		fmt.Println("Error reading directory:", err)
		return err
	}

	for i := range csvPath {
		filename := csvPath[i].Name()
		filepath := fileDirectory + filename
		s.loadCurrencyHistory(filepath, currencySymbol)
	}

	return nil
}

// loadCurrencyHistory loads exchange rates from a CSV into memory for an ExchangeRateService.
// The CSV file should have the following columns:
// Date,Open,High,Low,Close
func (s *ExchangeRateService) loadCurrencyHistory(filepath, currencySymbol string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	rates := make(map[string]float64)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Parse date (MM/DD/YYYY format)
		dateStr := record[0]
		date, err := time.Parse("01/02/2006", dateStr)
		if err != nil {
			continue
		}

		// Parse closing rate and remove quotes
		closeStr := strings.Trim(record[4], "\"")
		closeRate, err := strconv.ParseFloat(closeStr, 64)
		if err != nil {
			continue
		}

		// Store in map using YYYY-MM-DD format as key
		rates[date.Format("2006-01-02")] = closeRate
	}

	s.rates[currencySymbol] = rates
	fmt.Println("> Loaded", len(rates), "exchange rate entries for", currencySymbol)
	return nil
}

// GetRate provides the exchange rate for a specific date for a given currency
// It will backtrack up to 10 days if the exact date is not found to handle weekends and holidays.
func (s *ExchangeRateService) GetRate(originalCurrency string, date time.Time) (float64, error) {
	dateStr := date.Format("2006-01-02")
	fmt.Println("Getting exchange rate for currency:", originalCurrency, "on date:", dateStr)
	if s.rates == nil {
		fmt.Println("Rates map is nil")
		return 0.00, errors.New("Rates map is nil")
	}
	if s.rates[originalCurrency] == nil {
		fmt.Println("Rates for currency not found:", originalCurrency)
		return 0.00, errors.New("Rates for currency not found")
	}
	ratesForCurrency, ok := s.rates[originalCurrency]
	if !ok {
		return 0.00, errors.New("Currency not found")
	}

	fmt.Println(ratesForCurrency[dateStr])
	// Try exact date match
	if rateOnDate, ok := ratesForCurrency[dateStr]; ok {
		return rateOnDate, nil
	}

	// Exact date not found, backtrack to closest previous date
	// This handles weekends and holidays
	for range 10 { // Try up to 10 days back
		date = date.AddDate(0, 0, -1)
		dateStr = date.Format("2006-01-02")
		if rate, ok := ratesForCurrency[dateStr]; ok {
			return rate, nil
		}
	}

	// Crash if not found
	panic("No exchange rate found for the given date, or up to 10 days back")
}

// func main() {
// 	// Example usage
// 	rates, err := loadExchangeRates("../data/usdjpy/2024.csv")
// 	if err != nil {
// 		fmt.Println("Error loading exchange rates:", err)
// 		return
// 	}

// 	// Example lookup
// 	tradeDate := time.Date(2024, 10, 25, 0, 0, 0, 0, time.UTC)
// 	rate := getExchangeRate(rates, tradeDate)
// 	fmt.Printf("USD/JPY rate for %s: %.2f\n", tradeDate.Format("2006-01-02"), rate)
// }
