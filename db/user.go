package db

import (
	"fmt"

	"github.com/TKfleBR/GolangRawWeb/db/config"
	"github.com/TKfleBR/GolangRawWeb/models"
	"gopkg.in/mgo.v2/bson"
)

//InsertUser dsf
func InsertUser(u models.User) (err error) {
	u.ID = bson.NewObjectId()
	fmt.Println(u)
	conn := config.NewUserController()
	err = conn.DB("data").C("users").Insert(u)
	return
}

//GetUser sda
func GetUser(u *models.User) (newU *models.User) {
	conn := config.NewUserController()

	err := conn.DB("data").C("users").Find(bson.M{"name": u.Name}).One(&newU)
	if err != nil {
		fmt.Println("error")
	}
	return
}
