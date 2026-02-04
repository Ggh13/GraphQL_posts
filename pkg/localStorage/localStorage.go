package localstorage

import (
	"qraphQL_posts/api/graph/model"
)

type Storage struct {
	Users    []model.User
	Posts    []model.Post
	Comments []model.Comment
}

func NewLocalStorage() Storage {
	return Storage{
		Users:    []model.User{},
		Posts:    []model.Post{},
		Comments: []model.Comment{},
	}
}
