package sdk

import (
	"github.com/KuChainNetwork/go-kratos/types"
)

func (c *Client) QueryBlockResultsByHeight(height int64) (types.ResultBlockResults, error) {
	var res struct {
		Res types.ResultBlockResults `json:"result"`
	}

	if err := c.queryFromJSON(&res, c.RpcURL, "block_results?height=%d", height); err != nil {
		return types.ResultBlockResults{}, err
	}

	return res.Res, nil
}
