module github.com/irisnet/rainbow-sync

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.40.0-rc0
	github.com/irisnet/irismod v1.1.1-0.20201015064111-9d21ec6ad6fc
	github.com/jolestar/go-commons-pool v2.0.0+incompatible
	github.com/tendermint/tendermint v0.34.0-rc4.0.20201005135527-d7d0ffea13c6
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.34.4-0.20201014023301-f172e47973d0
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
)
