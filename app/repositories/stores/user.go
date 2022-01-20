package stores

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"

	"github.com/jackc/pgx"
)

type UserStore struct {
	db *pgx.ConnPool
}

func CreateUserRepository(db *pgx.ConnPool) repositories.UserRepository {
	return &UserStore{db: db}
}

func (userStore *UserStore) Create(user *models.User) (err error) {
	_, err = userStore.db.Exec("INSERT INTO users VALUES ($1, $2, $3, $4);",
		user.Nickname, user.Fullname, user.About, user.Email)
	return
}

func (userStore *UserStore) Update(user *models.User) (err error) {
	return userStore.db.QueryRow("UPDATE users SET"+
		"fullname = COALESCE(NULLIF(TRIM($1), ''), fullname), "+
		"about = COALESCE(NULLIF(TRIM($2), ''), about), "+
		"email = COALESCE(NULLIF(TRIM($3), ''), email) "+
		"WHERE nickname = $4 RETURNING fullname, about, email;",
		user.Fullname, user.About, user.Email, user.Nickname).Scan(&user.Fullname, &user.About, &user.Email)
}

func (userStore *UserStore) GetByNickname(nickname string) (user *models.User, err error) {
	user = &models.User{}
	err = userStore.db.QueryRow("SELECT nickname, fullname, about, email FROM users "+
		"WHERE nickname = $1;", nickname).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return
}

func (userStore *UserStore) GetAllMatchedUsers(user *models.User) (users *[]models.User, err error) {
	var usersSlice []models.User

	resultRows, err := userStore.db.Query("SELECT nickname, fullname, about, email FROM users "+
		"WHERE nickname = $1 OR email = $2;", user.Nickname, user.Email)
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
