module github.com/irisnet/rainbow-sync

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200825201020-d9fd4d2ca9a3
	github.com/irismod/coinswap v0.0.0-20200901103745-f38672ac63ec
    github.com/irismod/htlc v0.0.0-20200901103718-c3bf89708dce
    github.com/irismod/nft v1.1.1-0.20200827095318-d16861212579
    github.com/irismod/record v1.1.1-0.20200827095301-3e27fc43ae73
    github.com/irismod/service v1.1.1-0.20200901115916-d898b826bf10
    github.com/irismod/token v1.1.1-0.20200901121217-d3aa04e760e3
	github.com/jolestar/go-commons-pool v2.0.0+incompatible
	github.com/tendermint/tendermint v0.34.0-rc3
	go.uber.org/zap v1.13.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.34.4-0.20200901030027-1e0963031861
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
)
