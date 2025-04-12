package pools

import (
	"fmt"
	"math"
)

// PriceFromTick calcula o preço da cripto em USDC (AAVE/USDC)
func PriceFromTick(tick int) float64 {
	return 1 / math.Pow(1.0001, float64(tick))
}

// Inverte o preço (de AAVE/USDC para USDC/AAVE)
func InvertPrice(p float64) float64 {
	if p == 0 {
		return 0
	}
	return 1 / p
}

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
	Token1USDPrice  float64 // valor de AAVE em dólar
}

func FormatPoolSummary(data PoolData) string {
	// Faixas
	priceMin := PriceFromTick(data.TickUpper)
	priceMax := PriceFromTick(data.TickLower)

	// Invertido
	priceMinInv := InvertPrice(priceMin)
	priceMaxInv := InvertPrice(priceMax)

	// Liquidez efetiva (em tokens)
	liq0 := data.DepositedToken0 - data.WithdrawnToken0
	liq1 := data.DepositedToken1 - data.WithdrawnToken1

	// Valor atual da posição (considerando o valor da cripto)
	capital := liq0 + (liq1 * data.Token1USDPrice)

	// Tarifas em USD
	feesUSD := data.CollectedFees0 + (data.CollectedFees1 * data.Token1USDPrice)

	// String final
	return fmt.Sprintf(`
📊 Resumo da Pool %s/%s

🔸 Faixa: %.5f → %.5f %s/%s
🔸 Faixa: %.2f → %.2f %s/%s

🔸 Posição estimada: ≈ %.2f US$
🔸 Tarifas acumuladas: ≈ %.2f US$
🔸 🧮 Liquidez bruta (tokens): %.3f %s + %.4f %s
🔸 🧮 Liquidez USD: %.3f USDC + %.2f US$ (%.4f %s × %.2f US$)
`, data.Token0Symbol, data.Token1Symbol,
		priceMin, priceMax, data.Token1Symbol, data.Token0Symbol,
		priceMaxInv, priceMinInv, data.Token0Symbol, data.Token1Symbol,
		capital, feesUSD,
		liq0, data.Token0Symbol, liq1, data.Token1Symbol,
		liq0, liq1*data.Token1USDPrice, liq1, data.Token1Symbol, data.Token1USDPrice)
}
