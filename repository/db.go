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

func (d *PostgresDB) GetUser(u entity.User) {

	userQuery := `select email,password_hash from users where email=? and password_hash=?;`
	userEmail := u.Email
	userPass := u.Password
	res, err := d.db.Query(&u, userQuery, userEmail, userPass)
	fmt.Println(res.RowsReturned(),res.Model(), u,err)

	// err := d.db.Model(&u).from

	// 	Where("users.email = ?", u.Email).
	// 	Select()

	// fmt.Println(u, err)
}
