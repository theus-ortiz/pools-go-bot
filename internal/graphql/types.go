package graphql

type PositionRaw struct {
	ID                 string `json:"id"`
	DepositedToken0    string `json:"depositedToken0"`
	DepositedToken1    string `json:"depositedToken1"`
	WithdrawnToken0    string `json:"withdrawnToken0"`
	WithdrawnToken1    string `json:"withdrawnToken1"`
	AmountDepositedUSD string `json:"amountDepositedUSD"`
	AmountWithdrawnUSD string `json:"amountWithdrawnUSD"`
	Liquidity          string `json:"liquidity"`
}

type SubgraphResponse struct {
	Data struct {
		Positions []PositionRaw `json:"positions"`
	} `json:"data"`
}

type PositionDetailed struct {
	ID                  string `json:"id"`
	Liquidity           string `json:"liquidity"`
	DepositedToken0     string `json:"depositedToken0"`
	DepositedToken1     string `json:"depositedToken1"`
	WithdrawnToken0     string `json:"withdrawnToken0"`
	WithdrawnToken1     string `json:"withdrawnToken1"`
	CollectedFeesToken0 string `json:"collectedFeesToken0"`
	CollectedFeesToken1 string `json:"collectedFeesToken1"`
	TickLower           struct {
		TickIdx string `json:"tickIdx"`
	} `json:"tickLower"`
	TickUpper struct {
		TickIdx string `json:"tickIdx"`
	} `json:"tickUpper"`
	Pool struct {
		ID        string `json:"id"`
		FeeTier   string `json:"feeTier"`
		Tick      string `json:"tick"`
		SqrtPrice string `json:"sqrtPrice"`
		Token0    struct {
			Symbol   string `json:"symbol"`
			Decimals string `json:"decimals"`
		} `json:"token0"`
		Token1 struct {
			Symbol   string `json:"symbol"`
			Decimals string `json:"decimals"`
		} `json:"token1"`
	} `json:"pool"`
}

type SubgraphPositionByIDResponse struct {
	Data struct {
		Position PositionDetailed `json:"position"`
	} `json:"data"`
}