package impl

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"
	"Technopark_DB_Project/app/usecases"
	"Technopark_DB_Project/pkg/errors"
)

type ThreadUseCaseImpl struct {
	threadRepository repositories.ThreadRepository
	voteRepository   repositories.VoteRepository
}

func CreateThreadUseCase(threadRepository repositories.ThreadRepository, voteRepository repositories.VoteRepository) usecases.ThreadUseCase {
	return &ThreadUseCaseImpl{threadRepository: threadRepository, voteRepository: voteRepository}
}

func (threadUseCase *ThreadUseCaseImpl) CreatePosts(slugOrID string, posts *models.Posts) (err error) {
	thread, err := threadUseCase.threadRepository.GetBySlugOrID(slugOrID)
	if err != nil {
		err = errors.ErrThreadNotFound
		return
	}

	err = threadUseCase.threadRepository.CreatePosts(thread, posts)
	return
}

func (threadUseCase *ThreadUseCaseImpl) Get(slugOrID string) (thread *models.Thread, err error) {
	thread, err = threadUseCase.threadRepository.GetBySlugOrID(slugOrID)
	if err != nil {
		err = errors.ErrThreadNotFound
	}
	return
}

func (threadUseCase *ThreadUseCaseImpl) Update(slugOrID string, thread *models.Thread) (err error) {
	oldThread, err := threadUseCase.threadRepository.GetBySlugOrID(slugOrID)
	if err != nil {
		err = errors.ErrThreadNotFound
		return
	}

	oldThread.Title = thread.Title
	oldThread.Message = thread.Message

	err = threadUseCase.threadRepository.Update(oldThread)
	if err != nil {
		return
	}

	thread = oldThread

	return
}

func (threadUseCase *ThreadUseCaseImpl) GetPosts(slugOrID string, limit, since int, sort string, desc bool) (posts *models.Posts, err error) {
	thread, err := threadUseCase.threadRepository.GetBySlugOrID(slugOrID)
	if err != nil {
		err = errors.ErrThreadNotFound
		return
	}

	var postsSlice *[]models.Post

	switch sort {
	case "tree":
		postsSlice, err = threadUseCase.threadRepository.GetPostsTree(thread.ID, limit, since, desc)
	case "parent_tree":
		postsSlice, err = threadUseCase.threadRepository.GetPostsParentTree(thread.ID, limit, since, desc)
	default:
		postsSlice, err = threadUseCase.threadRepository.GetPostsFlat(thread.ID, limit, since, desc)
	}
	if err != nil {
		return
	}
	posts = new(models.Posts)
	*posts = *postsSlice

	return
}

func (threadUseCase *ThreadUseCaseImpl) Vote(slugOrID string, vote *models.Vote) (thread *models.Thread, err error) {
	thread, err = threadUseCase.threadRepository.GetBySlugOrID(slugOrID)
	if err != nil {
		err = errors.ErrForumNotExist
		return
	}

	err = threadUseCase.voteRepository.Vote(thread.ID, vote)
	thread.Votes, err = threadUseCase.threadRepository.GetVotes(thread.ID)
	//if threadUseCase.voteRepository.IsVoted(thread.ID, vote) {
	//	err = threadUseCase.voteRepository.Update(thread, vote)
	//} else {
	//	err = threadUseCase.voteRepository.Create(thread, vote)
	//}

	return
}
