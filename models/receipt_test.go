package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createRetailerTestReceipt(name string) Receipt {
	receipt := Receipt{
		Retailer:     name,
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "10.00"},
			{ShortDescription: "Item02", Price: "10.00"},
		},
		Total: "20.00",
	}
	return receipt
}

func TestGetNumAlphanumerical_NoChars(t *testing.T) {
	receipt := createRetailerTestReceipt("")
	value := receipt.getPointsAlphanumerical()
	assert.Equal(t, int64(0), value)
}

func TestGetNumAlphanumerical_OneChar(t *testing.T) {
	receipt := createRetailerTestReceipt("a")
	value := receipt.getPointsAlphanumerical()
	assert.Equal(t, int64(1), value)
}
func TestGetNumAlphanumerical_Valid01(t *testing.T) {
	receipt := createRetailerTestReceipt("hello123world")
	value := receipt.getPointsAlphanumerical()
	assert.Equal(t, int64(13), value)
}
func TestGetNumAlphanumerical_Valid02(t *testing.T) {
	receipt := createRetailerTestReceipt("ABCDEF12345")
	value := receipt.getPointsAlphanumerical()
	assert.Equal(t, int64(11), value)
}
func TestGetNumAlphanumerical_Valid03(t *testing.T) {
	receipt := createRetailerTestReceipt(" hello world ")
	value := receipt.getPointsAlphanumerical()
	assert.Equal(t, int64(10), value)
}
func TestGetNumAlphanumerical_Valid04(t *testing.T) {
	receipt := createRetailerTestReceipt("h3110,w0r1d!")
	value := receipt.getPointsAlphanumerical()
	assert.Equal(t, int64(10), value)
}
func TestGetNumAlphanumerical_Invalid(t *testing.T) {
	receipt := createRetailerTestReceipt("!@#$%^&*")
	value := receipt.getPointsAlphanumerical()
	assert.Equal(t, int64(0), value)
}

func TestGetPointsRoundAmount_Valid(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items:        []Item{},
		Total:        "00.00",
	}
	value := receipt.getPointsRoundAmount()
	assert.Equal(t, int64(50), value)
}

func TestGetPointsRoundAmount_NoPoints(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.50"},
			{ShortDescription: "Item02", Price: "0.09"},
		},
		Total: "00.59",
	}
	value := receipt.getPointsRoundAmount()
	assert.Equal(t, int64(0), value)
}

func TestGetPointsMultipleOf25_00(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "5.00"},
			{ShortDescription: "Item02", Price: "5.00"},
		},
		Total: "10.00",
	}
	value := receipt.getPointsMultipleOf25()
	assert.Equal(t, int64(25), value)
}

func TestGetPointsMultipleOf25_25(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "5.00"},
			{ShortDescription: "Item02", Price: "5.25"},
		},
		Total: "10.25",
	}
	value := receipt.getPointsMultipleOf25()
	assert.Equal(t, int64(25), value)
}

func TestGetPointsMultipleOf25_50(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.25"},
			{ShortDescription: "Item02", Price: "0.25"},
		},
		Total: "00.50",
	}
	value := receipt.getPointsMultipleOf25()
	assert.Equal(t, int64(25), value)
}

func TestGetPointsMultipleOf25_75(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.50"},
			{ShortDescription: "Item02", Price: "0.25"},
		},
		Total: "00.75",
	}
	value := receipt.getPointsMultipleOf25()
	assert.Equal(t, int64(25), value)
}

func TestGetPointsMultipleOf25_NoPoints(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.50"},
			{ShortDescription: "Item02", Price: "0.49"},
		},
		Total: "00.99",
	}
	value := receipt.getPointsMultipleOf25()
	assert.Equal(t, int64(0), value)
}

func TestGetPointsForOddDate_Odd(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.50"},
			{ShortDescription: "Item02", Price: "0.49"},
		},
		Total: "00.99",
	}
	value := receipt.getPointsForOddDate()
	assert.Equal(t, int64(6), value)
}

func TestGetPointsForOddDate_Even(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.50"},
			{ShortDescription: "Item02", Price: "0.49"},
		},
		Total: "00.99",
	}
	value := receipt.getPointsForOddDate()
	assert.Equal(t, int64(0), value)
}

func TestGetPointsForTimeOfPurchase_1PM(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:00",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.50"},
			{ShortDescription: "Item02", Price: "0.49"},
		},
		Total: "00.99",
	}
	value := receipt.getPointsForTimeOfPurchase()
	assert.Equal(t, int64(0), value)
}

func TestGetPointsForTimeOfPurchase_2PM(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "14:00",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.50"},
			{ShortDescription: "Item02", Price: "0.49"},
		},
		Total: "00.99",
	}
	value := receipt.getPointsForTimeOfPurchase()
	assert.Equal(t, int64(10), value)
}
func TestGetPointsForTimeOfPurchase_3PM(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "15:00",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.50"},
			{ShortDescription: "Item02", Price: "0.49"},
		},
		Total: "00.99",
	}
	value := receipt.getPointsForTimeOfPurchase()
	assert.Equal(t, int64(10), value)
}
func TestGetPointsForTimeOfPurchase_4PM(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "16:00",
		Items: []Item{
			{ShortDescription: "Item01", Price: "0.50"},
			{ShortDescription: "Item02", Price: "0.49"},
		},
		Total: "00.99",
	}
	value := receipt.getPointsForTimeOfPurchase()
	assert.Equal(t, int64(0), value)
}

func TestGetPointsForItems(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
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

	value, err := receipt.getPointsForItems()
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
