package main

import (
	"encoding/base64"
	"fmt"

	"github.com/tendermint/go-amino"
)

func UnmarshalBinaryBareStr(cdc *amino.Codec, str string, v interface{}) error {
	bz, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return fmt.Errorf("DecodeString error by %s", err.Error())
	}

	if err := cdc.UnmarshalBinaryBare(bz, &v); err != nil {
		return fmt.Errorf("UnmarshalBinaryBare error by %s", err.Error())
	}

	return nil
}
