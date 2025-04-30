package main

import (
	//"fmt"
	"task1/Repository"
//	"task/pg_test"
)

func main() {
//	fmt.Println("this shit is working")
//	pg_test.ExampleDB_Model()

	cfg := repository.Config{
		User:     "postgres",
		Password: "123456",
		Database: "postgres",
	}
	Postdb := repository.New(cfg)
	Postdb.Register()

}
