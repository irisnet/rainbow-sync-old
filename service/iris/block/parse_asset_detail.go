package block

import (
	"github.com/irisnet/rainbow-sync/service/iris/logger"
	model "github.com/irisnet/rainbow-sync/service/iris/db"
	imodel "github.com/irisnet/rainbow-sync/service/iris/model"
	"github.com/irisnet/rainbow-sync/service/iris/helper"
	"strings"
	"gopkg.in/mgo.v2/txn"
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
)

var (
	assetDetailTriggers = map[string]bool{
		"stakeEndBlocker":   true,
		"slashBeginBlocker": true,
		"slashEndBlocker":   true,
		"govEndBlocker":     true,
	}

	// adapt multiple asset
	assetDenoms = []string{"iris-atto"}
)

const (
	triggerTxHashLength = 64
	separator           = "::" // tag value separator
	triggerTx           = "tx"
	unDelegationSubject = "Undelegation"
	IRIS                = "Iris"
)

type Iris_Block struct{}

func (iris *Iris_Block) Name() string {
	return IRIS
}

func (iris *Iris_Block) SaveDocsWithTxn(blockDoc *imodel.Block, irisAssetDetail []*imodel.IrisAssetDetail, irisTxs []*imodel.IrisTx, taskDoc imodel.SyncTask) error {
	var (
		ops, irisAssetDetailOps, irisTxsOps []txn.Op
	)

	if blockDoc.Height == 0 {
		return fmt.Errorf("invalid block, height equal 0")
	}

	blockOp := txn.Op{
		C:      imodel.CollectionNameBlock,
		Id:     bson.NewObjectId(),
		Insert: blockDoc,
	}

	length_assetdetail := len(irisAssetDetail)
	if length_assetdetail > 0 {
		irisAssetDetailOps = make([]txn.Op, 0, length_assetdetail)
		for _, v := range irisAssetDetail {
			op := txn.Op{
				C:      imodel.CollectionNameAssetDetail,
				Id:     bson.NewObjectId(),
				Insert: v,
			}
			irisAssetDetailOps = append(irisAssetDetailOps, op)
		}
	}
	length_txs := len(irisTxs)
	if length_txs > 0 {
		irisTxsOps = make([]txn.Op, 0, length_txs)
		for _, v := range irisTxs {
			op := txn.Op{
				C:      imodel.CollectionNameIrisTx,
				Id:     bson.NewObjectId(),
				Insert: v,
			}
			irisTxsOps = append(irisTxsOps, op)
		}
	}

	updateOp := txn.Op{
		C:      imodel.CollectionNameSyncTask,
		Id:     taskDoc.ID,
		Assert: txn.DocExists,
		Update: bson.M{
			"$set": bson.M{
				"current_height":   taskDoc.CurrentHeight,
				"status":           taskDoc.Status,
				"last_update_time": taskDoc.LastUpdateTime,
			},
		},
	}

	ops = make([]txn.Op, 0, length_assetdetail+length_txs+2)
	ops = append(append(ops, blockOp, updateOp), irisAssetDetailOps...)
	ops = append(ops, irisTxsOps...)

	if len(ops) > 0 {
		err := model.Txn(ops)
		if err != nil {
			return err
		}
	}

	return nil
}

func (iris *Iris_Block) ParseBlock(b int64, client *helper.Client) (resBlock *imodel.Block, resIrisAssetDetails []*imodel.IrisAssetDetail, resIrisTxs []*imodel.IrisTx, resErr error) {

	defer func() {
		if err := recover(); err != nil {
			logger.Error("parse iris block fail", logger.Int64("height", b),
				logger.Any("err", err), logger.String("Chain Block", iris.Name()))

			resBlock = &imodel.Block{}
			resIrisAssetDetails = nil
			resIrisTxs = nil
			resErr = fmt.Errorf("%v", err)
		}
	}()

	irisAssetDetails, err := iris.ParseIrisAssetDetail(b, client)
	if err != nil {
		logger.Error("parse iris asset detail error", logger.String("error", err.Error()), logger.String("Chain Block", iris.Name()))
	}

	irisTxs, err := iris.ParseIrisTxs(b, client)
	if err != nil {
		logger.Error("parse iris txs", logger.String("error", err.Error()), logger.String("Chain Block", iris.Name()))
	}

	resBlock = &imodel.Block{
		Height:     b,
		CreateTime: time.Now().Unix(),
	}
	resIrisAssetDetails = irisAssetDetails
	resIrisTxs = irisTxs
	resErr = err

	return
}

// parse iris asset detail from block result tags
func (iris *Iris_Block) ParseIrisAssetDetail(b int64, client *helper.Client) ([]*imodel.IrisAssetDetail, error) {
	var irisAssetDetails []*imodel.IrisAssetDetail
	res, err := client.BlockResults(&b)
	if err != nil {
		logger.Warn("get block result err, now try again", logger.String("err", err.Error()),
			logger.String("Chain Block", iris.Name()))
		// there is possible parse block fail when in iterator
		var err2 error
		client2 := helper.GetClient()
		res, err2 = client2.BlockResults(&b)
		client2.Release()
		if err2 != nil {
			return nil, err2
		}
	}

	tags := res.Results.EndBlock.Tags
	//fmt.Printf("======>>tags:%+v\n",tags)

	// filter asset detail trigger from tags and build asset detail model
	irisAssetDetails = make([]*imodel.IrisAssetDetail, 0, len(tags))
	for _, t := range tags {
		tagKey := string(t.Key)
		tagValue := string(t.Value)

		if assetDetailTriggers[tagKey] || len(tagKey) == triggerTxHashLength {
			values := strings.Split(tagValue, separator)
			if len(values) != 6 {
				logger.Warn("struct of iris asset detail changed in block result, skip parse this asset detail",
					logger.Int64("height", b), logger.String("tagKey", tagKey),
					logger.String("Chain Block", iris.Name()))
				continue
			}

			irisAssetDetails = append(irisAssetDetails, buildIrisAssetDetailFromTag(tagKey, values, b))
		}
	}

	return irisAssetDetails, nil
}

// get asset detail info by parse tag key and values
func buildIrisAssetDetailFromTag(tagKey string, keyValues []string, height int64) *imodel.IrisAssetDetail {
	values := keyValues
	coinAmount, coinUnit := parseCoinAmountAndUnitFromStr(values[2])

	irisAssetDetail := &imodel.IrisAssetDetail{
		From:        values[0],
		To:          values[1],
		CoinAmount:  coinAmount,
		CoinUnit:    coinUnit,
		Trigger:     tagKey,
		Subject:     values[3],
		Description: values[4],
		Timestamp:   values[5],
		Height:      height,
	}

	if len(tagKey) == triggerTxHashLength {
		irisAssetDetail.TxHash = tagKey
		irisAssetDetail.Trigger = triggerTx
	}

	if irisAssetDetail.Subject == unDelegationSubject {
		irisAssetDetail.TxHash = irisAssetDetail.Description
	}

	irisAssetDetail.TxHash = strings.ToUpper(irisAssetDetail.TxHash)
	return irisAssetDetail
}

func parseCoinAmountAndUnitFromStr(s string) (string, string) {
	for _, denom := range assetDenoms {
		if strings.HasSuffix(s, denom) {
			return strings.Replace(s, denom, "", -1), denom
		}
	}
	return "", ""
}
