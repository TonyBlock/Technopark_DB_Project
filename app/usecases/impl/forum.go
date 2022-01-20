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
	oldForum, err := forumUseCase.forumRepository.GetBySlug(forum.Slug)
	if err == nil {
		forum = oldForum
		err = errors.ErrForumAlreadyExists
		return
	}

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
	_, err = forumUseCase.forumRepository.GetBySlug(thread.Forum)
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
	if err == nil {
		thread = oldThread
		err = errors.ErrThreadAlreadyExists
		return
	}

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
	*users = *usersSlice

	return
}

func (forumUseCase *ForumUseCaseImpl) GetThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error) {
	_, err = forumUseCase.forumRepository.GetBySlug(slug)
	if err != nil {
		err = errors.ErrForumNotExist
		return
	}

	threadsSlice, err := forumUseCase.forumRepository.GetThreads(slug, limit, since, desc)
	if err != nil {
		return
	}
	threads = new(models.Threads)
	*threads = *threadsSlice

	return
}
