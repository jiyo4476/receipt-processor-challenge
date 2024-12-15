package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNumAlphanumerical(t *testing.T) {
	value := getNumAlphanumerical("")
	assert.Equal(t, int64(0), value)
	value = getNumAlphanumerical("a")
	assert.Equal(t, int64(1), value)
	value = getNumAlphanumerical("hello123world")
	assert.Equal(t, int64(13), value)
	value = getNumAlphanumerical("ABCDEF12345")
	assert.Equal(t, int64(11), value)
	value = getNumAlphanumerical("!@#$%^&*")
	assert.Equal(t, int64(0), value)
	value = getNumAlphanumerical(" hello world ")
	assert.Equal(t, int64(10), value)
	value = getNumAlphanumerical("h3110,w0r1d!")
	assert.Equal(t, int64(10), value)

}

func TestGetPointsRoundAmount(t *testing.T) {
	value := getPointsRoundAmount("00.00")
	assert.Equal(t, int64(50), value)
	value = getPointsRoundAmount("00.59")
	assert.Equal(t, int64(0), value)
}

func TestGetPointsMultipleOf25(t *testing.T) {
	value := getPointsMultipleOf25("10.00")
	assert.Equal(t, int64(25), value)
	value = getPointsMultipleOf25("00.25")
	assert.Equal(t, int64(25), value)
	value = getPointsMultipleOf25("00.50")
	assert.Equal(t, int64(25), value)
	value = getPointsMultipleOf25("00.75")
	assert.Equal(t, int64(25), value)
	value = getPointsMultipleOf25("00.99")
	assert.Equal(t, int64(0), value)
}

func TestGetPointsForOddDate(t *testing.T) {
	value := getPointsForOddDate("2022-01-01")
	assert.Equal(t, int64(6), value)
	value = getPointsForOddDate("2022-01-02")
	assert.Equal(t, int64(0), value)
}

func TestGetPointsForTimeOfPurchase(t *testing.T) {
	value := getPointsForTimeOfPurchase("14:00")
	assert.Equal(t, int64(10), value)
	value = getPointsForTimeOfPurchase("15:00")
	assert.Equal(t, int64(10), value)
	value = getPointsForTimeOfPurchase("13:00")
	assert.Equal(t, int64(0), value)
	value = getPointsForTimeOfPurchase("16:00")
	assert.Equal(t, int64(0), value)
}

func TestGetPointsForItems(t *testing.T) {
	items := []Item{
		{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
		{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
		{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
		{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
	}
	value, err := getPointsForItems(items)
	if err != nil {
		t.Fatalf("Error calculating points for items: %v", err)
	}
	assert.Equal(t, int64(6), value)
}

func TestReceiptPoints_SuccessPath(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "01.64",
	}
	points, err := receipt.Points()
	if err != nil {
		t.Fatalf("Error calculating points for receipt: %v", err)
	}
	assert.Equal(t, int64(28), points)
}

func TestReceiptPointsForTotal25(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "14:00",
		Items: []Item{
			{ShortDescription: "Dew 12PK", Price: "25.00"},
		},
		Total: "25.00",
	}
	points, err := receipt.Points()
	if err != nil {
		t.Fatalf("Error calculating points for receipt: %v", err)
	}
	assert.Equal(t, int64(103), points)
}

func TestReceiptPointsWithFourItems(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}
	points, err := receipt.Points()
	if err != nil {
		t.Fatalf("Error calculating points for receipt: %v", err)
	}
	assert.Equal(t, int64(28), points)
}

func TestPointsForOddPurchaseDate(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "10.00"},
		},
		Total: "10.00",
	}
	points, err := receipt.Points()
	if err != nil {
		t.Fatalf("Error calculating points for receipt: %v", err)
	}
	assert.Equal(t, int64(95), points)
}

func TestValidateReceiptTotal(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item1", Price: "10.00"},
			{ShortDescription: "Item2", Price: "10.00"},
		},
		Total: "20.00",
	}
	err := receipt.ValidateTotal()
	assert.Nil(t, err, "Receipt should be valid")
}

func TestValidateInvalidReceiptTotal(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "10.00"},
		},
		Total: "20.00",
	}
	err := receipt.ValidateTotal()
	assert.NotNil(t, err, "Receipt should not be valid")
}

func TestValidateInvalidTotal(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "10.00"},
		},
		Total: "20.00.00",
	}
	err := receipt.ValidateTotal()
	assert.NotNil(t, err, "Receipt should not be valid")
}

func TestValidateInvalidItemPrice(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "10.00.00"},
		},
		Total: "20.00",
	}
	err := receipt.ValidateTotal()
	assert.NotNil(t, err, "Receipt should not be valid")
}

func TestPointsForInvalidItems(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "10.00.12"},
		},
		Total: "20.00",
	}
	points, err := receipt.Points()
	assert.NotNil(t, err, "Error should be returned")
	assert.Equal(t, int64(-1), points)
}
