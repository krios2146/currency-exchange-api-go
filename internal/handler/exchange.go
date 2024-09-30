package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/krios2146/currency-exchange-api-go/internal/response"
	"github.com/krios2146/currency-exchange-api-go/internal/store"
	"github.com/krios2146/currency-exchange-api-go/internal/validator"
)

type ExchangeHandler struct {
	exchangeRateStore *store.ExchangeRateStore
	currencyStore     *store.CurrencyStore
}

func NewExchangeHandler(exchangeRateStore *store.ExchangeRateStore, currencyStore *store.CurrencyStore) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeRateStore: exchangeRateStore,
		currencyStore:     currencyStore,
	}
}

func (c *ExchangeHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	slog.Debug("GET /exchange was called")

	query := r.URL.Query()

	baseCurrencyCode := query.Get("from")
	targetCurrencyCode := query.Get("to")
	amountStr := query.Get("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{
			Message: fmt.Sprintf("Couldn't parse amount from '%s'", amountStr),
		})
		return
	}
	if err := validator.ValidateCurrencyCode(baseCurrencyCode); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}
	if err := validator.ValidateCurrencyCode(targetCurrencyCode); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	if amount < 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: "Amount cannot be negative"})
		return
	}

	baseCurrency, berr := c.currencyStore.FindByCode(baseCurrencyCode)
	targetCurrency, terr := c.currencyStore.FindByCode(targetCurrencyCode)

	if errors.Is(berr, store.CurrencyNotFoundError) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: berr.Error()})
		return
	}
	if berr != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: berr.Error()})
		return
	}
	if errors.Is(terr, store.CurrencyNotFoundError) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: terr.Error()})
		return
	}
	if terr != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: terr.Error()})
		return
	}

	// Direct exchange
	exchangeRate, err := c.exchangeRateStore.FindByCurrencyCodes(baseCurrencyCode, targetCurrencyCode)

	if exchangeRate != nil {
		convertedAmount := amount * exchangeRate.Rate

		exchangeResponse := response.Exchange{
			BaseCurrency:    *baseCurrency,
			TargetCurrency:  *targetCurrency,
			Amount:          amount,
			Rate:            exchangeRate.Rate,
			ConvertedAmount: round(convertedAmount, 2),
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(exchangeResponse)
		return
	}

	// Indirect exchange
	exchangeRate, err = c.exchangeRateStore.FindByCurrencyCodes(targetCurrencyCode, baseCurrencyCode)

	if exchangeRate != nil {
		convertedAmount := amount * (1 / exchangeRate.Rate)

		exchangeResponse := response.Exchange{
			BaseCurrency:    *baseCurrency,
			TargetCurrency:  *targetCurrency,
			Amount:          amount,
			Rate:            (1 / exchangeRate.Rate),
			ConvertedAmount: round(convertedAmount, 2),
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(exchangeResponse)
		return
	}

	// Cross exchange
	usdToBaseExchangeRate, err := c.exchangeRateStore.FindByCurrencyCodes("USD", baseCurrencyCode)
	usdToTargetExchangeRate, err := c.exchangeRateStore.FindByCurrencyCodes("USD", targetCurrencyCode)

	if usdToBaseExchangeRate != nil && usdToTargetExchangeRate != nil {
		convertedAmount := amount * (usdToTargetExchangeRate.Rate / usdToBaseExchangeRate.Rate)

		exchangeResponse := response.Exchange{
			BaseCurrency:    *baseCurrency,
			TargetCurrency:  *targetCurrency,
			Amount:          amount,
			Rate:            (usdToTargetExchangeRate.Rate / usdToBaseExchangeRate.Rate),
			ConvertedAmount: round(convertedAmount, 2),
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(exchangeResponse)
		return
	}

	if errors.Is(err, store.ExchangeRateNotFoundError) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
	return
}

func round(value float64, precision int) float64 {
	return math.Round(value*math.Pow10(precision)) / math.Pow10(precision)
}
