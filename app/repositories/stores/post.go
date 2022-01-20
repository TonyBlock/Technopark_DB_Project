package stores

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"

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
	err = postStore.db.QueryRow("SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts "+
		"WHERE id = $1", id).
		Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	return
}

func (postStore *PostStore) Update(post *models.Post) (err error) {
	_, err = postStore.db.Exec("UPDATE posts SET"+
		"message = $1 WHERE id = $2;", post.Message, post.ID)
	return
}
