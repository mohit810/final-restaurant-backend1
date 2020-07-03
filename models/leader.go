package models

type Leader struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Designation string `json:"designation"`
	Abbr        string `json:"abbr"`
	Featured    bool   `json:"featured"`
	Description string `json:"description"`
}
