package usecases

import "Technopark_DB_Project/app/models"

type ServiceUseCase interface {
	Clear() (err error)
	GetStatus() (status *models.Status, err error)
}
