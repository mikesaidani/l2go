package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Username    string        `bson:"username"`
	Password    string        `bson:"password"`
	AccessLevel int8          `bson:"access_level"`
}
