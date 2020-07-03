package models

type Feedback struct {
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Telnum      string    `json:"telnum"`
	Email       string    `json:"email"`
	Agree       bool      `json:"agree"`
	ContactType string    `json:"contactType"`
	Message     string    `json:"message"`
	Date        string    `json:"date"`
	ID          int       `json:"id"`
}