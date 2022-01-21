package stores

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"
	"time"

	"github.com/jackc/pgx"
)

type PostStore struct {
	db *pgx.ConnPool
}

func CreatePostRepository(db *pgx.ConnPool) repositories.PostRepository {
	return &PostStore{db: db}
}

func (postStore *PostStore) GetByID(id int64) (post *models.Post, err error) {
	post = &models.Post{}
	postTime := time.Time{}
	err = postStore.db.QueryRow("SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts "+
		"WHERE id = $1", id).
		Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &postTime)
	post.Created = postTime.Format(time.RFC3339)
	return
}

func (postStore *PostStore) Update(post *models.Post) (err error) {
	_, err = postStore.db.Exec("UPDATE posts SET message = $1, is_edited = $2 WHERE id = $3;", post.Message, post.IsEdited, post.ID)
	return
}
