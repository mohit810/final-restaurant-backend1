package models

type Dishes struct{
	Id  		int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Label       string `json:"label"`
	Rating		float32 `json:"rating"`
	Time        int64	`json:"time"`
	Featured    bool	`json:"featured"`
}