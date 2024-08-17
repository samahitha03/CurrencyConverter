package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Define a struct to match the API response
type ExchangeRates struct {
	ConversionRates map[string]float64 `json:"conversion_rates"`
	Base            string             `json:"base_code"`
}

func main() {
	apiKey := "8758b5e30c6b18bad6fc2130"

	var baseCurrency, targetCurrency string
	var amount float64

	// Get user input
	fmt.Print("Enter base currency (e.g., USD): ")
	fmt.Scanln(&baseCurrency)
	baseCurrency = strings.ToUpper(baseCurrency)
	fmt.Print("Enter target currency (e.g., EUR): ")
	fmt.Scanln(&targetCurrency)
	targetCurrency = strings.ToUpper(targetCurrency)

	fmt.Print("Enter amount: ")
	fmt.Scanln(&amount)

	// Fetch exchange rates and perform conversion
	rates, err := fetchExchangeRates(apiKey, baseCurrency)
	if err != nil {
		fmt.Println("Error fetching exchange rates:", err)
		return
	}

	convertedAmount := convertCurrency(rates, targetCurrency, amount)
	fmt.Printf("%.2f %s = %.2f %s\n", amount, baseCurrency, convertedAmount, targetCurrency)
}

func fetchExchangeRates(apiKey, baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", apiKey, baseCurrency)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rates ExchangeRates
	if err := json.Unmarshal(body, &rates); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	//fmt.Println("Received Rates:", rates.ConversionRates)

	return rates.ConversionRates, nil
}

func convertCurrency(rates map[string]float64, targetCurrency string, amount float64) float64 {
	rate, exists := rates[targetCurrency]
	if !exists {
		fmt.Printf("Target currency %s not found in the rates.\n", targetCurrency)
		os.Exit(1)
	}
	return amount * rate
}
