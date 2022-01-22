package stores

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"
	"Technopark_DB_Project/pkg/errors"
	"fmt"
	"time"

	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type ThreadStore struct {
	db *pgx.ConnPool
}

func CreateThreadRepository(db *pgx.ConnPool) repositories.ThreadRepository {
	return &ThreadStore{db: db}
}

func (threadStore *ThreadStore) Create(thread *models.Thread) (err error) {
	err = threadStore.db.QueryRow("INSERT INTO threads (title, author, forum, message, slug, created) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created;",
		thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).
		Scan(&thread.ID, &thread.Created)
	return
}

func (threadStore *ThreadStore) GetByID(id int64) (thread *models.Thread, err error) {
	thread = &models.Thread{}
	err = threadStore.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created FROM threads "+
		"WHERE id = $1;", id).
		Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	return
}

func (threadStore *ThreadStore) GetBySlug(slug string) (thread *models.Thread, err error) {
	thread = &models.Thread{}
	err = threadStore.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created FROM threads "+
		"WHERE slug = $1;", slug).
		Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	return
}

func (threadStore *ThreadStore) GetBySlugOrID(slugOrID string) (thread *models.Thread, err error) {
	thread = &models.Thread{}
	err = threadStore.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created FROM threads "+
		"WHERE id = $1 OR slug = $2;", slugOrID, slugOrID).
		Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	return
}

func (threadStore *ThreadStore) GetVotes(id int64) (votesAmount int32, err error) {
	err = threadStore.db.QueryRow("SELECT votes FROM threads WHERE id = $1;", id).Scan(&votesAmount)
	return
}

func (threadStore *ThreadStore) Update(thread *models.Thread) (err error) {
	_, err = threadStore.db.Exec("UPDATE threads SET "+
		"title = $1, message = $2 WHERE id = $3;", thread.Title, thread.Message, thread.ID)
	return
}

func (threadStore *ThreadStore) createPartPosts(thread *models.Thread, posts *models.Posts, from, to int, created time.Time, createdFormatted string) (err error) {
	query := "INSERT INTO posts (parent, author, message, forum, thread, created) VALUES "
	args := make([]interface{}, 0, 0)

	j := 0
	for i := from; i < to; i++ {
		(*posts)[i].Forum = thread.Forum
		(*posts)[i].Thread = thread.ID
		(*posts)[i].Created = createdFormatted
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", j*6+1, j*6+2, j*6+3, j*6+4, j*6+5, j*6+6)
		if (*posts)[i].Parent != 0 {
			args = append(args, (*posts)[i].Parent, (*posts)[i].Author, (*posts)[i].Message, thread.Forum, thread.ID, created)
		} else {
			args = append(args, nil, (*posts)[i].Author, (*posts)[i].Message, thread.Forum, thread.ID, created)
		}
		j++
	}
	query = query[:len(query)-1]
	query += " RETURNING id;"

	isSuccess := false
	k := 0

	for !isSuccess {

		resultRows, err := threadStore.db.Query(query, args...)
		if err != nil {
			fmt.Println(err)
			return errors.ErrParentPostNotExist
		}
		defer resultRows.Close()

		for i := from; resultRows.Next(); i++ {
			isSuccess = true
			var id int64
			if err = resultRows.Scan(&id); err != nil {
				return err
			}
			(*posts)[i].ID = id
		}
		k++
		if k >= 3 {
			break
		}
	}

	return
}

func (threadStore *ThreadStore) CreatePosts(thread *models.Thread, posts *models.Posts) (err error) {
	created := time.Now()
	createdFormatted := created.Format(time.RFC3339)

	parts := len(*posts) / 20
	for i := 0; i < parts+1; i++ {
		if i == parts {
			if i*20 != len(*posts) {
				err = threadStore.createPartPosts(thread, posts, i*20, len(*posts), created, createdFormatted)
				if err != nil {
					return err
				}
			}
		} else {
			err = threadStore.createPartPosts(thread, posts, i*20, i*20+20, created, createdFormatted)
			if err != nil {
				return err
			}
		}
	}

	return
}

