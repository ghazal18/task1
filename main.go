package main

import (
	"task1/Repository"
	"task1/delivery"
	"task1/service/userservice"
)

func main() {

	repo := repository.New()
	usersvc := userservice.New(repo)
	serve := delivery.New(usersvc)
	serve.Serve()
	

}
