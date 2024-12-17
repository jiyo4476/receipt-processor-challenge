package models

import (
	"fmt"
	"math"
	"regexp"
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

var oddRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d[13579]$`)
var roundAmountRegex = regexp.MustCompile(`^\d+\.00$`)
var multipleOf25Regex = regexp.MustCompile(`^\d+\.(00|25|50|75)$`)
var timeOfPurchaseRegex = regexp.MustCompile(`^1[4,5]:[0-5][0-9]$`)

func (r Receipt) ValidateTotal() error {
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
	points += r.getPointsAlphanumerical()

	// 50 points if the total is a round dollar amount with no cents
	points += r.getPointsRoundAmount()

	// 25 points if total is a multiple of .25
	points += r.getPointsMultipleOf25()

	// 5 points for every two items on the receipt
	points += r.getPointsForItemNum()

	// if trimmed length of item description is a multiple of 3, multiply price by
	// 0.2 and round up to the nearest int. The result is the number of points added
	itemPoints, err := r.getPointsForItems()
	if err != nil {
		// return error
		return -1, err
	}

	points += itemPoints

	// 6 points in the day in the purchsae date is odd
	points += r.getPointsForOddDate()

	// 10 points if the time of purchase is after 2pm but before 4pm
	points += r.getPointsForTimeOfPurchase()

	return points, err
}

func (r Receipt) getPointsAlphanumerical() int64 {
	count := int64(0)
	for _, c := range r.Retailer {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			count++
		}
	}
	return count
}

func (r Receipt) getPointsRoundAmount() int64 {
	if roundAmountRegex.MatchString(r.Total) {
		return 50
	}
	return 0
}

func (r Receipt) getPointsMultipleOf25() int64 {
	if multipleOf25Regex.MatchString(r.Total) {
		return 25
	}
	return 0
}

func (r Receipt) getPointsForOddDate() int64 {
	match := oddRegex.MatchString(r.PurchaseDate)
	if match {
		return 6
	}
	return 0
}

func (r Receipt) getPointsForTimeOfPurchase() int64 {
	if matched := timeOfPurchaseRegex.MatchString(r.PurchaseTime); matched {
		return 10
	}
	return 0
}

func (r Receipt) getPointsForItemNum() int64 {
	return int64((len(r.Items) / 2) * 5)
}

func (r Receipt) getPointsForItems() (int64, error) {
	points := int64(0)
	for _, curr_item := range r.Items {
		if len(strings.Trim(curr_item.ShortDescription, " "))%3 == 0 {
			i, err := strconv.ParseFloat(curr_item.Price, 64)
			if err != nil {
				// return error
				return -1, err
			}

			points += int64(math.Ceil(i * 0.2))
		}
	}
	return points, nil
}
