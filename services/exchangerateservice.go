package exchangerateservice

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type ExchangeRateService struct {
	rates map[string]float64
}

// NewExchangeRateService creates a new exchange rate service from a CSV file
func NewExchangeRateService(filepath string) (*ExchangeRateService, error) {
	rates, err := loadExchangeRates(filepath)
	if err != nil {
		return nil, err
	}

	return &ExchangeRateService{
		rates: rates,
	}, nil
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

// Get exchange rate for a specific date
func (s *ExchangeRateService) GetRate(date time.Time) (float64, error) {
	dateStr := date.Format("2006-01-02")

	// Try exact date match
	if rate, ok := s.rates[dateStr]; ok {
		return rate, nil
	}

	// Exact date not found, backtrack to closest previous date
	// This handles weekends and holidays
	currentDate := date
	for range 10 { // Try up to 10 days back
		currentDate = currentDate.AddDate(0, 0, -1)
		dateStr = currentDate.Format("2006-01-02")
		if rate, ok := s.rates[dateStr]; ok {
			return rate, nil
		}
	}

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
