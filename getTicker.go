package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Ticker struct {
	Bid    float64     `json:"bid,string"`
	Ask    float64     `json:"ask,string"`
	Volume interface{} `json:"volume"`
	Last   float64     `json:"last,string"`
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
	return t
}
