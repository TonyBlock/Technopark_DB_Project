package impl

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"
	"Technopark_DB_Project/app/usecases"
	"Technopark_DB_Project/pkg/errors"
	"Technopark_DB_Project/pkg/validator"
)

type UserUseCaseImpl struct {
	userRepository repositories.UserRepository
}

func CreateUserUseCase(userRepository repositories.UserRepository) usecases.UserUseCase {
	return &UserUseCaseImpl{userRepository: userRepository}
}

func (userUseCase *UserUseCaseImpl) Create(user *models.User) (users *models.Users, err error) {
	if !validator.ValidateUserData(user, false) {
		err = errors.ErrBadInputData
		return
	}

	usersSlice, err := userUseCase.userRepository.GetAllMatchedUsers(user)
	if err != nil {
		return
	} else if len(*usersSlice) > 0 {
		users = new(models.Users)
		*users = *usersSlice
		err = errors.ErrUserAlreadyExist
		return
	}

	err = userUseCase.userRepository.Create(user)
	return
}

func (userUseCase *UserUseCaseImpl) Get(nickname string) (user *models.User, err error) {
	user, err = userUseCase.userRepository.GetByNickname(nickname)
	if err != nil {
		err = errors.ErrUserNotFound
	}
	return
}

func (userUseCase *UserUseCaseImpl) Update(user *models.User) (err error) {
	if !validator.ValidateUserData(user, true) {
		err = errors.ErrBadInputData
		return
	}

	_, err = userUseCase.userRepository.GetByNickname(user.Nickname)
	if err != nil {
		err = errors.ErrUserNotFound
	}

	err = userUseCase.userRepository.Update(user)
	if err != nil {
		err = errors.ErrUserDataConflict
	}
	return
}
