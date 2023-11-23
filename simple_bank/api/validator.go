package api

import (
	"github.com/go-playground/validator/v10"
	"simple_bank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, value := fl.Field().Interface().(string); value {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
