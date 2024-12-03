package models

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required,min=1,correctShortDescription"`
	Price            string `json:"price" binding:"required,min=4,correctCashValue"`
}
