package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"task1/entity"
)

type Repository interface {
	Register(u entity.User) (entity.User, error)
	GetUser(u entity.User) (entity.User, bool, error)
	CreateProject(p entity.Project) (entity.Project, error)
	AllProject(uID int) (p []entity.Project, e error)
	AllOtherProject(uID int) (p []entity.Project, e error)
	FindProjectByID(pID int) (p entity.Project, e error)
	DeleteProjectByID(pID int) (p entity.Project, b bool, e error)
	UpdateProjectByID(p entity.Project) (entity.Project, bool, error)
	JoinProjectByID(pID, uID string) (bool, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
}

type ACLGenerator interface {
	CanViewProject(userID, projectID int) bool
	CanEditProject(userID, projectID int) bool
}

type Service struct {
	repo Repository
	auth AuthGenerator
	acl  ACLGenerator
}

func New(repo Repository, auth AuthGenerator, acl ACLGenerator) Service {
	return Service{repo: repo, auth: auth, acl: acl}
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID    int    `json: "id"`
	Email string `json:"email"`
}

func (s Service) SignUp(req SignUpRequest) (*SignUpResponse, error) {
	var resp SignUpResponse

	user := entity.User{
		Email:    req.Email,
		Password: getMD5Hash(req.Password),
	}

	user, err := s.repo.Register(user)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	resp = SignUpResponse{
		ID:    user.ID,
		Email: user.Email,
	}
	return &resp, nil

}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
type LoginResponse struct {
	User   LoginUserResponse `json:"user"`
	Tokens Token             `json:"tokens"`
}

func (s Service) Login(req LoginRequest) (*LoginResponse, error) {

	user := entity.User{
		Email:    req.Email,
		Password: getMD5Hash(req.Password),
	}

	user, exist, err := s.repo.GetUser(user)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}
	if !exist {
		return nil, fmt.Errorf("username or password isn't correct")
	}
	if user.Password != getMD5Hash(req.Password) {

		return nil, fmt.Errorf("username or password isn't correct")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	return &LoginResponse{
		User: LoginUserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		Tokens: Token{
			AccessToken: accessToken,
		},
	}, nil

}


func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
