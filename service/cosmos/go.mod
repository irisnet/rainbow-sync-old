module github.com/irisnet/rainbow-sync/service/cosmos

		go 1.13

		require (
		github.com/cosmos/cosmos-sdk v0.34.4-0.20191011153240-3d5c97e59cc4
		github.com/jolestar/go-commons-pool v2.0.0+incompatible
		github.com/tendermint/tendermint v0.32.6
		go.uber.org/zap v1.10.0
		gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
		gopkg.in/natefinch/lumberjack.v2 v2.0.0
		gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
		)

		replace github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.23.2-0.20191022102555-c1d4d1c8fb5c
