package models

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Receipt struct {
	Retailer     string `json:"retailer" binding:"required,min=1,correctRetailerName"`
	PurchaseDate string `json:"purchaseDate" binding:"required,len=10,correctDate" time_format:"2022-01-01"`
	PurchaseTime string `json:"purchaseTime" binding:"required,len=5,correctTime" time_format:"13:01"`
	Items        []Item `json:"items" binding:"required,dive"`
	Total        string `json:"total" binding:"required,min=4,correctCashValue"`
}

func (r Receipt) Validate() error {
	total, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return fmt.Errorf("invalid total: %w", err)
	}

	itemsTotal := 0.0
	for _, item := range r.Items {
		itemPrice, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return fmt.Errorf("invalid item price: %w", err)
		}
		itemsTotal += itemPrice
	}

	if fmt.Sprintf("%.2f", itemsTotal) != fmt.Sprintf("%.2f", total) {
		return fmt.Errorf("sum of item prices %.2f does not equal total %.2f", itemsTotal, total)
	}

	return nil
}

func (r Receipt) Points() (int64, error) {
	points := int64(0)
	// One point for every alphanumerical character in the retailer name
	points += getNumAlphanumerical(r.Retailer)

	cents := r.Total[len(r.Total)-2:]
	// 50 points if the total is a round dollar amount with no cents
	points += getPointsRoundAmount(cents)
	// 25 points if total is a multiple of .25
	points += getPointsMultipleOf25(cents)

	// 5 points for every two items on the receipt
	points += int64((len(r.Items) / 2) * 5)

	// if trimmed length of item description is a multiple of 3, multiply price by
	// 0.2 and round up to the nearest int. The result is the number of points added
	itemPoints, err := getPointsForItems(r.Items)
	if err != nil {
		// handle error
		log.Printf("Validation error: %v", err) // Log the error
		return -1, err
	}

	points += itemPoints

	// 6 points in the day in the purchsae date is odd
	i, err := strconv.Atoi(r.PurchaseDate[len(r.PurchaseDate)-1:])
	if err != nil {
		log.Printf("Validation error: %v", err) // Log the error
		return -1, err
	}

	points += getPointsForOddDate(i)

	// 10 points if the time of purchase is after 2pm but before 4pm
	points += getPointsForTimeOfPurchase(r.PurchaseTime[:2])

	return points, err
}

func getNumAlphanumerical(s string) int64 {
	count := int64(0)
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			count++
		}
	}
	return count
}

func getPointsRoundAmount(cents string) int64 {
	if cents == "00" {
		return 50
	}
	return 0
}

func getPointsMultipleOf25(cents string) int64 {
	if cents == "00" || cents == "25" || cents == "50" || cents == "75" {
		return 25
	}
	return 0
}

func getPointsForOddDate(day int) int64 {
	if day%2 == 1 {
		return 6
	}
	return 0
}

func getPointsForTimeOfPurchase(receiptHour string) int64 {
	if receiptHour == "14" || receiptHour == "15" {
		return 10
	}
	return 0
}

func getPointsForItems(items []Item) (int64, error) {
	points := int64(0)
	for _, curr_item := range items {
		if len(strings.Trim(curr_item.ShortDescription, " "))%3 == 0 {
			i, err := strconv.ParseFloat(curr_item.Price, 64)
			if err != nil {
				// handle error
				return -1, err
			}

			points += int64(math.Ceil(i * 0.2))
		}
	}
	return points, nil
}
