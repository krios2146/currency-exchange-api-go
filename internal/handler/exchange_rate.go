package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/krios2146/currency-exchange-api-go/internal/response"
	"github.com/krios2146/currency-exchange-api-go/internal/store"
)

type ExchangeRateHandler struct {
	exchangeRateStore *store.ExchangeRateStore
	currencyStore     *store.CurrencyStore
}

func NewExchangeRateHandler(exchangeRateStore *store.ExchangeRateStore, currencyStore *store.CurrencyStore) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		exchangeRateStore: exchangeRateStore,
		currencyStore:     currencyStore,
	}
}

func (c *ExchangeRateHandler) GetAllExchangeRates(w http.ResponseWriter, r *http.Request) {
	slog.Debug("GET /exchangeRates was called")

	exchangeRates, err := c.exchangeRateStore.FindAll()

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	var exchangeRateResponses []response.ExchangeRate

	for _, exchangeRate := range exchangeRates {
		baseCurrency, berr := c.currencyStore.FindById(exchangeRate.BaseCurrencyId)
		targetCurrency, terr := c.currencyStore.FindById(exchangeRate.TargetCurrencyId)

		if berr != nil || terr != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
			return
		}

		exchangeRateResponse := response.ExchangeRate{
			Id:             exchangeRate.Id,
			BaseCurrency:   *baseCurrency,
			TargetCurrency: *targetCurrency,
			Rate:           exchangeRate.Rate,
		}
		exchangeRateResponses = append(exchangeRateResponses, exchangeRateResponse)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exchangeRateResponses)
}
