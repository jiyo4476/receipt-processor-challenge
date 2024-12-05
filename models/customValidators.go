package models

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var CorrectShortDescription validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched, _ := regexp.MatchString("^[\\w\\s\\-]+$", value)
		return matched
	}
	return false
}

var CorrectCashValue validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched, _ := regexp.MatchString("^\\d+\\.\\d{2}$", value)
		return matched
	}
	return false
}

var CorrectRetailerName validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched, _ := regexp.MatchString("^[\\w\\s\\-&]+$", value)
		return matched
	}
	return false
}

var CorrectDate validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched, _ := regexp.MatchString("^(\\d{4})-(1[0-2]|0[1-9])-(3[01]|[1-2]\\d|0[1-9])$", value)
		return matched
	}
	return false
}

var CorrectTime validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched, _ := regexp.MatchString("(^24:00$)|(^([01][0-9]|[2][0-3]):[0-5][0-9]$)", value)
		return matched
	}
	return false
}
