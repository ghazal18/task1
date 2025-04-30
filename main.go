package main

import (
	"task1/Repository"
	"task1/delivery"
	"task1/service/userservice"
)

func main() {

	cfg := repository.Config{
		User:     "postgres",
		Password: "123456",
		Database: "postgres",
	}
	repo := repository.New(cfg)
	serve := delivery.New(userservice.New(repo))

	
	serve.Serve()

	

}
