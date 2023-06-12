package util

var supportedCurrencies = []string{"USD", "EUR", "RUB"}

func IsSupportedCurrency(currency string) bool {
	for _, pkgCurrency := range supportedCurrencies {
		if currency == pkgCurrency {
			return true
		}
	}
	return false
}
