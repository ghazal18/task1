package userservice

import (
	"fmt"
	"task1/entity"
)

type Repository interface {
	Register(u entity.User) error
	GetUser(u entity.User) (entity.User, bool, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
}

type Service struct {
	repo Repository
	auth AuthGenerator
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(repo Repository, auth AuthGenerator) Service {
	return Service{repo: repo, auth: auth}
}

func (s Service) Register(req LoginRequest) {

	user := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	s.repo.Register(user)

}

func (s Service) Login(req LoginRequest) {

	user := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}
	fmt.Println("userservice.service.go ", user)

	user, _, _ = s.repo.GetUser(user)
	fmt.Println("userservice.service.go ", user)
	str, _ := s.auth.CreateAccessToken(user)
	fmt.Println(str)

}
