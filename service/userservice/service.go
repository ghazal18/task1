package userservice

import (
	"fmt"
	"task1/entity"
)

type Repository interface {
	Register(u entity.User) (entity.User, error)
	GetUser(u entity.User) (entity.User, bool, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
}

type Service struct {
	repo Repository
	auth AuthGenerator
}

func New(repo Repository, auth AuthGenerator) Service {
	return Service{repo: repo, auth: auth}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    int    `json: "id"`
	Email string `json:"email"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	var resp RegisterResponse

	user := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user,_= s.repo.Register(user)
	resp = RegisterResponse{
		ID: 0,
		Email: user.Email,

	}
	return resp,nil


}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s Service) Login(req LoginRequest) {

	user := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user, _, _ = s.repo.GetUser(user)
	token, _ := s.auth.CreateAccessToken(user)
	fmt.Println(token)

}
