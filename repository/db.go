package repository

import (
	"fmt"
	"task1/entity"

	"github.com/go-pg/pg/v10"
)

type PostgresDB struct {
	db *pg.DB
}

func New() *PostgresDB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "123456",
		Database: "postgres",
	})

	return &PostgresDB{db: db}
}

func (d *PostgresDB) Register(u entity.User) error {
	_, err := d.db.Model(&u).Insert()
	fmt.Println(err)
	return nil

}

func (d *PostgresDB) GetUser(u entity.User) {

	userQuery := `select * from users where email=? and password_hash=?;`
	userEmail := u.Email
	userPass := u.Password
	userID := u.ID
	res, err := d.db.Query(&u, userQuery, userEmail, userPass, userID)
	fmt.Println(res.RowsReturned(), res.Model(), u, err)
}
