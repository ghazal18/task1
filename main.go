package main

import (
	repository "task1/Repository"
	"task1/controller"
	"task1/delivery"
	"task1/service/userservice"
	"time"
)

func main() {

	repo := repository.New()
	userControl := controller.New(createConfig())
	usersvc := userservice.New(repo, userControl)

	serve := delivery.New(usersvc, userControl)
	serve.Serve()

}

func createConfig() controller.Config {
	return controller.Config{
		SignKey:              "jwt_secret",
		AccessExpirationTime: time.Hour * 24,
	}

}
