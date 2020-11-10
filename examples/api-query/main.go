package main

import (
	"flag"
	"fmt"

	sdk "github.com/KuChainNetwork/go-kratos"
	chainCfg "github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/utils/log"
)

var (
	num = flag.Int64("num", 0, "block height to query")
	url = flag.String("url", "http://127.0.0.1:1317/", "kuchain rpc url")
)

func main() {
	flag.Parse()

	chainCfg.SealChainConfig()

	cli := sdk.NewClient(*url)
	cli.SetLogger(log.NewLoggerByZap(true, "*:debug"))

	data, err := cli.QueryFullBlock(*num)

	if err != nil {
		panic(fmt.Errorf("data err by %s", err.Error()))
	}

	fmt.Println(string(sdk.MustMarshalJSON(data)))
}
