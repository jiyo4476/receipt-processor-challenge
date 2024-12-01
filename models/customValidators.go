package models

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var CorrectShortDescription validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched, _ := regexp.MatchString("^[\\w\\s\\-]+$", value)
		if matched {
			return true
		}
	}
	return false
}

var CorrectCashValue validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched, _ := regexp.MatchString("^\\d+\\.\\d{2}$", value)
		valueAsFloat, _ := strconv.ParseFloat(value, 64)
		if matched && valueAsFloat > 0 {
			return true
		}
	}
	return false
}

var CorrectRetailerName validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched, _ := regexp.MatchString("^[\\w\\s\\-&]+$", value)
		if matched && len(value) > 0 {
			return true
		}
	}
	return false
}
