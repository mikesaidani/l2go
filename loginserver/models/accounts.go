package models

type Account struct {
	Username    string `bson:username`
	Password    string `bson:password`
	AccessLevel int8  `bson:access_level`
}
