package model

import "github.com/irisnet/rainbow-sync/db"

var (
	SyncTaskModel    SyncTask
	BlockModel       Block
	TxModel          IrisTx
	AssetDetailModel IrisAssetDetail

	Collections = []db.Docs{
		SyncTaskModel,
		BlockModel,
		TxModel,
		AssetDetailModel,
	}
)

func EnsureDocsIndexes() {
	if len(Collections) > 0 {
		for _, v := range Collections {
			v.EnsureIndexes()
		}
	}
}
