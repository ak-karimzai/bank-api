package util

import "log"

const (
	USD = "USD"
	EUR = "EUR"
	RUB = "RUB"
)

func IsSupportedCurrency(currency string) bool {
	log.Println("fuck")
	switch currency {
	case USD, EUR, RUB:
		return true
	}
	return false
}
