package models

type Promotions struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Label       string `json:"label"`
	Price       string `json:"price"`
	Featured    bool   `json:"featured"`
	Description string `json:"description"`
}
