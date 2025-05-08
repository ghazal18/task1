package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
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
		return nil, fmt.Errorf("%w",err)
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
		return nil,fmt.Errorf("unexpected error: %w", err)
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

type NewProjectRequest struct {
	ID          int               `pg:"id"`
	Name        string            `json:"name"`
	Company     string            `json:"company"`
	Description string            `json:"description"`
	SocialLinks map[string]string `json:"social_links"`
}

type NewProjectResponse struct {
	ID          int               `pg:"id"`
	OwnerID     int               `pg:"owner_id"`
	Name        string            `pg:"name"`
	Company     string            `pg:"company"`
	Description string            `pg:"description"`
	SocialLinks map[string]string `pg:"social_links"`
}

func (s Service) NewProject(req NewProjectRequest, userID int) (*NewProjectResponse, error) {
	
	pr := entity.Project{
		Name:        req.Name,
		Company:     req.Company,
		OwnerID:     req.ID,
		Description: req.Description,
		SocialLinks: req.SocialLinks,
	}
	project, err := s.repo.CreateProject(pr)
	if err != nil {
		return nil , fmt.Errorf("%w",err)
	}

	return &NewProjectResponse{
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

func (s Service) GetAllProject(request AllProjectRequest) (*[]entity.Project, error) {

	project, err := s.repo.AllProject(request.ID)
	if len(project)==0 {
		return nil, fmt.Errorf("you don't have project")
	}

	if err != nil {
		return nil, fmt.Errorf("%w",err)
	}

	return &project, nil

}

type AllOtherProjectRequest struct {
	ID int
}

func (s Service) GetAllOthersProject(request AllOtherProjectRequest) (*[]entity.Project, error) {

	project, err := s.repo.AllOtherProject(request.ID)
	if len(project)==0 {
		return nil, fmt.Errorf("there is no project")
	}
	if err != nil {
		return nil, fmt.Errorf("%w",err)
	}

	return &project, nil

}

type GetProjectByIDRequest struct {
	UserID    int
	ProjectID int
}

func (s Service) GetProjectByID(request GetProjectByIDRequest) (entity.Project, error) {

	project, err := s.repo.FindProjectByID(request.ProjectID)
	if err != nil {
		
		return project, fmt.Errorf("something unexpected happend")
	}

	return project, nil

}

func (s Service) DeleteProjectByID(id int) (entity.Project, bool, error) {

	project, affected, err := s.repo.DeleteProjectByID(id)
	if err != nil {
		return project, affected, err
	}

	return project, affected, nil

}

type JoinProjectByIDRequest struct {
	ProjectID string
	UserID    string
}

func (s Service) JoinProjectByID(req JoinProjectByIDRequest) (bool, error) {

	done, err := s.repo.JoinProjectByID(req.ProjectID, req.UserID)
	if err != nil {
		fmt.Println(err)
	}
	return done, err
}

func (s Service) UpdateProjectByID(pID string, p PutProjectByIDRequest) (PutProjectByIDRespons, bool, error) {
	id, err := strconv.Atoi(pID)
	if err != nil {
		fmt.Errorf("There is a problem in casting")
	}
	project := entity.Project{
		ID:          id,
		Name:        p.Name,
		Company:     p.Company,
		Description: p.Description,
		SocialLinks: p.SocialLinks,
	}

	project, ok, err := s.repo.UpdateProjectByID(project)

	resp := PutProjectByIDRespons{
		Name:        p.Name,
		Company:     p.Company,
		Description: p.Description,
		SocialLinks: p.SocialLinks,
	}
	return resp, ok, err

}

type PutProjectByIDRequest struct {
	ID          int               `json:"id"`
	OwnerID     int               `json:"owner_id"`
	Name        string            `json:"name"`
	Company     string            `json:"company"`
	Description string            `json:"description"`
	SocialLinks map[string]string `json:"social_links"`
}

type PutProjectByIDRespons struct {
	Name        string            `json:"name"`
	Company     string            `json:"company"`
	Description string            `json:"description"`
	SocialLinks map[string]string `json:"social_links"`
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
