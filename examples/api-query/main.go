package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	urlLCD = "http://127.0.0.1:1317/"
)

func main() {
	resp, err := http.Get(urlLCD)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 {
		fmt.Println("ok")
	}
}
