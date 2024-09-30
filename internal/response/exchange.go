package response

import "github.com/krios2146/currency-exchange-api-go/internal/model"

type Exchange struct {
	BaseCurrency    model.Currency `json:"baseCurrency"`
	TargetCurrency  model.Currency `json:"targetCurrency"`
	Rate            float64        `json:"rate"`
	Amount          float64        `json:"amount"`
	ConvertedAmount float64        `json:"convertedAmount"`
}
