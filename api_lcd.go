package sdk

import (
	"github.com/KuChainNetwork/kuchain/chain/client/rest/block"
	"github.com/pkg/errors"
)

func (c *Client) QueryBlockByNum(num int64) (block.DecodeResultBlock, error) {
	var res block.DecodeResultBlock

	if err := c.queryFromJSON(&res, "blocks/%d/decode", num); err != nil {
		return res, errors.Wrapf(err, "error by query block decode %d", num)
	}

	return res, nil
}
