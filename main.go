package main

import (
	//"fmt"
	"task1/Repository"
	"task1/delivery"
	// "task/pg_test"
)

func main() {

	cfg := repository.Config{
		User:     "postgres",
		Password: "123456",
		Database: "postgres",
	}
	Postdb := repository.New(cfg)
	Postdb.Register()
	delivery.Serve2()

}
