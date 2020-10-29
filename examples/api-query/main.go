package main

import (
	"fmt"

	sdk "github.com/KuChainNetwork/go-kratos"
)

const (
	urlLCD = "http://127.0.0.1:1231/"
)

func main() {
	cli := sdk.NewClient(urlLCD)
	data, err := cli.Query("blocks/%d/decode", 20)
	if err != nil {
		fmt.Errorf("data err by %s", err.Error())
	}
	fmt.Println(string(data))
}
