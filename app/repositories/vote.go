package repositories

import "Technopark_DB_Project/app/models"

type VoteRepository interface {
	Vote(threadID int64, vote *models.Vote) (err error)
}
