package sdk

import (
	"github.com/KuChainNetwork/kuchain/app"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Client struct {
	cdc     *codec.Codec
	NodeURL string
}

func NewClient() *Client {
	return &Client{
		cdc: app.MakeCodec(),
	}
}
