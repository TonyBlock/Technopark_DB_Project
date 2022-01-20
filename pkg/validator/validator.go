package validator

import (
	"Technopark_DB_Project/app/models"
	"regexp"
)

var (
	regLatinNumUnderline = regexp.MustCompile("[a-zA-Z0-9_]*")
	regEmail             = regexp.MustCompile(".+@.+")
)

func ValidateUserData(user *models.User, isUpdate bool) (isValid bool) {
	if !isUpdate && (len(user.Fullname) == 0 || len(user.Email) == 0) {
		return false
	}
	if !regLatinNumUnderline.MatchString(user.Fullname) {
		return false
	}
	if !regEmail.MatchString(user.Email) {
		return false
	}
	return true
}
