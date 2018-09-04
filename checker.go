package main

import (
	"net/http"
	"time"
)


type info struct {
	url string
	status string
}

func main() {
	// Url - Status Channel
	uChan := make(chan info)

	uURls := []info{
		{"google.com",""},
		{"golang.org",""},
		{"wikipedia.com",""},

	}

	for _, i := range uURls{
		go checkUrlLife(uChan, i.url)
	}

	for {
		go func(i info) {
			time.Sleep(time.Second)
			checkUrlLife(uChan,i.url)
		}(<-uChan)
	}

}

func checkUrlLife(ch chan info,  url string) {
	r, err := http.Get("http://"+url)

	in := info{url,r.Status}
	if err != nil{
		println("💢",err.Error())
		ch <- info{url, err.Error()}
		return
	}
	println("✅ ",in.status, in.url)
	ch <- info{url,r.Status}
}

