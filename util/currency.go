package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	//check if the currency is supported
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
