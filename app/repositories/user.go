package repositories

import "Technopark_DB_Project/app/models"

type UserRepository interface {
	Create(user *models.User) (err error)
	Update(user *models.User) (err error)
	GetByNickname(nickname string) (user *models.User, err error)
	GetAllMatchedUsers(user *models.User) (users *[]models.User, err error)
}
