package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"receipt-processor-challenge/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

// Integration Tests

func SetUpRouter() *gin.Engine {
	validate = validator.New(validator.WithRequiredStructEnabled())
	router := gin.Default()

	// Register custom validation functions for the test router
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("correctRetailerName", models.CorrectRetailerName)
		v.RegisterValidation("correctShortDescription", models.CorrectShortDescription)
		v.RegisterValidation("correctCashValue", models.CorrectCashValue)
	}

	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getReceiptsPoints)
	return router
}

func makeRequest(method string, url string, body interface{}) (*httptest.ResponseRecorder, error) {
	router := SetUpRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, url, nil)
	if body != nil {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("error encoding JSON: %w", err)
		}
		c.Request = httptest.NewRequest(method, "http://127.0.0.1:8080"+url, buf)
		c.Request.Header.Set("Content-Type", "application/json")
	}
	router.ServeHTTP(w, c.Request)
	return w, nil
}

func TestProcessReceipt_ValidReceipt(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
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

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON response: %v", err)
	}
	assert.NotEmpty(t, response.ID, "Response should contain a non-empty ID")
}

func TestProcessReceipt_InvalidReceipt(t *testing.T) {
	// Example of an invalid receipt (missing required fields or incorrect format)
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

func TestProcessReceipt_InvalidReceipt_02(t *testing.T) {
	// Example of an invalid receipt (missing required fields or incorrect format)
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
	// Example of an invalid receipt (missing required fields or incorrect format)
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01", // Missing PurchaseDate
		PurchaseTime: "13:01",      // Missing PurchaseTime
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

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error processing receipt: %v", err)
	}
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 for POST")

	var postResponse struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &postResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling POST response: %v", err)
	}
	receiptID := postResponse.ID
	assert.NotEmpty(t, receiptID, "POST response should contain a non-empty ID")

	// Now, fetch points using the valid ID
	endpoint := fmt.Sprintf("/receipts/%s/points", receiptID)
	res, err := makeRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatalf("Error making GET request: %v", err)
	}

	assert.Equal(t, http.StatusOK, res.Code, "Expected status code 200 for GET")

	var pointsResponse struct {
		Points int `json:"points"`
	}
	err = json.Unmarshal(res.Body.Bytes(), &pointsResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling GET response: %v", err)
	}

	assert.GreaterOrEqual(t, pointsResponse.Points, 0, "Expected non-negative points")
	assert.Equal(t, 28, pointsResponse.Points, "Expected 28 points for this receipt")
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
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total: "9.00",
	}

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error processing receipt: %v", err)
	}
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 for POST")

	var postResponse struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &postResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling POST response: %v", err)
	}
	receiptID := postResponse.ID
	assert.NotEmpty(t, receiptID, "POST response should contain a non-empty ID")

	// Now, fetch points using the valid ID
	endpoint := fmt.Sprintf("/receipts/%s/points", receiptID)
	res, err := makeRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatalf("Error making GET request: %v", err)
	}

	assert.Equal(t, http.StatusOK, res.Code, "Expected status code 200 for GET")

	var pointsResponse struct {
		Points int `json:"points"`
	}
	err = json.Unmarshal(res.Body.Bytes(), &pointsResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling GET response: %v", err)
	}
	assert.GreaterOrEqual(t, pointsResponse.Points, 0, "Expected non-negative points")
	assert.Equal(t, 109, pointsResponse.Points, "Expected 109 points for this receipt")
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

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error processing receipt: %v", err)
	}
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 for POST")

	var postResponse struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &postResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling POST response: %v", err)
	}
	receiptID := postResponse.ID
	assert.NotEmpty(t, receiptID, "POST response should contain a non-empty ID")

	// Now, fetch points using the valid ID
	endpoint := fmt.Sprintf("/receipts/%s/points", receiptID)
	res, err := makeRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatalf("Error making GET request: %v", err)
	}

	assert.Equal(t, http.StatusOK, res.Code, "Expected status code 200 for GET")

	var pointsResponse struct {
		Points int `json:"points"`
	}
	err = json.Unmarshal(res.Body.Bytes(), &pointsResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling GET response: %v", err)
	}
	assert.GreaterOrEqual(t, pointsResponse.Points, 0, "Expected non-negative points")
	assert.Equal(t, 15, pointsResponse.Points, "Expected 15 points for this receipt")
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

	w, err := makeRequest("POST", "/receipts/process", receipt)
	if err != nil {
		t.Fatalf("Error processing receipt: %v", err)
	}
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 for POST")

	var postResponse struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &postResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling POST response: %v", err)
	}
	receiptID := postResponse.ID
	assert.NotEmpty(t, receiptID, "POST response should contain a non-empty ID")

	// Now, fetch points using the valid ID
	endpoint := fmt.Sprintf("/receipts/%s/points", receiptID)
	res, err := makeRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatalf("Error making GET request: %v", err)
	}

	assert.Equal(t, http.StatusOK, res.Code, "Expected status code 200 for GET")

	var pointsResponse struct {
		Points int `json:"points"`
	}
	err = json.Unmarshal(res.Body.Bytes(), &pointsResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling GET response: %v", err)
	}
	assert.GreaterOrEqual(t, pointsResponse.Points, 0, "Expected non-negative points")
	assert.Equal(t, 31, pointsResponse.Points, "Expected 31 points for this receipt")
}

func TestGetReceiptsPoints_InvalidID(t *testing.T) {
	// Try to fetch points with an invalid ID
	endpoint := "/receipts/invalid-id/points"
	res, err := makeRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatalf("Error making GET request: %v", err)
	}

	assert.Equal(t, http.StatusNotFound, res.Code, "Expected status code 404 for invalid ID")
}