func (threadStore *ThreadStore) GetPostsTree(threadID int64, limit, since int, desc bool) (posts *[]models.Post, err error) {
	var rows *pgx.Rows

	if since == -1 {
		if desc {
			query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts " +
				"WHERE thread = $1 ORDER BY path DESC LIMIT NULLIF($2, 0);"
			rows, err = threadStore.db.Query(query, threadID, limit)
		} else {
			query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts " +
				"WHERE thread = $1 ORDER BY path LIMIT NULLIF($2, 0);"
			rows, err = threadStore.db.Query(query, threadID, limit)
		}
	} else {
		if desc {
			query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts " +
				"WHERE thread = $1 AND path < (SELECT path FROM posts WHERE id = $2) ORDER BY path DESC LIMIT NULLIF($3, 0);"
			rows, err = threadStore.db.Query(query, threadID, since, limit)
		} else {
			query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts " +
				"WHERE thread = $1 AND path > (SELECT path FROM posts WHERE id = $2) ORDER BY path LIMIT NULLIF($3, 0);"
			rows, err = threadStore.db.Query(query, threadID, since, limit)
		}
	}

	if err != nil {
		return
	}
	defer rows.Close()

	posts = new([]models.Post)
	for rows.Next() {
		post := models.Post{}
		postTime := time.Time{}

		err = rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &postTime)
		if err != nil {
			return
		}

		post.Created = postTime.Format(time.RFC3339)
		*posts = append(*posts, post)
	}

	return
}

func (threadStore *ThreadStore) GetPostsParentTree(threadID int64, limit, since int, desc bool) (posts *[]models.Post, err error) {
	var rows *pgx.Rows

	if since == -1 {
		if desc {
			rows, err = threadStore.db.Query(`
					SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN 
						(SELECT id FROM posts WHERE thread = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path ASC, id ASC;`, threadID, limit)
		} else {
			rows, err = threadStore.db.Query(`
					SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN 
						(SELECT id FROM posts WHERE thread = $1 AND parent IS NULL ORDER BY id LIMIT $2)
					ORDER BY path;`, threadID, limit)
		}
	} else {
		if desc {
			rows, err = threadStore.db.Query(`
					SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN 
						(SELECT id FROM posts WHERE thread = $1 AND parent IS NULL AND path[1] <
							(SELECT path[1] FROM posts WHERE id = $2) 
						ORDER BY id DESC LIMIT $3)
					ORDER BY path[1] DESC, path ASC, id ASC;`, threadID, since, limit)
		} else {
			rows, err = threadStore.db.Query(`
					SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN 
						(SELECT id FROM posts WHERE thread = $1 AND parent IS NULL AND path[1] >
							(SELECT path[1] FROM posts WHERE id = $2) 
						ORDER BY id LIMIT $3) 
					ORDER BY path;`, threadID, since, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts = new([]models.Post)
	for rows.Next() {
		post := models.Post{}
		postTime := time.Time{}

		err = rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &postTime)
		if err != nil {
			return
		}

		post.Created = postTime.Format(time.RFC3339)
		*posts = append(*posts, post)
	}

	return
}

func (threadStore *ThreadStore) GetPostsFlat(threadID int64, limit, since int, desc bool) (posts *[]models.Post, err error) {
	var rows *pgx.Rows

	if since == -1 {
		if desc {
			query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts WHERE thread = $1 ORDER BY id DESC LIMIT NULLIF($2, 0);"
			rows, err = threadStore.db.Query(query, threadID, limit)
		} else {
			query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts WHERE thread = $1 ORDER BY id LIMIT NULLIF($2, 0);"
			rows, err = threadStore.db.Query(query, threadID, limit)
		}
	} else {
		if desc {
			query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts WHERE thread = $1 AND id < $2 ORDER BY id DESC LIMIT NULLIF($3, 0);"
			rows, err = threadStore.db.Query(query, threadID, since, limit)
		} else {
			query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts WHERE thread = $1 AND id > $2 ORDER BY id LIMIT NULLIF($3, 0);"
			rows, err = threadStore.db.Query(query, threadID, since, limit)
		}
	}
	if err != nil {
		return
	}

	defer rows.Close()
	posts = new([]models.Post)
	for rows.Next() {
		post := models.Post{}
		postTime := time.Time{}

		err = rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &postTime)
		if err != nil {
			return
		}

		post.Created = postTime.Format(time.RFC3339)
		*posts = append(*posts, post)
	}

	return
}
