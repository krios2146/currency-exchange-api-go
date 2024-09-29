package store

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/krios2146/currency-exchange-api-go/internal/model"
)

type ExchangeRateStore struct {
	db *sql.DB
}

var ExchangeRateNotFoundError error = errors.New("Exchange rate not found")
var ExchangeRateAlreadyExistsError error = errors.New("Exchange rate already exists")

func NewExchangeRateStore(db *sql.DB) *ExchangeRateStore {
	return &ExchangeRateStore{
		db: db,
	}
}

func (s *ExchangeRateStore) FindAll() ([]model.ExchangeRate, error) {
	rows, err := s.db.Query("SELECT id, base_currency_id, target_currency_id, rate FROM Exchange_rates;")
	defer rows.Close()

	if err != nil {
		slog.Error("SQL Query execution failed", "error", err)
		return nil, err
	}

	var exchangeRates []model.ExchangeRate

	for rows.Next() {
		var exchangeRate model.ExchangeRate
		err := rows.Scan(
			&exchangeRate.Id,
			&exchangeRate.BaseCurrencyId,
			&exchangeRate.TargetCurrencyId,
			&exchangeRate.Rate,
		)

		if err != nil {
			slog.Error("Unable to map row to model", "error", err)
			return nil, err
		}

		exchangeRates = append(exchangeRates, exchangeRate)
	}

	return exchangeRates, nil
}

func (s *ExchangeRateStore) FindByCurrencyCodes(baseCurrencyCode string, targetCurrencyCode string) (*model.ExchangeRate, error) {
	row := s.db.QueryRow(
		`SELECT er.id, er.base_currency_id, er.target_currency_id, er.rate FROM Exchange_rates er
		JOIN Currencies bc ON bc.id = er.base_currency_id
		JOIN Currencies tc ON tc.id = er.target_currency_id
		WHERE bc.code = ? AND tc.code = ?`,
		baseCurrencyCode, targetCurrencyCode,
	)

	var exchangeRate model.ExchangeRate

	err := row.Scan(
		&exchangeRate.Id,
		&exchangeRate.BaseCurrencyId,
		&exchangeRate.TargetCurrencyId,
		&exchangeRate.Rate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ExchangeRateNotFoundError
	}

	if err != nil {
		slog.Error("Unable to map row to model", "error", err)
		return nil, err
	}

	return &exchangeRate, nil
}
