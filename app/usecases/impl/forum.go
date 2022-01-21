package impl

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"
	"Technopark_DB_Project/app/usecases"
	"Technopark_DB_Project/pkg/errors"
)

type ForumUseCaseImpl struct {
	forumRepository  repositories.ForumRepository
	threadRepository repositories.ThreadRepository
	userRepository   repositories.UserRepository
}

func CreateForumUseCase(forumRepository repositories.ForumRepository, threadRepository repositories.ThreadRepository, userRepository repositories.UserRepository) usecases.ForumUseCase {
	return &ForumUseCaseImpl{forumRepository: forumRepository, threadRepository: threadRepository, userRepository: userRepository}
}

func (forumUseCase *ForumUseCaseImpl) CreateForum(forum *models.Forum) (err error) {
	user, err := forumUseCase.userRepository.GetByNickname(forum.User)
	if err != nil {
		err = errors.ErrUserNotFound
		return
	}

	oldForum, err := forumUseCase.forumRepository.GetBySlug(forum.Slug)
	if oldForum.Slug != "" {
		*forum = *oldForum
		err = errors.ErrForumAlreadyExists
		return
	}

	forum.User = user.Nickname
	err = forumUseCase.forumRepository.Create(forum)
	return
}

func (forumUseCase *ForumUseCaseImpl) Get(slug string) (forum *models.Forum, err error) {
	forum, err = forumUseCase.forumRepository.GetBySlug(slug)
	if err != nil {
		err = errors.ErrForumNotExist
	}
	return
}

func (forumUseCase *ForumUseCaseImpl) CreateThread(thread *models.Thread) (err error) {
	forum, err := forumUseCase.forumRepository.GetBySlug(thread.Forum)
	if err != nil {
		err = errors.ErrForumOrTheadNotFound
		return
	}

	_, err = forumUseCase.userRepository.GetByNickname(thread.Author)
	if err != nil {
		err = errors.ErrForumOrTheadNotFound
		return
	}

	oldThread, err := forumUseCase.threadRepository.GetBySlug(thread.Slug)
	if oldThread.Slug != "" {
		*thread = *oldThread
		err = errors.ErrThreadAlreadyExists
		return
	}

	thread.Forum = forum.Slug
	err = forumUseCase.threadRepository.Create(thread)
	return
}

func (forumUseCase *ForumUseCaseImpl) GetUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error) {
	_, err = forumUseCase.forumRepository.GetBySlug(slug)
	if err != nil {
		err = errors.ErrForumNotExist
		return
	}

	usersSlice, err := forumUseCase.forumRepository.GetUsers(slug, limit, since, desc)
	if err != nil {
		return
	}
	users = new(models.Users)
	if len(*usersSlice) == 0 {
		*users = []models.User{}
	} else {
		*users = *usersSlice
	}

	return
}

func (forumUseCase *ForumUseCaseImpl) GetThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error) {
	forum, err := forumUseCase.forumRepository.GetBySlug(slug)
	if err != nil {
		err = errors.ErrForumNotExist
		return
	}

	threadsSlice, err := forumUseCase.forumRepository.GetThreads(forum.Slug, limit, since, desc)
	if err != nil {
		return
	}
	threads = new(models.Threads)
	if len(*threadsSlice) == 0 {
		*threads = []models.Thread{}
	} else {
		*threads = *threadsSlice
	}

	return
}
