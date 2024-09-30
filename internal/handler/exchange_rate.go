package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/krios2146/currency-exchange-api-go/internal/response"
	"github.com/krios2146/currency-exchange-api-go/internal/store"
	"github.com/krios2146/currency-exchange-api-go/internal/validator"
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

func (c *ExchangeRateHandler) GetExchangeRateByCodes(w http.ResponseWriter, r *http.Request) {
	codePair := r.PathValue("code_pair")

	slog.Debug("GET /exchangeRate/{code_pair} was called, with", "code_pair", codePair)

	if len(codePair) != 6 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: "Code pair must contain exactly 6 letters"})
		return
	}

	baseCurrencyCode := codePair[0:3]
	targetCurrencyCode := codePair[3:6]

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

	exchangeRate, err := c.exchangeRateStore.FindByCurrencyCodes(baseCurrencyCode, targetCurrencyCode)

	if errors.Is(err, store.ExchangeRateNotFoundError) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exchangeRateResponse)
}

func (c *ExchangeRateHandler) AddExchangeRate(w http.ResponseWriter, r *http.Request) {
	slog.Debug("POST /exchangeRates was called")

	if err := r.ParseForm(); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	baseCurrencyCode := r.Form.Get("baseCurrencyCode")
	targetCurrencyCode := r.Form.Get("targetCurrencyCode")
	rateStr := r.Form.Get("rate")
	rate, err := strconv.ParseFloat(rateStr, 64)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: fmt.Sprintf("Couldn't parse rate from '%s'", rateStr)})
		return
	}

	if err := validator.ValidateCurrencyCode(baseCurrencyCode); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}
	if err := validator.ValidateCurrencyCode(targetCurrencyCode); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	if rate <= 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: "Rate cannot be negative or zero"})
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

	exchangeRate, err := c.exchangeRateStore.Save(baseCurrency.Id, targetCurrency.Id, rate)

	if errors.Is(err, store.ExchangeRateAlreadyExistsError) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	if err != nil {
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exchangeRateResponse)
}

func (c *ExchangeRateHandler) UpdateExchangeRate(w http.ResponseWriter, r *http.Request) {
	codePair := r.PathValue("code_pair")

	slog.Debug("PATCH /exchangeRate/{code_pair} was called, with", "code_pair", codePair)

	if len(codePair) != 6 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: "Code pair must contain exactly 6 letters"})
		return
	}

	baseCurrencyCode := codePair[0:3]
	targetCurrencyCode := codePair[3:6]

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

	if err := r.ParseForm(); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	rateStr := r.Form.Get("rate")
	rate, err := strconv.ParseFloat(rateStr, 64)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: fmt.Sprintf("Couldn't parse rate from '%s'", rateStr)})
		return
	}

	if rate <= 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: "Rate cannot be negative or zero"})
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

	exchangeRate, err := c.exchangeRateStore.Update(baseCurrency.Id, targetCurrency.Id, rate)

	if errors.Is(err, store.ExchangeRateNotFoundError) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	if err != nil {
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exchangeRateResponse)
}
