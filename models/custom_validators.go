package models

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var correctShortDescriptionFormat = regexp.MustCompile(`^[\w\s\-]+$`)
var correctCashValueFormat = regexp.MustCompile(`^\d+\.\d{2}$`)
var retailerNameFormat = regexp.MustCompile(`^[\w\s\-&]+$`)
var correctDateFormat = regexp.MustCompile(`^(\d{4})-(1[0-2]|0[1-9])-(3[01]|[1-2]\d|0[1-9])$`)
var correctTimeFormat = regexp.MustCompile(`(^24:00$)|(^([01][0-9]|[2][0-3])):[0-5][0-9]$`)

var CorrectShortDescription validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched := correctShortDescriptionFormat.MatchString(value)
		return matched
	}
	return false
}

var CorrectCashValue validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched := correctCashValueFormat.MatchString(value)
		return matched
	}
	return false
}

var CorrectRetailerName validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched := retailerNameFormat.MatchString(value)
		return matched
	}
	return false
}

var CorrectDate validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched := correctDateFormat.MatchString(value)
		return matched
	}
	return false
}

var CorrectTime validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		matched := correctTimeFormat.MatchString(value)
		return matched
	}
	return false
}
