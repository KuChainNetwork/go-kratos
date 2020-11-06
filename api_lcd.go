package sdk

import (
	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/KuChainNetwork/kuchain/chain/client/rest/block"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (c *Client) QueryBlockByNum(num int64) (block.DecodeResultBlock, error) {
	var res block.DecodeResultBlock

	if err := c.queryFromJSON(&res, "blocks/%d/decode", num); err != nil {
		return res, errors.Wrapf(err, "error by query block decode %d", num)
	}

	return res, nil
}

func (c *Client) QueryTxByHash(hash string) (sdk.TxResponse, error) {
	var tx sdk.TxResponse

	if err := c.queryFromJSON(&tx, "txs/%s", hash); err != nil {
		return tx, errors.Wrapf(err, "error by query tx %s", hash)
	}

	return tx, nil
}

func (c *Client) QueryFullBlock(num int64) (types.FullBlock, error) {
	var decodeRes block.DecodeResultBlock

	if err := c.queryFromJSON(&decodeRes, "blocks/%d/decode", num); err != nil {
		return types.FullBlock{}, errors.Wrapf(err, "error by query block decode %d", num)
	}

	res := types.FullBlock{
		Header:     decodeRes.DecodeBlock.Header,
		BlockID:    decodeRes.BlockID.Hash.String(),
		Evidence:   decodeRes.DecodeBlock.Evidence,
		LastCommit: *decodeRes.DecodeBlock.LastCommit,
		Txs:        make([]types.TxData, 0, len(decodeRes.DecodeBlock.TxsHash)),
	}

	// get each tx
	for _, txHash := range decodeRes.DecodeBlock.TxsHash {
		txRes, err := c.QueryTxByHash(txHash.String())
		if err != nil {
			return types.FullBlock{}, errors.Wrapf(err, "error by tx query %s", txHash.String())
		}

		res.Txs = append(res.Txs, types.TxData(txRes))
	}

	return res, nil
}
