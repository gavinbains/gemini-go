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
	"net/http"
	"os"
	"time"
)

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

func makeTimestamp() string {
	return fmt.Sprint(int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond))
}

func Base64Encode(message []byte) []byte {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(b, message)
	return b
}

func computeHmacInHex(message []byte, secret []byte) string {
	var h hash.Hash = hmac.New(sha512.New384, secret)

	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

func placeOrder(symbol string, amount string, price string, url string) string {
	// execute maker buy, round to 8 decimal places for precision, multiply price by 2 so your limit order always gets fully filled
	geminiApiKey := os.Getenv("SANDBOX_API_KEY")
	geminiApiSecret := []byte(os.Getenv("SANDBOX_API_SECRET"))

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

	req, err := http.NewRequest("POST", url+"/v1/order/new", bytes.NewBuffer(requestBody))
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

	return string(body)
}
