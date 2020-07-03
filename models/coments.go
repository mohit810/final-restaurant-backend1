package models

type Comments struct {
	DishID  int       `json:"dishId"`
	Rating  string    `json:"rating"`
	Author  string     `json:"author"`
	Comment string    `json:"comment"`
	Date    string    `json:"date"`
	ID      int       `json:"id"`
}
