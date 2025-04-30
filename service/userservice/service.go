package userservice

import "task1/entity"

type Repository interface {
	Register(u entity.User) error
}

type Service struct {
	repo Repository
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Register(req LoginRequest) {

	user := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	s.repo.Register(user)

}
