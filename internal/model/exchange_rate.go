package model

type ExchangeRate struct {
	Id               int64
	BaseCurrencyId   int64
	TargetCurrencyId int64
	Rate             float64
}
