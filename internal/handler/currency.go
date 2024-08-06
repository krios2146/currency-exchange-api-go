package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/krios2146/currency-exchange-api-go/internal/response"
	"github.com/krios2146/currency-exchange-api-go/internal/store"
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

}

func (c *CurrencyHandler) AddCurrency(w http.ResponseWriter, r *http.Request) {
	slog.Debug("POST /currency was called")
}
