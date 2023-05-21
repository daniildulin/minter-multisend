package multisend

import "math/big"

type Recipient struct {
	Value   *big.Int `json:"value"`
	Address string   `json:"address"`
	CoinID  uint64   `json:"coin_id"`
}
