package stores

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"

	"github.com/jackc/pgx"
)

type ForumStore struct {
	db *pgx.ConnPool
}

func CreateForumRepository(db *pgx.ConnPool) repositories.ForumRepository {
	return &ForumStore{db: db}
}

func (forumStore *ForumStore) Create(forum *models.Forum) (err error) {
	_, err = forumStore.db.Exec("INSERT INTO forums (title, user_, slug) VALUES ($1, $2, $3);",
		forum.Title, forum.User, forum.Slug)
	return
}

func (forumStore *ForumStore) GetBySlug(slug string) (forum *models.Forum, err error) {
	forum = &models.Forum{}
	err = forumStore.db.QueryRow("SELECT title, user_, slug, posts, threads FROM forums WHERE slug = $1", slug).
		Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	return
}

func (forumStore *ForumStore) GetUsers(slug string, limit int, since string, desc bool) (users *[]models.User, err error) {
	var usersSlice []models.User

	var resultRows *pgx.Rows

	query := "SELECT users.nickname, users.fullname, users.about, users.email FROM users " +
		"LEFT JOIN user_forum ON users.nickname = user_forum.nickname WHERE user_forum.forum = $1"

	if since != "" {
		if desc {
			query += " AND users.nickname <= $2 ORDER BY users.nickname DESC"
		} else {
			query += " AND users.nickname >= $2 ORDER BY users.nickname"
		}
		query += " LIMIT = $3;"
		resultRows, err = forumStore.db.Query(query, slug, since, limit)
	} else {
		if desc {
			query += " ORDER BY users.nickname DESC"
		} else {
			query += " ORDER BY users.nickname"
		}
		query += " LIMIT = $2;"
		resultRows, err = forumStore.db.Query(query, slug, limit)
	}

	if err != nil {
		return
	}
	defer resultRows.Close()

	for resultRows.Next() {
		user := models.User{}
		err = resultRows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return
		}
		usersSlice = append(usersSlice, user)
	}
	return &usersSlice, nil
}

func (forumStore *ForumStore) GetThreads(slug string, limit int, since string, desc bool) (threads *[]models.Thread, err error) {
	var threadsSlice []models.Thread

	var resultRows *pgx.Rows

	query := "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE forum = $1"

	if since != "" {
		if desc {
			query += " AND created <= $2 ORDER BY created DESC"
		} else {
			query += " AND created >= $2 ORDER BY created"
		}
		query += " LIMIT = $3;"
		resultRows, err = forumStore.db.Query(query, slug, since, limit)
	} else {
		if desc {
			query += " ORDER BY created DESC"
		} else {
			query += " ORDER BY created"
		}
		query += " LIMIT = $2;"
		resultRows, err = forumStore.db.Query(query, slug, limit)
	}

	if err != nil {
		return
	}
	defer resultRows.Close()

	for resultRows.Next() {
		thread := models.Thread{}
		err = resultRows.Scan(&thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return
		}
		threadsSlice = append(threadsSlice, thread)
	}
	return &threadsSlice, nil
}
