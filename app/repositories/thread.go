package repositories

import "Technopark_DB_Project/app/models"

type ThreadRepository interface {
	Create(thread *models.Thread) (err error)
	GetByID(id int64) (thread *models.Thread, err error)
	GetBySlug(slug string) (thread *models.Thread, err error)
	GetBySlugOrID(slugOrID string) (thread *models.Thread, err error)
	GetVotes(id int64) (votesAmount int32, err error)
	Update(thread *models.Thread) (err error)
	CreatePosts(thread *models.Thread, posts *models.Posts) (err error)
	GetPostsTree(threadID int64, limit, since int, desc bool) (posts *[]models.Post, err error)
	GetPostsParentTree(threadID int64, limit, since int, desc bool) (posts *[]models.Post, err error)
	GetPostsFlat(threadID int64, limit, since int, desc bool) (posts *[]models.Post, err error)
}
