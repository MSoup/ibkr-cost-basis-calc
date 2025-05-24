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

func (s *ExchangeRateService) HasCurrency(currency string) bool {
	return s.rates[currency] != nil
}

func (s *ExchangeRateService) AddCurrency(currency string) {
	if s.rates[currency] == nil {
		s.rates[currency] = map[string]float64{}
	}

	s.LoadExchangeRates("./data/usdjpy/", currency)
}

// NewExchangeRateService creates a new exchange rate service from a CSV file
func (s *ExchangeRateService) LoadExchangeRates(fileDirectory, currency string) (int, error) {
	fmt.Println("Loading exchange for currency:", currency)
	fmt.Println("Reading CSV file from path:", fileDirectory)

	csvPath, err := os.ReadDir(fileDirectory)

	if err != nil {
		fmt.Println("Error reading directory:", err)
		return 1, err
	}
	fmt.Println("CSV files in directory:", csvPath)

	for i := range csvPath {
		filename := csvPath[i].Name()
		fmt.Println(filename)

	}
	// rates, err := loadExchangeRates(filepath)
	// if err != nil {
	// 	return nil, err
	// }

	return 0, nil
}

// Loads exchange rates from CSV
func loadExchangeRates(filepath string) (map[string]float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
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

	return rates, nil
}

// GetRate provides the exchange rate for a specific date for a given currency
// It will backtrack up to 10 days if the exact date is not found
func (s *ExchangeRateService) GetRate(originalCurrency string, date time.Time) (float64, error) {
	dateStr := date.Format("2006-01-02")

	// Check if the currency exists in the rates map
	fmt.Println(1)

	if s.rates == nil {
		fmt.Println("Rates map is nil")
		return 0.00, errors.New("Rates map is nil")
	}
	ratesForCurrency, ok := s.rates[originalCurrency]
	if !ok {
		return 0.00, errors.New("Currency not found")
	}

	fmt.Println(1)
	fmt.Println(ratesForCurrency)
	fmt.Println(2)
	fmt.Println(ratesForCurrency[dateStr])
	// // Try exact date match
	// if rateOnDate, ok := ratesForCurrency[dateStr]; ok {
	// 	return rate, nil
	// }

	// // Exact date not found, backtrack to closest previous date
	// // This handles weekends and holidays
	// currentDate := date
	// for range 10 { // Try up to 10 days back
	// 	currentDate = currentDate.AddDate(0, 0, -1)
	// 	dateStr = currentDate.Format("2006-01-02")
	// 	if rate, ok := s.rates[dateStr]; ok {
	// 		return rate, nil
	// 	}
	// }

	// Default fallback
	return 0.00, errors.New("No exchange rate found for the given date, or up to 10 days back")
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
