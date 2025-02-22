package store

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/krios2146/currency-exchange-api-go/internal/model"
	"github.com/mattn/go-sqlite3"
)

type CurrencyStore struct {
	db *sql.DB
}

var CurrencyNotFoundError error = errors.New("Currency not found")
var CurrencyAlreadyExistsError error = errors.New("Currency already exists")

var cache = make(map[int64]model.Currency)

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

func (s *CurrencyStore) FindByCode(code string) (*model.Currency, error) {
	row := s.db.QueryRow("SELECT * FROM Currencies WHERE code = ?;", code)

	var currency model.Currency

	err := row.Scan(&currency.Id, &currency.Code, &currency.FullName, &currency.Sign)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, CurrencyNotFoundError
	}

	if err != nil {
		slog.Error("Unable to map row to model", "error", err)
		return nil, err
	}

	return &currency, nil
}

func (s *CurrencyStore) FindById(id int64) (*model.Currency, error) {
	if currency, exists := cache[id]; exists {
		return &currency, nil
	}

	row := s.db.QueryRow("SELECT * FROM Currencies WHERE id = ?;", id)

	var currency model.Currency

	err := row.Scan(&currency.Id, &currency.Code, &currency.FullName, &currency.Sign)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, CurrencyNotFoundError
	}

	if err != nil {
		slog.Error("Unable to map row to model", "error", err)
		return nil, err
	}

	cache[id] = currency

	return &currency, nil
}

func (s *CurrencyStore) Save(name string, code string, sign string) (*model.Currency, error) {
	row := s.db.QueryRow(
		"INSERT INTO Currencies (full_name, code, sign) VALUES (?, ?, ?) RETURNING id, code, full_name, sign;",
		name, code, sign)

	var currency model.Currency

	err := row.Scan(&currency.Id, &currency.Code, &currency.FullName, &currency.Sign)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return nil, CurrencyAlreadyExistsError
	}

	if err != nil {
		slog.Error("Unable to map row to model", "error", err)
		return nil, err
	}

	return &currency, nil
}
