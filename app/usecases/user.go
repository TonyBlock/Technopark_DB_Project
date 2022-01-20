package usecases

import (
	"Technopark_DB_Project/app/models"
)

type UserUseCase interface {
	Create(user *models.User) (users *models.Users, err error)
	Get(nickname string) (user *models.User, err error)
	Update(user *models.User) (err error)
}
