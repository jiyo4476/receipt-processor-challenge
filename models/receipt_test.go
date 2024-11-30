package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNumAlphanumerical(t *testing.T) {
	value := getNumAlphanumerical("")
	assert.Equal(t, 0, value)
	value = getNumAlphanumerical("a")
	assert.Equal(t, 1, value)
	value = getNumAlphanumerical("hello123world")
	assert.Equal(t, 13, value)
	value = getNumAlphanumerical("ABCDEF12345")
	assert.Equal(t, 11, value)
	value = getNumAlphanumerical("!@#$%^&*")
	assert.Equal(t, 0, value)
	value = getNumAlphanumerical(" hello world ")
	assert.Equal(t, 10, value)
	value = getNumAlphanumerical("h3110,w0r1d!")
	assert.Equal(t, 10, value)

}

func TestGetPointsRoundAmount(t *testing.T) {
	value := getPointsRoundAmount(".00")
	assert.Equal(t, 50, value)
	value = getPointsRoundAmount(".59")
	assert.Equal(t, 0, value)
}

func TestGetPointsMultipleOf25(t *testing.T) {
	value := getPointsMultipleOf25(".00")
	assert.Equal(t, 25, value)
	value = getPointsMultipleOf25(".25")
	assert.Equal(t, 25, value)
	value = getPointsMultipleOf25(".50")
	assert.Equal(t, 25, value)
	value = getPointsMultipleOf25(".75")
	assert.Equal(t, 25, value)
	value = getPointsMultipleOf25(".99")
	assert.Equal(t, 0, value)
}

func TestGetPointsForOddDate(t *testing.T) {
	value := getPointsForOddDate(1)
	assert.Equal(t, 6, value)
	value = getPointsForOddDate(2)
	assert.Equal(t, 0, value)
}

func TestGetPointsForTimeOfPurchase(t *testing.T) {
	value := getPointsForTimeOfPurchase("14")
	assert.Equal(t, 10, value)
	value = getPointsForTimeOfPurchase("15")
	assert.Equal(t, 10, value)
	value = getPointsForTimeOfPurchase("13")
	assert.Equal(t, 0, value)
	value = getPointsForTimeOfPurchase("16")
	assert.Equal(t, 0, value)
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
	assert.Equal(t, 6, value)
}

func TestReceiptPoints_SuccessPath(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "14:30",
		Items: []Item{
			{ShortDescription: "Item 1", Price: "10.00"},
			{ShortDescription: "Item 2", Price: "20.00"},
			{ShortDescription: "Item 3", Price: "30.00"},
			{ShortDescription: "Item 4", Price: "40.00"},
		},
		Total: "100.00",
	}
	points := receipt.Points()
	assert.Equal(t, 100, points)
}

func TestReceiptPointsForTotal25(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "14:00",
		Items:        []Item{},
		Total:        "25.00",
	}
	points := receipt.Points()
	assert.Equal(t, 81, points)
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
	points := receipt.Points()
	assert.Equal(t, 28, points)
}

func TestPointsForOddPurchaseDate(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items:        []Item{},
		Total:        "10.00",
	}
	points := receipt.Points()
	assert.Equal(t, 56, points)
}
