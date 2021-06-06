package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func getAllTickers(url string) string {
	resp, err := http.Get(url + "/v1/symbols")
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
