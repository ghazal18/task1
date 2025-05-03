package main

import (
	"task1/Repository"
	"task1/controller"
	"task1/delivery"
	"task1/repository/acl"
	"task1/service/userservice"
	"time"
)

func main() {

	repo := repository.New()
	userControl := controller.New(createConfig())
	acl := acl.Service{
		Repo :repo,
		Podb: repo,
	}

	usersvc := userservice.New(repo, userControl,acl)

	serve := delivery.New(usersvc, userControl,acl)
	serve.Serve()

}

func createConfig() controller.Config {
	return controller.Config{
		SignKey:              "jwt_secret",
		AccessExpirationTime: time.Hour * 24,
	}

}
