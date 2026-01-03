package models

type User struct {
	UUID       string `json:"uuid"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email"`
}
