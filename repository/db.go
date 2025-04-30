package repository 

import (
	"github.com/go-pg/pg/v10"
	"task1/entity"

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

//	defer db.Close()
	return &PostgresDB{config: config, db: db}
}

func (d *PostgresDB) Register() {
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
}
