package handlers

import "github.com/jiyo4476/receipt-processor-challenge/models"

// Map for in-memory data storage
var memory_cache = make(map[string]models.Receipt)
