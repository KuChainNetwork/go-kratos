package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/KuChainNetwork/kuchain/app"
	"github.com/KuChainNetwork/kuchain/chain/types"
	stakingTypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

var (
	NumsLarge map[int][]string = make(map[int][]string)
	NumsLess  map[int][]string = make(map[int][]string)
)

func parseCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse datas marshal by cdc in kuchain",
	}

	cmd.AddCommand(
		ForEachCheckSet(cdc),
	)

	return cmd
}

// {"module": "kuStoreLog", "storeKey": "KVStoreKey{0xc00018f900, kustaking}", "num": 0, "key": "Eg==", "value": "ATA="}
type LogData struct {
	Module   string `json:"module"`
	StoreKey string `json:"storeKey"`
	Num      int    `json:"num"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

func ForEachCheckSet(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "foreach-set [logpath]",
		Short: "foreach set data from log",
		RunE: func(cmd *cobra.Command, args []string) error {

			file, err := os.Open(args[0])
			if err != nil {
				log.Fatal(err)
				return err
			}

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				text := scanner.Bytes()
				// process log

				//log.Printf("%s", string(text))

				var logData LogData
				if err := json.Unmarshal(text, &logData); err != nil {
					log.Fatal(err)
					return err
				}

				processLogData(logData)
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
				return err
			}

			fmt.Printf("\n\nNumsLarge to\n")
			for n, v := range NumsLarge {
				fmt.Printf("NumsLarg %d %v\n", n, v)
			}

			fmt.Printf("\n\\nNumsLess to\n")
			for n, v := range NumsLess {
				fmt.Printf("NumsLess %d %v\n", n, v)
			}

			for n, vv := range NumsLess {
				if vvv, ok := NumsLarge[n]; ok {
					fmt.Printf("Both In %d large %v <--> less %v", n, vvv, vv)
				}
			}

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}

func parseKey(keyStr string) ([]byte, types.AccountID, error) {
	bz, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return []byte{}, types.AccountID{}, fmt.Errorf("hex decode error %s", err.Error())
	}

	if len(bz) <= 1 {
		return []byte{}, types.AccountID{}, fmt.Errorf("len")
	}

	prefixLen := 1

	id := bz[:prefixLen]

	n := types.AccountID{
		Value: bz[prefixLen:],
	}

	return id, n, nil
}

func parseVal(valStr string) (stakingTypes.Validator, error) {
	bz, err := base64.StdEncoding.DecodeString(valStr)
	if err != nil {
		return stakingTypes.Validator{}, fmt.Errorf("hex val decode error %s", err.Error())
	}

	res := stakingTypes.MustUnmarshalValidator(app.MakeCodec(), bz)
	return res, nil
}

func processLogData(logData LogData) {
	// parse key
	id, acc, err := parseKey(logData.Key)
	if err != nil {
		return
	}

	//	log.Printf("id %s %v %s\n", logData.Key, id, acc.String())

	if bytes.Equal(id, stakingTypes.ValidatorsKey) {
		//log.Printf("id %s %v %s\n", logData.Key, id, acc.String())
		//fmt.Printf("key set %s\n", acc.String())

		val, err := parseVal(logData.Value)
		if err != nil {
			panic(err)
		}

		//fmt.Printf("key set %s\n", val.String())

		if len(val.OperatorAccount.Value) == types.AccIDStoreKeyLen {
			fmt.Printf("oper Val is StoreLen  %s %d, %v\n", acc.String(), logData.Num, val.OperatorAccount.String())

			n, ok := NumsLarge[logData.Num]
			if ok {
				n = append(n, val.OperatorAccount.String())
				NumsLarge[logData.Num] = n
			} else {
				NumsLarge[logData.Num] = []string{val.OperatorAccount.String()}
			}
		} else {
			fmt.Printf("oper Val is LesLen %d  %s %d, %v\n", len(val.OperatorAccount.Value), acc.String(), logData.Num, val.OperatorAccount.String())
			n, ok := NumsLess[logData.Num]
			if ok {
				n = append(n, val.OperatorAccount.String())
				NumsLess[logData.Num] = n
			} else {
				NumsLess[logData.Num] = []string{val.OperatorAccount.String()}
			}
		}
	}
}
