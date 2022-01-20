package usecases

import "Technopark_DB_Project/app/models"

type ForumUseCase interface {
	CreateForum(forum *models.Forum) (err error)
	Get(slug string) (forum *models.Forum, err error)
	CreateThread(thread *models.Thread) (err error)
	GetUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error)
	GetThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error)
}
