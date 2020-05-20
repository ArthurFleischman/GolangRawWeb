package auth

import (
	"github.com/TKfleBR/GolangRawWeb/db"
	"github.com/TKfleBR/GolangRawWeb/models"
	"golang.org/x/crypto/bcrypt"
)

//User wdf
func User(u *models.User, password []byte) bool {
	gotUser := db.GetUser(u)
	if gotUser != nil {
		if err := bcrypt.CompareHashAndPassword(gotUser.Password, password); err == nil {
			return true
		}
	}
	return false
}
