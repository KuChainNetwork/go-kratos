package main

import (
	"flag"
	"fmt"

	sdk "github.com/KuChainNetwork/go-kratos"
)

const (
	urlLCD = "http://10.1.1.30:1317/"
)

var (
	num = flag.Int64("num", 1, "block height to query")
)

func main() {
	flag.Parse()

	cli := sdk.NewClient(urlLCD)
	data, err := cli.QueryBlockByNum(*num)
	if err != nil {
		panic(fmt.Errorf("data err by %s", err.Error()))
	}

	fmt.Println(string(sdk.MustMarshalJSON(data.DecodeBlock)))
}
