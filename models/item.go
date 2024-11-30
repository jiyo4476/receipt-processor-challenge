package models

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required,correctShortDescription"`
	Price            string `json:"price" binding:"required,correctCashValue"`
}
