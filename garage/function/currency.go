package function

import "github.com/biter777/countries"

func GetCurrencyInfo(countryCode string) (name string, code string, err error) {
	country := countries.ByName(countryCode)
	currency := country.Currency()
	cCode := currency.Alpha()

	return currency.String(), cCode, nil
}
