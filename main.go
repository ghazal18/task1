package main

import (
	"task1/Repository"
	"task1/controller"
	"task1/delivery"
	"task1/repository/acl"
	"task1/service/projectservice"
	"task1/service/userservice"
	"task1/validator/user_validator"
	"time"
)

func main() {

	repo := repository.New()
	userControl := controller.New(createConfig())
	acl := acl.Service{
		Repo: repo,
		Podb: repo,
	}

	usersvc := userservice.New(repo, userControl, acl)
	projectsvc := projectservice.New(repo, userControl, acl)


	uservalid := uservalidator.New()

	serve := delivery.New(usersvc, userControl, acl, uservalid,projectsvc)
	serve.Serve()

}

func createConfig() controller.Config {
	return controller.Config{
		SignKey:              "jwt_secret",
		AccessExpirationTime: time.Hour * 24,
	}

}
