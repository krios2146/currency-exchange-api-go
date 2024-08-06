package store

import (
	"database/sql"
	"log/slog"

	"github.com/krios2146/currency-exchange-api-go/internal/model"
)

type CurrencyStore struct {
	db *sql.DB
}

func NewCurrencyStore(db *sql.DB) *CurrencyStore {
	return &CurrencyStore{
		db: db,
	}
}

func (s *CurrencyStore) FindAll() ([]model.Currency, error) {
	rows, err := s.db.Query("SELECT * FROM Currencies;")

	if err != nil {
		slog.Error("SQL Query execution failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	var currencies []model.Currency

	for rows.Next() {
		var currency model.Currency
		err := rows.Scan(&currency.Id, &currency.Code, &currency.FullName, &currency.Sign)

		if err != nil {
			slog.Error("Unable to map row to model", "error", err)
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}
