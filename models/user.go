package models

type User struct {
	Name   		string        `json:"name" `
	Email 		string        `json:"email" `
	Password    string        `json:"password" `
	Mobile		int64 		  `json:"mobile" `
	Address  	string 		  `json:"address" `
	ProfilePic  string		  `json:"image"`
}