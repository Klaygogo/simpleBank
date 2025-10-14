package api

import (
	"strings"

	"github.com/Klaygogo/simplebank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLever validator.FieldLevel) bool {
	if currency, ok := fieldLever.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
