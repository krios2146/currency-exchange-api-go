package api

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/krios2146/currency-exchange-api-go/internal/handler"
	"github.com/krios2146/currency-exchange-api-go/internal/store"
)

type Server struct {
	db *sql.DB
}

type CurrenciesHandler struct {
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	slog.Debug("Registering handlers")

	currencyStore := store.NewCurrencyStore(s.db)
	currencyHandler := handler.NewCurrencyHandler(currencyStore)

	mux.HandleFunc("GET /currencies", currencyHandler.GetAllCurrencies)
	mux.HandleFunc("GET /currency/{code}", currencyHandler.GetCurrencyByCode)
	mux.HandleFunc("GET /currency/", currencyHandler.GetCurrencyByCode)
	mux.HandleFunc("POST /currency", currencyHandler.AddCurrency)

	slog.Info("Starting server")

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
