package main

import (
	"github.com/mijia/sweb/log"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	url = "http://127.0.0.1:8888"

	COUNTS = 100000
)

func get() {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		log.Error(err)
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)
}
func main() {
	start := time.Now()
	for i := 0; i < COUNTS; i++ {
		get()
	}
	costs := time.Now().Sub(start).Seconds()
	log.Info("costs time:", costs, " qps:", (COUNTS / costs))
}
