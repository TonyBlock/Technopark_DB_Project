package models

const (
	Dislike = iota - 1
	Like    = iota
)

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int32  `json:"voice"`
}
