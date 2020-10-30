package main

import (
	"encoding/json"
	"flag"
	"fmt"

	sdk "github.com/KuChainNetwork/go-kratos"
)

const (
	urlLCD = "http://127.0.0.1:1231/"
)

var (
	num = flag.Int64("num", 1, "block height to query")
)

func main() {
	flag.Parse()

	cli := sdk.NewClient(urlLCD)
	data, err := cli.QueryBlockByNum(*num)
	if err != nil {
		fmt.Errorf("data err by %s", err.Error())
	}

	datas, _ := json.Marshal(*data.DecodeBlock)
	fmt.Println(string(datas))
}
