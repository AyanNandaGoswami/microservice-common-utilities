package models

// User represents core user profile fields used across microservice data layers and APIs.
type User struct {
	UUID          string `json:"uuid" bson:"uuid"`
	Firstname     string `json:"firstname" bson:"firstname"`
	Fullname      string `json:"fullname" bson:"fullname"`
	ProfileImg    string `json:"profile_img" bson:"profile_img"`
	Country       string `json:"country" bson:"country"`
	ContactNumber string `json:"contact_number" bson:"contact_number"`
}

// UserDetail extends User with detailed personal profile parameters including email and additional name fields.
type UserDetail struct {
	User          `bson:",inline"`
	Lastname      string `json:"lastname" bson:"lastname"`
	Middlename    string `json:"middlename" bson:"middlename"`
	Email         string `json:"email" bson:"email"`
	IsSystemAdmin bool   `json:"is_system_admin" bson:"is_system_admin"`
}
