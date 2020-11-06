module github.com/KuChainNetwork/go-kratos

go 1.15

require (
	github.com/KuChainNetwork/kuchain v0.6.0-rc
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/gorilla/websocket v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.33.8
	go.uber.org/zap v1.16.0
)

replace github.com/KuChainNetwork/kuchain v0.6.0-rc => github.com/KuChainNetwork/kratos v0.6.0-rc
