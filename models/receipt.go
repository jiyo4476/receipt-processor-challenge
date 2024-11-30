package models

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Receipt struct {
	Retailer     string `json:"retailer" binding:"required,correctRetailerName"`
	PurchaseDate string `json:"purchaseDate" binding:"required" time_format:"2022-01-01"`
	PurchaseTime string `json:"purchaseTime" binding:"required" time_format:"13:01"`
	Items        []Item `json:"items" binding:"required"`
	Total        string `json:"total" binding:"required,correctCashValue"`
}

type points interface {
	Points() int
}

func (r Receipt) Points() int {
	points := 0
	// One point for every alphanumerical character in the retailer name
	points += getNumAlphanumerical(r.Retailer)

	cents := r.Total[len(r.Total)-3:]
	// 50 points if the total is a round dollar amount with no cents
	points += getPointsRoundAmount(cents)
	// 25 points if total is a multiple of .25
	points += getPointsMultipleOf25(cents)

	// 5 points for every two items on the receipt
	points += (len(r.Items) / 2) * 5

	// if trimmed length of item description is a multiple of 3, multiply price by
	// 0.2 and round up to the nearest int. The result is the number of points added
	itemPoints, err := getPointsForItems(r.Items)
	if err != nil {
		// handle error
		return -1
	}

	points += itemPoints

	// 6 points in the day in the purchsae date is odd
	i, err := strconv.Atoi(r.PurchaseDate[len(r.PurchaseDate)-1:])
	if err != nil {
		return -1
	}

	points += getPointsForOddDate(i)

	// 10 points if the time of purchase is after 2pm but before 4pm
	points += getPointsForTimeOfPurchase(r.PurchaseTime[:2])

	return points
}

func getNumAlphanumerical(s string) int {
	count := 0
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			count++
		}
	}
	return count
}

func getPointsRoundAmount(cents string) int {
	if cents == ".00" {
		return 50
	}
	return 0
}

func getPointsMultipleOf25(cents string) int {
	if cents == ".00" || cents == ".25" || cents == ".50" || cents == ".75" {
		return 25
	}
	return 0
}

func getPointsForOddDate(day int) int {
	if day%2 == 1 {
		return 6
	}
	return 0
}

func getPointsForTimeOfPurchase(receiptHour string) int {
	if receiptHour == "14" || receiptHour == "15" {
		return 10
	}
	return 0
}

func getPointsForItems(items []Item) (int, error) {
	points := 0
	for _, curr_item := range items {
		if len(strings.Trim(curr_item.ShortDescription, " "))%3 == 0 {
			i, err := strconv.ParseFloat(curr_item.Price, 64)
			if err != nil {
				// handle error
				return -1, err
			}

			points += int(math.Ceil(i * 0.2))
		}
	}
	return points, nil
}
