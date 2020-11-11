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
	cdc    *codec.Codec
	LcdURL string
	logger tmlog.Logger
}

func NewClient(ctx Context) *Client {
	lcdURL := ctx.LcdURL()

	if lcdURL == "" {
		panic(errors.New("url error"))
	}

	if lcdURL[len(lcdURL)-1] != '/' {
		lcdURL += "/"
	}

	return &Client{
		LcdURL: lcdURL,
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

func (c *Client) Query(format string, a ...interface{}) ([]byte, error) {
	path := c.LcdURL + fmt.Sprintf(format, a...)

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

func (c *Client) queryFromJSON(res interface{}, format string, args ...interface{}) error {
	c.logger.Debug("query", "path", fmt.Sprintf(format, args...))

	data, err := c.Query(format, args...)
	if err != nil {
		return err
	}

	//c.logger.Debug("get %s", string(data))

	//fmt.Println(string(data))

	if err := c.cdc.UnmarshalJSON(data, res); err != nil {
		return errors.Wrapf(err, "unmarshal json err by query %s", fmt.Sprintf(format, args...))
	}

	return nil
}
