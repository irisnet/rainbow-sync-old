module github.com/irisnet/rainbow-sync/service/iris

		go 1.13

		require (
		github.com/btcsuite/btcd v0.0.0-20190807005414-4063feeff79a // indirect
		github.com/cosmos/cosmos-sdk v0.34.4-0.20191011012532-e69029151525
		github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d // indirect
		github.com/golang/mock v1.3.1 // indirect
		github.com/irisnet/cosmos-sdk v0.23.1 // indirect
		github.com/jolestar/go-commons-pool v2.0.0+incompatible
		github.com/onsi/ginkgo v1.10.2 // indirect
		github.com/onsi/gomega v1.7.0 // indirect
		github.com/prometheus/client_golang v1.1.0 // indirect
		github.com/rcrowley/go-metrics v0.0.0-20190706150252-9beb055b7962 // indirect
		github.com/spf13/afero v1.2.2 // indirect
		github.com/tendermint/tendermint v0.32.6
		go.uber.org/zap v1.10.0
		golang.org/x/crypto v0.0.0-20190909091759-094676da4a83 // indirect
		golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
		golang.org/x/text v0.3.2 // indirect
		google.golang.org/genproto v0.0.0-20190801165951-fa694d86fc64 // indirect
		gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
		gopkg.in/natefinch/lumberjack.v2 v2.0.0
		gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
		)

		replace (
		github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.23.2-0.20191015002325-ee5d7f3d62d9
		github.com/tendermint/tendermint => github.com/tendermint/tendermint v0.32.4
		)
