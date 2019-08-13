package iris

import "gopkg.in/mgo.v2/bson"

type (
	IrisAssetDetail struct {
		From        string `bson:"from"`
		To          string `bson:"to"`
		CoinAmount  string `bson:"coin_amount"`
		CoinUnit    string `bson:"coin_unit"`
		Trigger     string `bson:"trigger"`
		Subject     string `bson:"subject"`
		Description string `bson:"description"`
		Timestamp   string `bson:"timestamp"`
		Height      int64  `bson:"height"`
		TxHash      string `bson:"tx_hash"`
		Ext         string `bson:"ext"`
	}
)

const (
	CollectionNameAssetDetail = "sync_iris_asset_detail"
)

func (d IrisAssetDetail) Name() string {
	return CollectionNameAssetDetail
}

func (d IrisAssetDetail) PkKvPair() map[string]interface{} {
	return bson.M{}
}
