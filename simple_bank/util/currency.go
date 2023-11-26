package util

func getCurrencies() []string {
	return []string{"USD", "EUR", "YEN", "CAD", "AUD"}
}

func GetCurrency(index int) string {
	return getCurrencies()[index]
}
func IsSupportedCurrency(currency string) bool {
	for _, c := range getCurrencies() {
		if c == currency {
			return true
		}
	}
	return false
}
