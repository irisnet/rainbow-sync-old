package msg

type (
	CoinStr struct {
		Amount string `bson:"amount" json:"amount"`
		Denom  string `bson:"denom" json:"denom"`
	}
)
