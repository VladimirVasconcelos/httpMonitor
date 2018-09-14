package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type info struct {
	url      string
	status   string
	interval time.Duration
}

type monitorConfig struct {
	Urls           []string      `json:"urls"`
	Interval       time.Duration `json:"interval"`
	PeopleInCharge []struct {
		User   string   `json:"user"`
		Phones []string `json:"phones"`
	} `json:"peopleInCharge"`
}

func main() {
	// Url - Status Channel
	uChan := make(chan info)

	file, _ := ioutil.ReadFile("./config.json")

	var config monitorConfig
	err := json.Unmarshal(file, &config)
	if err != nil {
		log.Panic(err.Error())
	}

	cUrls := config.Urls
	for _, urlVal := range cUrls {
		in := info{urlVal, *new(string), config.Interval}
		go checkURLLife(uChan, in)
	}

	for {
		go func(c info) {
			time.Sleep(time.Second * config.Interval)
			checkURLLife(uChan, c)
		}(<-uChan)
	}

}

func checkURLLife(ch chan info, config info) {
	r, err := http.Get("http://" + config.url)

	config.status = r.Status
	if err != nil {
		println("💢", err.Error())
		config.status = err.Error()
		ch <- config
		return
	}
	println("✅ ", config.status, config.url)
	ch <- config
}
