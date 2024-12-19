package store

import "github.com/jiyo4476/receipt-processor-challenge/models"

// Map for in-memory data storage
var Receipts = make(map[string]models.Receipt)
