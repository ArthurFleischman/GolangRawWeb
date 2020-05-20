package config

import (
	"gopkg.in/mgo.v2"
)

//UserController connector

//NewUserController jf
func NewUserController() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	return s
}
