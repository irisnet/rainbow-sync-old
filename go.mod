module github.com/irisnet/rainbow-sync

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.40.0
	github.com/jolestar/go-commons-pool v2.0.0+incompatible
	github.com/tendermint/tendermint v0.34.1
	github.com/weichang-bianjie/msg-sdk v0.0.0-20210127095957-c6c9a54ce1f2
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20200930145003-4acb6c075d10
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
