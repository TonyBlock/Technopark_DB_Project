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
		"WHERE LOWER(slug) = LOWER($1);", slug).
		Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	return
}

func (threadStore *ThreadStore) GetBySlugOrID(slugOrID string) (thread *models.Thread, err error) {
	thread = &models.Thread{}
	err = threadStore.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created FROM threads "+
		"WHERE id = $1 OR LOWER(slug) = LOWER($2);", slugOrID, slugOrID).
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
	resultRows, err := threadStore.db.Query(query, args...)
	if err != nil {
		return errors.ErrParentPostNotExist
	}
	defer resultRows.Close()

	for i := from; resultRows.Next(); i++ {
		var id int64
		if err = resultRows.Scan(&id); err != nil {
			return err
		}
		(*posts)[i].ID = id
	}
	return
}

func (threadStore *ThreadStore) CreatePosts(thread *models.Thread, posts *models.Posts) (err error) {
	created := time.Now()
	createdFormatted := created.Format(time.RFC3339)

	parts := len(*posts) / 30
	for i := 0; i < parts+1; i++ {
		if i == parts {
			err = threadStore.createPartPosts(thread, posts, i*30, len(*posts), created, createdFormatted)
			if err != nil {
				return err
			}
		} else {
			err = threadStore.createPartPosts(thread, posts, i*30, i*30+30, created, createdFormatted)
			if err != nil {
				return err
			}
		}
	}

	//j := 0
	//if len(*posts) > 45 {
	//	half := len(*posts) / 2
	//	{
	//		query := "INSERT INTO posts (parent, author, message, forum, thread, created) VALUES "
	//		args := make([]interface{}, 0, 0)
	//
	//		for i := 0; i < half; i++ {
	//			(*posts)[i].Forum = thread.Forum
	//			(*posts)[i].Thread = thread.ID
	//			(*posts)[i].Created = createdFormatted
	//			query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
	//			if (*posts)[i].Parent != 0 {
	//				args = append(args, (*posts)[i].Parent, (*posts)[i].Author, (*posts)[i].Message, thread.Forum, thread.ID, created)
	//			} else {
	//				args = append(args, nil, (*posts)[i].Author, (*posts)[i].Message, thread.Forum, thread.ID, created)
	//			}
	//		}
	//		query = query[:len(query)-1]
	//		query += " RETURNING id;"
	//		resultRows, err := threadStore.db.Query(query, args...)
	//		if err != nil {
	//			return errors.ErrParentPostNotExist
	//		}
	//		defer resultRows.Close()
	//
	//		for i := 0; resultRows.Next(); i++ {
	//			var id int64
	//			if err = resultRows.Scan(&id); err != nil {
	//				return err
	//			}
	//			(*posts)[i].ID = id
	//		}
	//	}
	//	{
	//		query := "INSERT INTO posts (parent, author, message, forum, thread, created) VALUES "
	//		args := make([]interface{}, 0, 0)
	//
	//		j := 0
	//		for i := half; i < len(*posts); i++ {
	//			(*posts)[i].Forum = thread.Forum
	//			(*posts)[i].Thread = thread.ID
	//			(*posts)[i].Created = createdFormatted
	//			query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", j*6+1, j*6+2, j*6+3, j*6+4, j*6+5, j*6+6)
	//			if (*posts)[i].Parent != 0 {
	//				args = append(args, (*posts)[i].Parent, (*posts)[i].Author, (*posts)[i].Message, thread.Forum, thread.ID, created)
	//			} else {
	//				args = append(args, nil, (*posts)[i].Author, (*posts)[i].Message, thread.Forum, thread.ID, created)
	//			}
	//			j++
	//		}
	//		query = query[:len(query)-1]
	//		query += " RETURNING id;"
	//		resultRows, err := threadStore.db.Query(query, args...)
	//		if err != nil {
	//			return errors.ErrParentPostNotExist
	//		}
	//		defer resultRows.Close()
	//
	//		for i := half; resultRows.Next(); i++ {
	//			var id int64
	//			if err = resultRows.Scan(&id); err != nil {
	//				return err
	//			}
	//			(*posts)[i].ID = id
	//		}
	//	}
	//} else {
	//	query := "INSERT INTO posts (parent, author, message, forum, thread, created) VALUES "
	//	args := make([]interface{}, 0, 0)
	//
	//	for i, post := range *posts {
	//		(*posts)[i].Forum = thread.Forum
	//		(*posts)[i].Thread = thread.ID
	//		(*posts)[i].Created = createdFormatted
	//		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
	//		if post.Parent != 0 {
	//			args = append(args, post.Parent, post.Author, post.Message, thread.Forum, thread.ID, created)
	//		} else {
	//			args = append(args, nil, post.Author, post.Message, thread.Forum, thread.ID, created)
	//		}
	//	}
	//	query = query[:len(query)-1]
	//	query += " RETURNING id;"
	//	resultRows, err := threadStore.db.Query(query, args...)
	//	if err != nil {
	//		return errors.ErrParentPostNotExist
	//	}
	//	defer resultRows.Close()
	//
	//	for i := 0; resultRows.Next(); i++ {
	//		var id int64
	//		if err = resultRows.Scan(&id); err != nil {
	//			return err
	//		}
	//		(*posts)[i].ID = id
	//	}
	//}

	return
}

func (threadStore *ThreadStore) GetPostsTree(threadID int64, limit, since int, desc bool) (posts *[]models.Post, err error) {
	var rows *pgx.Rows
	query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts " +
		"WHERE thread = $1 "

	if since != -1 && desc {
		query += " AND path < "
	} else if since != -1 && !desc {
		query += " AND path > "
	} else if since != -1 {
		query += " AND path > "
	}
	if since != -1 {
		query += fmt.Sprintf(` (SELECT path FROM posts WHERE id = %d) `, since)
	}
	if desc {
		query += " ORDER BY path DESC"
	} else if !desc {
		query += " ORDER BY path ASC, id"
	} else {
		query += " ORDER BY path, id"
	}
	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0);", limit)

	rows, err = threadStore.db.Query(query, threadID)
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
					WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path ASC, id ASC;`, threadID, limit)
		} else {
			rows, err = threadStore.db.Query(`
					SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL ORDER BY id ASC LIMIT $2)
					ORDER BY path ASC, id ASC;`, threadID, limit)
		}
	} else {
		if desc {
			rows, err = threadStore.db.Query(`
					SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL AND path[1] <
					(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3)
					ORDER BY path[1] DESC, path ASC, id ASC;`, threadID, since, limit)
		} else {
			rows, err = threadStore.db.Query(`
					SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL AND path[1] >
					(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id ASC LIMIT $3) 
					ORDER BY path ASC, id ASC;`, threadID, since, limit)
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
	query := "SELECT id, COALESCE(parent, 0), author, message, is_edited, forum, thread, created FROM posts WHERE thread = $1 "

	if since != -1 && desc {
		query += " AND id < $2"
	} else if since != -1 && !desc {
		query += " AND id > $2"
	} else if since != -1 {
		query += " AND id > $2"
	}
	if desc {
		query += " ORDER BY created DESC, id DESC"
	} else if !desc {
		query += " ORDER BY created ASC, id"
	} else {
		query += " ORDER BY created, id"
	}
	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0);", limit)

	if since == -1 {
		rows, err = threadStore.db.Query(query, threadID)
	} else {
		rows, err = threadStore.db.Query(query, threadID, since)
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
