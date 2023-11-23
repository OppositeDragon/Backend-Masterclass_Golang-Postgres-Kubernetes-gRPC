package util

func getCurrencies() []string {
	return []string{"USD", "EUR", "CAD", "YEN", "AUD"}
}

func IsSupportedCurrency(currency string) bool {
	for _, c := range getCurrencies() {
		if c == currency {
			return true
		}
	}
	return false
}
