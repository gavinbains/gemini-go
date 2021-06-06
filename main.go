package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

const sandboxUrl = "https://api.sandbox.gemini.com"
const geminiUrl = "https://api.gemini.com"
const baseUrl = sandboxUrl

const buySize = 300.00
const tickerSymbol = "BTCUSD"
const precision = "%.8f" // BTC is 8 decimals, ETH is 6, etc

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	allTickers := getAllTickers(baseUrl)
	log.Println(allTickers)

	ticker := getTicker(baseUrl, tickerSymbol)
	log.Println(ticker)

	priceAsString := fmt.Sprintf("%.2f", ticker.Ask*.999)
	price, err := strconv.ParseFloat(priceAsString, 64)
	if err != nil {
		log.Fatalln(err)
	}
	// must round to precision of currency and *.999 for fee inclusion
	amount := fmt.Sprintf(precision, (buySize*.999)/price)
	log.Println(amount)

	response := placeOrder(tickerSymbol, amount, priceAsString, baseUrl)
	log.Println(response)
}
