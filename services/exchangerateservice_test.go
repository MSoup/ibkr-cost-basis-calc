package exchangerateservice

import (
	"fmt"
	"os"
	"testing"
	"time"
)

const rateDirectory = "../data/usdjpy/2024.csv"

var rateService *ExchangeRateService

// Setup shared test resources
func TestMain(m *testing.M) {
	var err error
	rateService, err = NewExchangeRateService()
	if err != nil {
		fmt.Println("Error creating ExchangeRateService:", err)
		os.Exit(1)
	}

	rateService.AddCurrency("USD")

	// Run the tests
	code := m.Run()

	os.Exit(code)
}

// TestGetRate calls exchangerateservice.GetRate with a date
// for a valid return value.
func TestGetRate(t *testing.T) {
	tradeDate := time.Date(2024, 10, 25, 0, 0, 0, 0, time.UTC) // 10/25/2024
	want := 152.30                                             // Rate on 10/25/2024
	rate, err := rateService.GetRate("USD", tradeDate)
	if want != rate || err != nil {
		t.Errorf(`GetRate(%v) = %v, %v, want match for %v, nil`, want, rate, err, want)
	}
}

// TestGetRateEmpty calls GetRate with a non-existent (in csv) date
// checking for an error.
func TestGetRateInvalidTime(t *testing.T) {
	currentDate := time.Time{}
	rate, err := rateService.GetRate("USD", currentDate)
	if err == nil {
		t.Errorf(`GetRate(%v) = %v, %v, want match for 0.00, %v`, currentDate, rate, err, nil)
	}
}
