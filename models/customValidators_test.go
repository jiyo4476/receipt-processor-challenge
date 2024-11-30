package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestCorrectShortDescription_Valid(t *testing.T) {
	//var validate *validator.Validate
	var validate = validator.New()
	validate.RegisterValidation("correctShortDescription", CorrectShortDescription)
	input := "Valid Description 123"
	err := validate.Var(input, "correctShortDescription")
	assert.NoError(t, err)
}

func TestCorrectShortDescription_ValidShortDescription(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("correctShortDescription", CorrectShortDescription)

	input := "This is a valid description"
	err := validate.Var(input, "correctShortDescription")

	assert.NoError(t, err)
}

func TestCorrectShortDescription_ValidHyphenatedDescription(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("correctShortDescription", CorrectShortDescription)

	input := "This-is-a-valid-description"
	err := validate.Var(input, "correctShortDescription")

	assert.NoError(t, err)
}

func TestCorrectShortDescription_ValidSingleWord(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("correctShortDescription", CorrectShortDescription)

	input := "Apple"
	err := validate.Var(input, "correctShortDescription")

	assert.NoError(t, err)
}

func TestCorrectShortDescription_ValidShortDescription_02(t *testing.T) {
	var validate = validator.New()
	validate.RegisterValidation("correctShortDescription", CorrectShortDescription)

	input := "This-is-a-valid-description"
	err := validate.Var(input, "correctShortDescription")

	assert.NoError(t, err)
}

func TestCorrectShortDescription_ValidLength(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("correctShortDescription", CorrectShortDescription)
	input := "This is a valid short description that is exactly 50"
	err := validate.Var(input, "correctShortDescription")
	assert.NoError(t, err)
}

func TestCorrectCashValue_ValidCashValue_Ten(t *testing.T) {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("correct_cash_value", CorrectCashValue)

	type TestStruct struct {
		CashValue string `binding:"required,correct_cash_value,min=4"`
	}

	testData := TestStruct{CashValue: "10.00"}

	err := validate.Struct(testData)
	assert.NoError(t, err)
}

func TestCorrectCashValue_ValidCashValue_Change(t *testing.T) {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("correct_cash_value", CorrectCashValue)

	type TestStruct struct {
		CashValue string `binding:"required,correct_cash_value,min=4"`
	}

	testData := TestStruct{CashValue: "0.99"}
	err := validate.Struct(testData)
	assert.NoError(t, err)
}

func TestCorrectCashValue_ValidCashValue_Len_06(t *testing.T) {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("correctCashValue", CorrectCashValue)

	type TestStruct struct {
		CashValue string `binding:"required,correctCashValue,min=4"`
	}

	valid := TestStruct{CashValue: "123.45"}

	err := validate.Struct(valid)
	assert.NoError(t, err)
}

func TestCorrectCashValue_ValidCashValue(t *testing.T) {
	var validate = validator.New()
	validate.RegisterValidation("correct_cash_value", CorrectCashValue)

	type TestStruct struct {
		CashValue string `binding:"required,correct_cash_value,min=4"`
	}

	valid := TestStruct{CashValue: "9999.00"}

	err := validate.Struct(valid)

	assert.NoError(t, err)
}

func TestCorrectCashValue_ValidCashValue_02(t *testing.T) {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("correct_cash_value", CorrectCashValue)

	type TestStruct struct {
		CashValue string `binding:"required,correct_cash_value,min=4"`
	}

	valid := TestStruct{CashValue: "50.50"}

	err := validate.Struct(valid)
	assert.NoError(t, err)
}

func TestCorrectRetailerName_ValidRetailerName(t *testing.T) {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("correctRetailerName", CorrectRetailerName)

	type TestStruct struct {
		RetailerName string `binding:"required,correctRetailerName,min=1"`
	}

	valid := TestStruct{RetailerName: "Valid Retailer Name"}

	err := validate.Struct(valid)
	assert.NoError(t, err)
}

func TestCorrectRetailerName_ValidRetailerName_Non_alphanumeric(t *testing.T) {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("correctRetailerName", CorrectRetailerName)

	type TestStruct struct {
		RetailerName string `binding:"required,correctRetailerName,min=1"`
	}

	valid := TestStruct{RetailerName: "Target~"}

	err := validate.Struct(valid)
	assert.NoError(t, err)
}

func TestCorrectRetailerName_InvalidRetailerName(t *testing.T) {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("correctRetailerName", CorrectRetailerName)

	type TestStruct struct {
		RetailerName string `binding:"required,correctRetailerName,min=1"`
	}

	valid := TestStruct{RetailerName: ""}

	err := validate.Struct(valid)
	assert.Error(t, err)
}
