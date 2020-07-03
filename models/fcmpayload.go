package models

type Fcmpayload struct {
	Data Fcmdata `json:"data"`
	To string `json:"to"`
}