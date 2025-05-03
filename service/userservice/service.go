package userservice

import (
	"fmt"
	"task1/entity"
)

type Repository interface {
	Register(u entity.User) (entity.User, error)
	GetUser(u entity.User) (entity.User, bool, error)
	CreateProject(p entity.Project) (entity.Project, bool, error)
	AllProject(uID int) (p []entity.Project, b bool, e error)
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

	user, _ = s.repo.Register(user)
	resp = RegisterResponse{
		ID:    0,
		Email: user.Email,
	}
	return resp, nil

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

type NewProjectRequest struct {
	OwnerID     int    `json:"owner_id"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Description string `json:"description"`
	SocialLinks string `json:"social_links"`
}

type NewProjectResponse struct {
	Name string `json:"name"`
}

func (s Service) NewProject(req NewProjectRequest) (NewProjectResponse, error) {
	pr := entity.Project{
		Name:        req.Name,
		Company:     req.Company,
		OwnerID:     req.OwnerID,
		Description: req.Description,
		SocialLinks: req.SocialLinks,
	}
	s.repo.CreateProject(pr)
	return NewProjectResponse{Name: "created"}, nil
}

type AllProjectRequest struct {
	ID int `json:"id"`
}

type AllProjectResponse struct {
	OwnerID     int    `json:"owner_id"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Description string `json:"description"`
	SocialLinks string `json:"social_links"`
}

func (s Service) GetAllProject(id int) ([]entity.Project, error) {

	project, _, _ := s.repo.AllProject(id)

	return project, nil

}
