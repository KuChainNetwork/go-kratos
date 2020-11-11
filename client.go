package sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/KuChainNetwork/kuchain/app"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/pkg/errors"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

var (
	Codec *codec.Codec
)

func init() {
	Codec = app.MakeCodec()
}

type Client struct {
	cdc *codec.Codec

	LcdURL string
	RpcURL string
	logger tmlog.Logger
}

func NewClient(ctx Context) *Client {
	lcdURL := ctx.LcdURL()
	if lcdURL != "" {
		if lcdURL[len(lcdURL)-1] != '/' {
			lcdURL += "/"
		}
	}

	rpcURL := ctx.RpcURL()
	if rpcURL != "" {
		if rpcURL[len(lcdURL)-1] != '/' {
			rpcURL += "/"
		}
	}

	return &Client{
		LcdURL: lcdURL,
		RpcURL: rpcURL,

		cdc:    app.MakeCodec(),
		logger: ctx.Logger(),
	}
}

func (c *Client) SetLogger(l tmlog.Logger) {
	c.logger = l
}

func (c Client) Cdc() *codec.Codec {
	return c.cdc
}

func (c *Client) query(url, format string, a ...interface{}) ([]byte, error) {
	path := url + fmt.Sprintf(format, a...)

	resp, err := http.Get(path)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "error by get with %s", path)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "error by read all with %s", path)
	}

	if resp.StatusCode == 200 {
		return body, nil
	} else {
		return []byte{}, fmt.Errorf("resp code by %d with %s", resp.StatusCode, path)
	}
}

func (c *Client) queryFromJSON(res interface{}, url, format string, args ...interface{}) error {
	c.logger.Debug("query", "path", fmt.Sprintf(format, args...))

	data, err := c.query(url, format, args...)
	if err != nil {
		return err
	}

	if err := c.cdc.UnmarshalJSON(data, res); err != nil {
		return errors.Wrapf(err, "unmarshal json err by query %s", fmt.Sprintf(format, args...))
	}

	return nil
}
