package stores

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/repositories"

	"github.com/jackc/pgx"
)

type VoteStore struct {
	db *pgx.ConnPool
}

func CreateVoteRepository(db *pgx.ConnPool) repositories.VoteRepository {
	return &VoteStore{db: db}
}

func (voteStore *VoteStore) Vote(threadID int64, vote *models.Vote) (err error) {
	_, err = voteStore.db.Exec("INSERT INTO votes (nickname, thread, voice) "+
		"VALUES ($1, $2, $3) ON CONFLICT (nickname, thread) DO UPDATE SET voice = EXCLUDED.voice;",
		vote.Nickname, threadID, vote.Voice)
	return
}
