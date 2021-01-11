package model

import "github.com/irisnet/rainbow-sync/db"

var (
	SyncTaskModel SyncTask
	BlockModel    Block
	TxModel       Tx

	Collections = []db.Docs{
		SyncTaskModel,
		BlockModel,
		TxModel,
		new(TxMsg),
	}
)

func EnsureDocsIndexes() {
	if len(Collections) > 0 {
		for _, v := range Collections {
			v.EnsureIndexes()
		}
	}
}
