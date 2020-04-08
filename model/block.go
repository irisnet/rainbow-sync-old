package model

const (
	CollectionNameBlock = "sync_%v_block"
)

type (
	Block struct {
		Height     int64 `bson:"height"`
		CreateTime int64 `bson:"create_time"`
	}
)
