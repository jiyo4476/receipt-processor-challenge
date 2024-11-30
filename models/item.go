package models

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required,correctShortDescription,min=1"`
	Price            string `json:"price" binding:"required,correctCashValue,min=4"`
}
