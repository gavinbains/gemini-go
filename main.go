package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const baseUrl = "https://api.gemini.com"
const sandboxUrl = "https://api.sandbox.gemini.com"

const buySize = 400.00
const tickerSymbol = "BTCUSD"
const precision = "%.8f"

type Ticker struct {
	Bid    float64     `json:"bid,string"`
	Ask    float64     `json:"ask,string"`
	Volume interface{} `json:"volume"`
	Last   float64     `json:"last,string"`
}

type Order struct {
	Request string   `json:"request"`
	Nonce   string   `json:"nonce"`
	Symbol  string   `json:"symbol"`
	Amount  string   `json:"amount"`
	Price   string   `json:"price"`
	Side    string   `json:"side"`
	Type    string   `json:"type"`
	Options []string `json:"options"`
}

func roundTo2(x float64) float64 {
	return math.Round(x*100) / 100
}

func makeTimestamp() string {
	return fmt.Sprint(int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond))
}

func Base64Encode(message []byte) []byte {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(b, message)
	return b
}

func computeHmacInHex(message []byte, secret []byte) string {
	var h hash.Hash

	h = hmac.New(sha512.New384, secret)

	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

func getAllTickers(url string) {
	resp, err := http.Get(url + "/v1/symbols")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln()
	}

	log.Println(string(body))
}

func getTicker(url string, symbol string) Ticker {
	resp, err := http.Get(url + "/v1/pubticker/" + symbol)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln()
	}

	var t Ticker
	json.Unmarshal(body, &t)

	log.Println(t)
	return t
}

func newOrder(symbol string, amount string, price string) {
	// execute maker buy, round to 8 decimal places for precision, multiply price by 2 so your limit order always gets fully filled
	geminiApiKey := os.Getenv("API_KEY")
	geminiApiSecret := []byte(os.Getenv("API_SECRET"))

	payloadNonce := makeTimestamp()

	requestBody, err := json.Marshal(Order{
		Request: "/v1/order/new",
		Nonce:   payloadNonce,
		Symbol:  symbol,
		Amount:  amount,
		Price:   price,
		Side:    "buy",
		Type:    "exchange limit",
		Options: []string{"maker-or-cancel"},
	})
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", sandboxUrl+"/v1/order/new", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	encodedPayloadAsString := base64.StdEncoding.EncodeToString(requestBody)
	encodedPayload := Base64Encode(requestBody)
	signature := computeHmacInHex(encodedPayload, geminiApiSecret)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-GEMINI-APIKEY", geminiApiKey)
	req.Header.Set("X-GEMINI-PAYLOAD", encodedPayloadAsString)
	req.Header.Set("X-GEMINI-SIGNATURE", signature)
	req.Header.Set("Cache-Control", "no-cache")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln()
	}

	log.Println(string(body))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// getAllTickers(sandboxUrl)
	ticker := getTicker(sandboxUrl, tickerSymbol)

	priceAsString := fmt.Sprintf("%.2f", ticker.Ask*.999)
	price, err := strconv.ParseFloat(priceAsString, 64)
	if err != nil {
		log.Fatalln(err)
	}
	// must round to precision of currency and *.999 for fee inclusion
	// BTC is 8 decimals, ETH is 6, etc
	amount := fmt.Sprintf(precision, (buySize*.999)/price)

	newOrder(tickerSymbol, amount, priceAsString)
}
