package pools

import (
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/theus-ortiz/pools-bot/internal/graphql"
)

type PoolData struct {
	Token0Symbol    string
	Token1Symbol    string
	Decimals0       int
	Decimals1       int
	DepositedToken0 float64
	WithdrawnToken0 float64
	DepositedToken1 float64
	WithdrawnToken1 float64
	CollectedFees0  float64
	CollectedFees1  float64
	TickLower       int
	TickUpper       int
	CurrentTick     int
	Token1USDPrice  float64
}

// Calcula preço a partir de sqrtPrice
func GetPriceFromSqrtPrice(sqrtPriceStr string) float64 {
	sqrtPrice, ok := new(big.Int).SetString(sqrtPriceStr, 10)
	if !ok {
		return 0
	}

	sqrtPriceFloat := new(big.Float).SetInt(sqrtPrice)
	twoPow96 := new(big.Float).SetFloat64(math.Pow(2, 96))
	ratio := new(big.Float).Quo(sqrtPriceFloat, twoPow96)

	price := new(big.Float).Mul(ratio, ratio)
	result, _ := price.Float64()
	return result
}

// Constrói PoolData a partir de PositionDetailed
func BuildPoolDataFromPosition(p graphql.PositionDetailed) PoolData {
	// Converte strings para números
	parse := func(s string) float64 {
		f, _ := strconv.ParseFloat(s, 64)
		return f
	}

	tickLower, _ := strconv.Atoi(p.TickLower.TickIdx)
	tickUpper, _ := strconv.Atoi(p.TickUpper.TickIdx)
	currentTick, _ := strconv.Atoi(p.Pool.Tick)
	dec0, _ := strconv.Atoi(p.Pool.Token0.Decimals)
	dec1, _ := strconv.Atoi(p.Pool.Token1.Decimals)

	price := GetPriceFromSqrtPrice(p.Pool.SqrtPrice)

	var token1USDPrice float64
	if p.Pool.Token0.Symbol == "USDC" {
		token1USDPrice = 1.0 / price
	} else if p.Pool.Token1.Symbol == "USDC" {
		token1USDPrice = 1.0 * price
	} else {
		token1USDPrice = 0 // Nenhum dos tokens é USDC
	}

	return PoolData{
		Token0Symbol:    p.Pool.Token0.Symbol,
		Token1Symbol:    p.Pool.Token1.Symbol,
		Decimals0:       dec0,
		Decimals1:       dec1,
		DepositedToken0: parse(p.DepositedToken0),
		WithdrawnToken0: parse(p.WithdrawnToken0),
		DepositedToken1: parse(p.DepositedToken1),
		WithdrawnToken1: parse(p.WithdrawnToken1),
		CollectedFees0:  parse(p.CollectedFeesToken0),
		CollectedFees1:  parse(p.CollectedFeesToken1),
		TickLower:       tickLower,
		TickUpper:       tickUpper,
		CurrentTick:     currentTick,
		Token1USDPrice:  token1USDPrice,
	}
}

func FormatPoolSummary(p PoolData) string {
	return fmt.Sprintf(
		"💧 Token0 (%s): depositado `%.2f`, retirado `%.2f`, fees `%.2f`\n"+
			"💧 Token1 (%s): depositado `%.2f`, retirado `%.2f`, fees `%.2f`\n"+
			"📈 Ticks: atual `%d`, faixa [`%d`, `%d`]\n"+
			"💵 Preço estimado de Token1: `$%.4f`",
		p.Token0Symbol, p.DepositedToken0, p.WithdrawnToken0, p.CollectedFees0,
		p.Token1Symbol, p.DepositedToken1, p.WithdrawnToken1, p.CollectedFees1,
		p.CurrentTick, p.TickLower, p.TickUpper,
		p.Token1USDPrice,
	)
}
