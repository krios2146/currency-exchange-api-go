package response

import "github.com/krios2146/currency-exchange-api-go/internal/model"

type ExchangeRate struct {
	Id             int64          `json:"id"`
	BaseCurrency   model.Currency `json:"baseCurrency"`
	TargetCurrency model.Currency `json:"targetCurrency"`
	Rate           float64        `json:"rate"`
}
