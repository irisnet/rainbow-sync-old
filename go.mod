module github.com/irisnet/rainbow-sync

        go 1.13

        require (
        github.com/irisnet/irishub v0.16.0
        github.com/jolestar/go-commons-pool v2.0.0+incompatible
        github.com/tendermint/tendermint v0.32.8
        go.uber.org/zap v1.13.0
        gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
        gopkg.in/natefinch/lumberjack.v2 v2.0.0
        gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
        )

        replace (
        github.com/tendermint/iavl => github.com/irisnet/iavl v0.12.3
        github.com/tendermint/tendermint => github.com/irisnet/tendermint v0.32.1
        golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
        )
