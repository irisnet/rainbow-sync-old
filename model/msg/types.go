package msg

type (
	CoinStr struct {
		Denom  string `bson:"denom" `
		Amount string `bson:"amount"`
	}
)
