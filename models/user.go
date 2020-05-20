package models

import "gopkg.in/mgo.v2/bson"

//User la
type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Password []byte        `json:"password" bson:"password"`
}
