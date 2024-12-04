package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"receipt-processor-challenge/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Map for in-memory data storage
var memory_cache = make(map[string]models.Receipt)

// Returns an ID for the receipt
func processReceipt(c *gin.Context) {
	var receipt models.Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		log.Printf("Error binding JSON: %v", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := receipt.Validate(); err != nil {
		log.Printf("Validation error: %v", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id = uuid.New().String()
	memory_cache[id] = receipt
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getReceiptsPoints(c *gin.Context) {
	receipt, ok := memory_cache[c.Param("id")]
	if ok {
		points, _ := receipt.Points()
		c.JSON(http.StatusOK, gin.H{
			"points": points,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"code": "ID_NOT_FOUND", "message": "No receipt found for that id"})
	}
}

func loadConfig(configFile string) (*libopenapi.DocumentModel[v3.Document], error) {
	// Load config
	spec, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load config from %s: %e", configFile, err)
	}

	specDocument, err := libopenapi.NewDocument(spec)
	if err != nil {
		return nil, fmt.Errorf("failed creating config from %s: %e", configFile, err)
	}

	docModel, errors := specDocument.BuildV3Model()

	if len(errors) > 0 {
		for i := range errors {
			log.Printf("error: %e\n", errors[i])
		}
		log.Printf("cannot create openApi v3 model from document: %d errors reported", len(errors))
		return nil, fmt.Errorf("cannot create openApi v3 model from document: %d errors reported", len(errors))
	}

	return docModel, nil
}

// Use a single instance of Validate, it caches struct info
func SetUpRouter() *gin.Engine {
	router := gin.Default()

	// Register custom validation functions for the test router
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("correctRetailerName", models.CorrectRetailerName)
		v.RegisterValidation("correctShortDescription", models.CorrectShortDescription)
		v.RegisterValidation("correctCashValue", models.CorrectCashValue)
		v.RegisterValidation("correctDate", models.CorrectDate)
		v.RegisterValidation("correctTime", models.CorrectTime)
	}

	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getReceiptsPoints)
	return router
}

func main() {
	docModel, err := loadConfig("api.yml")
	if err != nil {
		log.Printf("Error loading config: %v", err) // Log the error
		return
	}

	fmt.Printf("Starting %s %s - %s\n\n", docModel.Model.Info.Title, docModel.Model.Info.Version, docModel.Model.Info.Description)

	router := SetUpRouter()

	router.Run()
}
