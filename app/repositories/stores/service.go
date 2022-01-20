package stores

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"

	"github.com/jackc/pgx"
)

type ServiceStore struct {
	db *pgx.ConnPool
}

func CreateServiceRepository(db *pgx.ConnPool) repositories.ServiceRepository {
	return &ServiceStore{db: db}
}

func (serviceStore *ServiceStore) Clear() (err error) {
	_, err = serviceStore.db.Exec("TRUNCATE TABLE forums, posts, threads, user_forum, users, votes CASCADE;")
	return
}

func (serviceStore *ServiceStore) GetStatus() (status *models.Status, err error) {
	status = &models.Status{}
	err = serviceStore.db.QueryRow("(SELECT count(*) FROM users) AS users,"+
		"SELECT (SELECT count(*) FROM forums) AS forums, "+
		"(SELECT count(*) FROM threads) AS threads, "+
		"(SELECT count(*) FROM posts) AS posts;").
		Scan(&status.User, &status.Forum, &status.Thread, &status.Post)
	return
}
