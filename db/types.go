// interface for a document

package db

const (
	CollectionNameTxn = "sync_mgo_txn"
	ExistError        = "record exist"
)

type (
	Docs interface {
		// collection name
		Name() string
		// primary key pair(used to find a unique record)
		PkKvPair() map[string]interface{}
		// ensure indexes
		EnsureIndexes()
	}
)
