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
	AllProject(uID int) (p []entity.Project, b bool, e error) 
	AllOtherProject(uID int) (p []entity.Project, b bool, e error)
	FindProjectByID(pID int) (p entity.Project, b bool, e error)
	DeleteProjectByID(pID int) (p entity.Project, b bool, e error)
	UpdateProjectByID(pID string, p entity.Project)
	JoinProjectByID(pID, uID string)
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

func (s Service) SignUp(req SignUpRequest) (SignUpResponse, error) {
	var resp SignUpResponse

	user := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user, _ = s.repo.Register(user)
	resp = SignUpResponse{
		ID:    user.ID,
		Email: user.Email,
	}
	return resp, nil

}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User   entity.User `json:"user"`
	Tokens Token       `json:"tokens"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {

	user := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user, exist, err := s.repo.GetUser(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("something unexpected happend")
	}
	if !exist {
		return LoginResponse{}, fmt.Errorf("username or password isn't correct")
	}
	//if user.Password != getMD5Hash(req.Password) {
	if user.Password != req.Password {
		return LoginResponse{}, fmt.Errorf("username or password isn't correct")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return LoginResponse{
		User: entity.User{
			ID:       user.ID,
			Email:    user.Email,
			Password: user.Password,
		},
		Tokens: Token{
			AccessToken: accessToken,
		},
	}, nil

}

type NewProjectRequest struct {
	ID          int    `pg:"id"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Description string `json:"description"`
	SocialLinks string `json:"social_links"`
}

type NewProjectResponse struct {
	ID          int    `pg:"id"`
	OwnerID     int    `pg:"owner_id"`
	Name        string `pg:"name"`
	Company     string `pg:"company"`
	Description string `pg:"description"`
	SocialLinks string `pg:"social_links"`
}

func (s Service) NewProject(req NewProjectRequest, userID int) (NewProjectResponse, error) {
	pr := entity.Project{
		Name:        req.Name,
		Company:     req.Company,
		OwnerID:     req.ID,
		Description: req.Description,
		SocialLinks: req.SocialLinks,
	}
	project, err := s.repo.CreateProject(pr)
	if err != nil {
		fmt.Errorf("something unexpected happend")
	}
	return NewProjectResponse{
		ID:          project.ID,
		OwnerID:     project.OwnerID,
		Name:        project.Name,
		Company:     project.Company,
		Description: project.Description,
		SocialLinks: project.SocialLinks,
	}, nil
}

type AllProjectRequest struct {
	ID int
}

func (s Service) GetAllProject(request AllProjectRequest) ([]entity.Project, bool, error) {

	project, exist, err := s.repo.AllProject(request.ID)

	if err != nil {
		fmt.Errorf("something unexpected happend")
		return project, exist, err
	}

	return project, exist, nil

}

type AllOtherProjectRequest struct {
	ID int
}

func (s Service) GetAllOthersProject(request AllOtherProjectRequest) ([]entity.Project, bool, error) {

	project, exist, err := s.repo.AllOtherProject(request.ID)
	if err != nil {
		fmt.Errorf("something unexpected happend")
		return project, exist, err
	}

	return project, exist, nil

}

func (s Service) GetProjectByID(id int) (entity.Project, error) {

	project, _, _ := s.repo.FindProjectByID(id)

	return project, nil

}

func (s Service) DeleteProjectByID(id int) (entity.Project, error) {

	project, _, _ := s.repo.DeleteProjectByID(id)

	return project, nil

}

func (s Service) JoinProjectByID(pID, uID string) {

	s.repo.JoinProjectByID(pID, uID)
}
func (s Service) UpdateProjectByID(pID string, p entity.Project) {

	s.repo.UpdateProjectByID(pID, p)

}

func (s Service) PutProjectByID(id string, update map[string]interface{}) error {
	return nil
}

type UpdateProjectRequest struct {
	OwnerID     int    `json:"owner_id"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Description string `json:"description"`
	SocialLinks string `json:"social_links"`
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
