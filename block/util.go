package block

import "github.com/kaifei-bianjie/msg-parser/types"
import m "github.com/kaifei-bianjie/msg-parser/modules"

func parseDenoms(coins []types.Coin) []string {
	if len(coins) == 0 {
		return nil
	}
	var denoms []string
	for _, v := range coins {
		denoms = append(denoms, v.Denom)
	}

	return denoms
}

// convert coins defined in modules to coins defined in types
func convertCoins(mCoins []m.Coin) types.Coins {
	var coins types.Coins
	if len(mCoins) == 0 {
		return coins
	}
	for _, v := range mCoins {
		coins = append(coins, types.Coin{
			Denom:  v.Denom,
			Amount: v.Amount,
		})
	}

	return coins
}
