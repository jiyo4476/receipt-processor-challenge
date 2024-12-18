package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func tryValidateShortDescription(t *testing.T, input string, isValid bool) {
	validate := validator.New()
	validate.RegisterValidation("valid_func", CorrectShortDescription)
	err := validate.Var(input, "valid_func")
	if isValid {
		assert.NoError(t, err)
		return
	}
	assert.Error(t, err)
}

func TestCorrectShortDescription_ValidSingleWord(t *testing.T) {
	tryValidateShortDescription(t, "Valid", true)
}

func TestCorrectShortDescription_ValidMultiWord(t *testing.T) {
	tryValidateShortDescription(t, "Valid Description", true)
	tryValidateShortDescription(t, "This is a valid description", true)
}

func TestCorrectShortDescription_ValidHyphenatedDescription(t *testing.T) {
	tryValidateShortDescription(t, "This-is-a-valid-description", true)
}

func TestCorrectShortDescription_ValidLength(t *testing.T) {
	tryValidateShortDescription(t, "This is a valid short description that is exactly 50", true)
}

func TestCorrectShortDescription_Invalid(t *testing.T) {
	tryValidateShortDescription(t, "Hello@world", false)
}

func TestCorrectShortDescription_InvalidType(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("valid_func", CorrectShortDescription)
	err := validate.Var(123, "valid_func")
	assert.Error(t, err)
}

func tryValidateCashValue(t *testing.T, input string, isValid bool) {
	validate := validator.New()
	validate.RegisterValidation("correctCashValue", CorrectCashValue)

	type TestStruct struct {
		CashValue string `validate:"required,min=4,correctCashValue"`
	}

	valid := TestStruct{CashValue: input}

	err := validate.Struct(valid)
	if isValid {
		assert.NoError(t, err)
		return
	}
	assert.Error(t, err)
}

func TestCorrectCashValue_ValidCashValue_Valid(t *testing.T) {
	tryValidateCashValue(t, "00.00", true)
	tryValidateCashValue(t, "99.99", true)
	tryValidateCashValue(t, "999999.99", true)
}

func TestCorrectCashValue_ValidCashValue_Len(t *testing.T) {
	tryValidateCashValue(t, "1234.56", true)
	tryValidateCashValue(t, "12345.67", true)
	tryValidateCashValue(t, "12345.6755", false)
}

func TestCorrectCashValue_InvalidChar(t *testing.T) {
	tryValidateCashValue(t, "abc", false)
	tryValidateCashValue(t, "%.^?", false)
	tryValidateCashValue(t, " ", false)
}

func TestCorrectCashValue_InvalidType(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("correctCashValue", CorrectCashValue)

	type TestStruct struct {
		CashValue int64 `validate:"required,min=4,correctCashValue"`
	}

	instance := TestStruct{CashValue: int64(1234)}

	err := validate.Struct(instance)
	assert.Error(t, err)
}

func tryValidateRetailerName(t *testing.T, input string, isValid bool) {
	validate := validator.New()
	validate.RegisterValidation("correctRetailerName", CorrectRetailerName)

	type TestStruct struct {
		RetailerName string `validate:"required,min=1,correctRetailerName"`
	}

	valid := TestStruct{RetailerName: input}

	err := validate.Struct(valid)
	if isValid {
		assert.NoError(t, err)
		return
	}
	assert.Error(t, err)
}

func TestCorrectRetailerName_InvalidType(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("correctRetailerName", CorrectRetailerName)

	type TestStruct struct {
		RetailerName int `validate:"required,min=1,correctRetailerName"`
	}

	instance := TestStruct{RetailerName: 1234}

	err := validate.Struct(instance)
	assert.Error(t, err)
}

func TestCorrectRetailerName_ValidRetailerName(t *testing.T) {
	tryValidateRetailerName(t, "Valid Retailer Name", true)
}

func TestCorrectRetailerName_InvalidRetailerName_Non_alphanumeric(t *testing.T) {
	tryValidateRetailerName(t, "Valid Retailer Name & Co.", false)
}

func tryValidateCorrectDate(t *testing.T, input string, isValid bool) {
	validate := validator.New()
	validate.RegisterValidation("correctDate", CorrectDate)

	type TestStruct struct {
		PurchaseDate string `validate:"required,len=10,correctDate" time_format:"2022-01-01"`
	}

	valid := TestStruct{PurchaseDate: input}

	err := validate.Struct(valid)
	if isValid {
		assert.NoError(t, err)
		return
	}
	assert.Error(t, err)
}

func TestValidateCorrectDateInvalid(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("correctDate", CorrectDate)

	type TestStruct struct {
		PurchaseDate int `validate:"required,len=10,correctDate" time_format:"2022-01-01"`
	}

	valid := TestStruct{PurchaseDate: 10}

	err := validate.Struct(valid)
	assert.Error(t, err)
}

func TestCorrectDateInvalid_Type(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("correctDate", CorrectDate)

	type TestStruct struct {
		PurchaseDate int `validate:"required,len=10,correctDate" time_format:"2022-01-01"`
	}

	valid := TestStruct{PurchaseDate: 1234567899}

	err := validate.Struct(valid)
	assert.Error(t, err)
}

func TestCorrectDateValid_CharRange(t *testing.T) {
	tryValidateCorrectDate(t, "2022-01-01", true)
	tryValidateCorrectDate(t, "9999-12-01", true)
	tryValidateCorrectDate(t, "2022-12-01", true)
	tryValidateCorrectDate(t, "2022-12-31", true)
	tryValidateCorrectDate(t, "9999-00-01", false)
	tryValidateCorrectDate(t, "9999-13-01", false)
	tryValidateCorrectDate(t, "9999-12-00", false)
	tryValidateCorrectDate(t, "9999-12-32", false)
}

func TestCorrectDateInvalid_SectionLen(t *testing.T) {
	tryValidateCorrectDate(t, "2024-01-01", true)
	tryValidateCorrectDate(t, "20224-01-01", false)
	tryValidateCorrectDate(t, "2024-011-01", false)
	tryValidateCorrectDate(t, "2024-01-011", false)
}

func TestValidateCorrectTime_InvalidType(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("correctTime", CorrectTime)

	type TestStruct struct {
		PurchaseTime int64 `validate:"required,correctTime,len=5" time_format:"13:01"`
	}

	valid := TestStruct{PurchaseTime: int64(123)}

	err := validate.Struct(valid)
	assert.Error(t, err)
}

func tryValidateCorrectTime(t *testing.T, input string, isValid bool) {
	validate := validator.New()
	validate.RegisterValidation("correctTime", CorrectTime)

	type TestStruct struct {
		PurchaseTime string `validate:"required,correctTime,len=5" time_format:"13:01"`
	}

	valid := TestStruct{PurchaseTime: input}

	err := validate.Struct(valid)
	if isValid {
		assert.NoError(t, err)
		return
	}
	assert.Error(t, err)
}

func TestCorrectTime_CharRange(t *testing.T) {
	tryValidateCorrectTime(t, "00:00", true)
	tryValidateCorrectTime(t, "24:00", true)
	tryValidateCorrectTime(t, "23:59", true)
	tryValidateCorrectTime(t, "25:00", false)
	tryValidateCorrectTime(t, "23:60", false)
	tryValidateCorrectTime(t, "30:00", false)
}

func TestCorrectTime_SectionLen(t *testing.T) {
	tryValidateCorrectTime(t, "00:00", true)
	tryValidateCorrectTime(t, "001:00", false)
	tryValidateCorrectTime(t, "00:001", false)
}
