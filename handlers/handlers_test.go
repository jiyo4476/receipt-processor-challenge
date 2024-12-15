package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/jiyo4476/receipt-processor-challenge/models"
	"github.com/jiyo4476/receipt-processor-challenge/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func makeRequest(method string, url string, body interface{}) (*httptest.ResponseRecorder, error) {
	test_router := router.SetUpRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, url, nil)
	if body != nil {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("error encoding JSON: %v", err)
		}
		c.Request = httptest.NewRequest(method, "http://localhost:8080"+url, buf)
		c.Request.Header.Set("Content-Type", "application/json")
	}
	test_router.ServeHTTP(w, c.Request)
	return w, nil
}

func attemptProcessReceipt(t *testing.T, receipt models.Receipt) (string, error) {
	w, err := makeRequest("POST", "/receipts/process", receipt)
	assert.NoError(t, err, "error making request")
	if err != nil {
		t.Fatalf("error making POST request")
		return "", fmt.Errorf("error making POST request: %w", err)
	}

	var response struct {
		ID string `json:"id"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error in unmarshaling JSON response")
	if err != nil {
		t.Fatalf("error unmarshaling JSON response")
		return "", fmt.Errorf("error unmarshaling JSON response: %w", err)
	}

	assert.NotEmpty(t, response.ID, "Response should contain a non-empty ID")
	return response.ID, nil
}

func attemptGetPoints(t *testing.T, id string) (int64, error) {
	endpoint := fmt.Sprintf("/receipts/%s/points", id)
	w, err := makeRequest("GET", endpoint, nil)
	assert.NoError(t, err, "error making request")
	if err != nil {
		t.Fatalf("error making GET request")
		return -1, fmt.Errorf("error making GET request: %w", err)
	}

	var pointsResponse struct {
		Points int64 `json:"points"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &pointsResponse)
	assert.NoError(t, err, "Error in unmarshaling JSON response")
	if err != nil {
		t.Fatalf("error unmarshaling JSON response")
		return -1, fmt.Errorf("error unmarshaling JSON response: %w", err)
	}
	return int64(pointsResponse.Points), nil
}

func TestGetReceiptsPoints_NotFound(t *testing.T) {
	// Try to fetch points with an invalid ID
	endpoint := "/receipts/adb6b560-0eef-42bc-9d16-df48f30e89b2/points"
	res, err := makeRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatalf("Error making GET request: %v", err)
	}

	assert.Equal(t, http.StatusNotFound, res.Code, "Expected status code 404 for Not Found")
}

// Test Invalid
func TestGetReceiptsPoints_InvalidUUID_Hyphen(t *testing.T) {
	endpoint := fmt.Sprintf("/receipts/%s/points", "0000000--0000-0000-0000-000000000000")
	w, err := makeRequest("GET", endpoint, nil)
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected 404 status code for id")
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
}

func TestGetReceiptsPoints_InvalidUUID_Whitespace(t *testing.T) {
	endpoint := fmt.Sprintf("/receipts/%s/points", "0000000%20-0000-0000-0000-000000000000")
	w, err := makeRequest("GET", endpoint, nil)
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected 404 status code for id")
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
}

func TestProcessReceipt_ValidReceipt(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-12-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		},
		Total: "18.74",
	}

	id, err := attemptProcessReceipt(t, receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	assert.NotEmpty(t, id, "Response should contain a non-empty ID")

	matched, err := regexp.MatchString("^\\S+$", id)
	if err != nil {
		t.Fatalf("Error matching ID: %v", err)
	}
	assert.True(t, matched, "value id does not match spec %s", id)
}

func TestProcessReceipt_InvalidReceipt_Date(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01-01", // Invalid date format
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		},
		Total: "18.74",
	}

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")
}

func TestProcessReceipt_InvalidReceipt_Time(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "abcd",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		},
		Total: "18.74",
	}

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")
}

func TestProcessReceipt_InvalidReceipt_Total(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "", // Missing retailer
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items:        []models.Item{},
		Total:        "abc", // Invalid total
	}

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 status code for invalid receipt")
}

func TestProcessReceipt_InvalidReceipt_ItemTotal(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-12-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		},
		Total: "20.74",
	}

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 status code for invalid receipt")
}

func TestProcessReceipt_Invalid_Body(t *testing.T) {
	w, err := makeRequest("POST", "/receipts/process", nil)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 status code for invalid receipt")
}

func TestProcessReceipt_InvalidReceipt_No_DateAndTime(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "", // Missing PurchaseDate
		PurchaseTime: "", // Missing PurchaseTime
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "01.64",
	}

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 status code for invalid receipt")
}

func TestProcessReceipt_InvalidTotal(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "01.64.00", // Invalid Total
	}

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 status code for invalid receipt")
}

func TestProcessReceipt_Invalid_Item_Price(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49.00"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "01.64", // Invalid Total
	}

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 status code for invalid receipt")
}

func TestGetReceiptsPoints_ValidID_01(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}

	receiptID, err := attemptProcessReceipt(t, receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	assert.NotEmpty(t, receiptID, "Response should contain a non-empty ID")

	points, err := attemptGetPoints(t, receiptID)
	if err != nil {
		t.Fatalf("Error getting points: %v", err)
	}

	// Check points
	assert.Equal(t, int64(28), int64(points), "Expected 28 points for this receipt")
}

func TestGetReceiptsPoints_ValidID_02(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total: "9.00",
	}

	receiptID, err := attemptProcessReceipt(t, receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	points, err := attemptGetPoints(t, receiptID)
	if err != nil {
		t.Fatalf("Error getting points: %v", err)
	}

	assert.Equal(t, int64(109), points, "Expected 109 points for this receipt")
}

func TestGetReceiptsPoints_ValidID_03(t *testing.T) {
	// First, process a receipt to get a valid ID
	receipt := models.Receipt{
		Retailer:     "Walgreens",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "08:13",
		Items: []models.Item{
			{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
			{ShortDescription: "Dasani", Price: "1.40"},
		},
		Total: "2.65",
	}

	receiptID, err := attemptProcessReceipt(t, receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	points, err := attemptGetPoints(t, receiptID)
	if err != nil {
		t.Fatalf("Error getting points: %v", err)
	}

	assert.Equal(t, int64(15), points, "Expected 15 points for this receipt")
}

func TestGetReceiptsPoints_ValidID_04(t *testing.T) {
	// First, process a receipt to get a valid ID
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Items: []models.Item{
			{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
		},
		Total: "1.25",
	}

	receiptID, err := attemptProcessReceipt(t, receipt)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	points, err := attemptGetPoints(t, receiptID)
	if err != nil {
		t.Fatalf("Error getting points: %v", err)
	}
	assert.Equal(t, int64(31), points, "Expected 31 points for this receipt")
}
