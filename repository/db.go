package repository

import (
	"fmt"
	"task1/entity"

	"github.com/go-pg/pg/v10"
)

type Config struct {
	User     string
	Password string
	Database string
}

type PostgresDB struct {
	config Config
	db     *pg.DB
}

func New(config Config) *PostgresDB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "123456",
		Database: "postgres",
	})

	return &PostgresDB{config: config, db: db}
}

func (d *PostgresDB) Register(u entity.User)error {
	/*

	user1 := &entity.User{
	//	ID:       9,
		Email:    "sdfwewer@gmail",
		Password: "345123",
	}
	_, err := d.db.Model(user1).Insert()
	if err != nil {
		panic(err)
	}

	_, err = d.db.Model(&entity.User{
	//	ID:       7,
		Email:    "wefewew@gmail",
		Password: "12334",
	}).Insert()
	if err != nil {
		panic(err)
	}
	*/

	_, err := d.db.Model(&u).Insert()
	fmt.Println(err)
	return nil
	 
	
}
