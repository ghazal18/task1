package acl

import repository "task1/Repository"

type Repo interface {
	IsOwner(userID, projectID int) bool
	IsMember(userID, projectID int) bool
}

type Service struct {
	Repo Repo
	Podb *repository.PostgresDB
}

func (s Service) CanViewProject(userID, projectID int) bool {
	return s.Repo.IsOwner(userID, projectID) || s.Repo.IsMember(userID, projectID)
}

func (s Service) CanEditProject(userID, projectID int) bool {
	return s.Repo.IsOwner(userID, projectID)
}
