package models

type User struct {
	UUID       string `json:"uuid" bson:"uuid"`
	Firstname  string `json:"firstname" bson:"firstname"`
	Fullname   string `json:"fullname" bson:"fullname"`
	ProfileImg string `json:"profile_img" bson:"profile_img"`
}

type UserDetail struct {
	User       `bson:",inline"`
	Lastname   string `json:"lastname" bson:"lastname"`
	Middlename string `json:"middlename" bson:"middlename"`
	Email      string `json:"email" bson:"email"`
}
