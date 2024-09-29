package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/krios2146/currency-exchange-api-go/internal/response"
	"github.com/krios2146/currency-exchange-api-go/internal/store"
	"github.com/krios2146/currency-exchange-api-go/internal/validator"
)

type CurrencyHandler struct {
	store *store.CurrencyStore
}

func NewCurrencyHandler(store *store.CurrencyStore) *CurrencyHandler {
	return &CurrencyHandler{
		store: store,
	}
}

func (c *CurrencyHandler) GetAllCurrencies(w http.ResponseWriter, r *http.Request) {
	slog.Debug("GET /currencies was called")

	currencies, err := c.store.FindAll()

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currencies)
}

func (c *CurrencyHandler) GetCurrencyByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	slog.Debug("GET /currency/{code} was called with", "code", code)

	if err := validator.ValidateCurrencyCode(code); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	currency, err := c.store.FindByCode(code)

	if errors.Is(err, store.CurrencyNotFoundError) {
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currency)
}

func (c *CurrencyHandler) AddCurrency(w http.ResponseWriter, r *http.Request) {
	slog.Debug("POST /currency was called")

	if err := r.ParseForm(); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}

	name := r.Form.Get("name")
	code := r.Form.Get("code")
	sign := r.Form.Get("sign")

	if len(name) == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: "Currency name is not present in the request"})
		return
	}

	if err := validator.ValidateCurrencyCode(code); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: err.Error()})
		return
	}
	if len(sign) == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response.ErrorResponse{Message: "Currency sign is not present in the request"})
		return
	}

	currency, err := c.store.Save(name, code, sign)

	if errors.Is(err, store.CurrencyAlreadyExistsError) {
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(currency)
}
