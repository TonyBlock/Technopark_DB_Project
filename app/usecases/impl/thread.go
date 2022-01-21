package impl

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"
	"Technopark_DB_Project/app/usecases"
	"Technopark_DB_Project/pkg/errors"
	"strconv"
)

type ThreadUseCaseImpl struct {
	threadRepository repositories.ThreadRepository
	voteRepository   repositories.VoteRepository
	postRepository   repositories.PostRepository
	userRepository   repositories.UserRepository
}

func CreateThreadUseCase(
	threadRepository repositories.ThreadRepository,
	voteRepository repositories.VoteRepository,
	postRepository repositories.PostRepository,
	userRepository repositories.UserRepository,
) usecases.ThreadUseCase {
	return &ThreadUseCaseImpl{threadRepository: threadRepository, voteRepository: voteRepository, postRepository: postRepository, userRepository: userRepository}
}

func (threadUseCase *ThreadUseCaseImpl) CreatePosts(slugOrID string, posts *models.Posts) (err error) {
	id, errConv := strconv.Atoi(slugOrID)
	var thread *models.Thread
	if errConv != nil {
		thread, err = threadUseCase.threadRepository.GetBySlug(slugOrID)
	} else {
		thread, err = threadUseCase.threadRepository.GetByID(int64(id))
	}

	if err != nil {
		err = errors.ErrThreadNotFound
		return
	}

	if len(*posts) == 0 {
		return
	}

	if (*posts)[0].Parent != 0 {
		var parentPost *models.Post
		parentPost, err = threadUseCase.postRepository.GetByID((*posts)[0].Parent)
		if parentPost.Thread != thread.ID {
			err = errors.ErrParentPostFromOtherThread
			return
		}
	}
	_, err = threadUseCase.userRepository.GetByNickname((*posts)[0].Author)
	if err != nil {
		err = errors.ErrUserNotFound
		return
	}

	err = threadUseCase.threadRepository.CreatePosts(thread, posts)
	return
}

func (threadUseCase *ThreadUseCaseImpl) Get(slugOrID string) (thread *models.Thread, err error) {
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = threadUseCase.threadRepository.GetBySlug(slugOrID)
	} else {
		thread, err = threadUseCase.threadRepository.GetByID(int64(id))
	}
	if err != nil {
		err = errors.ErrThreadNotFound
		return
	}
	return
}

func (threadUseCase *ThreadUseCaseImpl) Update(slugOrID string, thread *models.Thread) (err error) {
	id, errConv := strconv.Atoi(slugOrID)
	var oldThread *models.Thread
	if errConv != nil {
		oldThread, err = threadUseCase.threadRepository.GetBySlug(slugOrID)
	} else {
		oldThread, err = threadUseCase.threadRepository.GetByID(int64(id))
	}

	if err != nil {
		err = errors.ErrThreadNotFound
		return
	}

	if thread.Title != "" {
		oldThread.Title = thread.Title
	}
	if thread.Message != "" {
		oldThread.Message = thread.Message
	}

	err = threadUseCase.threadRepository.Update(oldThread)
	if err != nil {
		return
	}

	*thread = *oldThread

	return
}

func (threadUseCase *ThreadUseCaseImpl) GetPosts(slugOrID string, limit, since int, sort string, desc bool) (posts *models.Posts, err error) {
	id, errConv := strconv.Atoi(slugOrID)
	var thread *models.Thread
	if errConv != nil {
		thread, err = threadUseCase.threadRepository.GetBySlug(slugOrID)
	} else {
		thread, err = threadUseCase.threadRepository.GetByID(int64(id))
	}

	if err != nil {
		err = errors.ErrThreadNotFound
		return
	}

	postsSlice := new([]models.Post)
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
	if len(*postsSlice) == 0 {
		*posts = []models.Post{}
	} else {
		*posts = *postsSlice
	}

	return
}

func (threadUseCase *ThreadUseCaseImpl) Vote(slugOrID string, vote *models.Vote) (thread *models.Thread, err error) {
	id, errConv := strconv.Atoi(slugOrID)

	if errConv != nil {
		thread, err = threadUseCase.threadRepository.GetBySlug(slugOrID)
	} else {
		thread, err = threadUseCase.threadRepository.GetByID(int64(id))
	}

	if err != nil {
		err = errors.ErrThreadNotFound
		return
	}

	err = threadUseCase.voteRepository.Vote(thread.ID, vote)
	if err != nil {
		err = errors.ErrUserNotFound
		return
	}
	thread.Votes, err = threadUseCase.threadRepository.GetVotes(thread.ID)
	//if threadUseCase.voteRepository.IsVoted(thread.ID, vote) {
	//	err = threadUseCase.voteRepository.Update(thread, vote)
	//} else {
	//	err = threadUseCase.voteRepository.Create(thread, vote)
	//}

	return
}
