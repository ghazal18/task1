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

func (d *PostgresDB) Register(u entity.User) (entity.User, error) {
	_, err := d.db.Model(&u).Insert()
	fmt.Println(err)
	return u, nil

}

func (d *PostgresDB) GetUser(u entity.User) (entity.User, bool, error) {

	userQuery := `select * from users where email=? and password_hash=?;`
	userEmail := u.Email
	userPass := u.Password
	userID := u.ID
	res, err := d.db.Query(&u, userQuery, userEmail, userPass, userID)
	fmt.Println(res.RowsReturned(), res.Model(), u, err)
	return u, true, nil
}

func (d *PostgresDB) CreateProject(p entity.Project) (entity.Project, bool, error) {

	projectQuery := `Insert into projects(owner_id,name,company,description,social_links) values (?,?,?,?,?);`
	projectOwner := p.OwnerID
	projectName := p.Name
	projectCompany := p.Company
	projectDesc := p.Description
	projectSocial := p.SocialLinks

	res, err := d.db.Query(&p, projectQuery, projectOwner, projectName, projectCompany, projectDesc, projectSocial)
	fmt.Println(res.RowsReturned(), res.Model(), p, err)
	fmt.Println("this is p ", p)
	return p, true, nil
}

func (d *PostgresDB) AllProject(uID int) (p []entity.Project, b bool, e error) {

	projectQuery := `SELECT DISTINCT p.*
FROM projects p
LEFT JOIN project_members pm ON p.id = pm.project_id
WHERE p.owner_id = ? OR pm.user_id = ?;`
	userId := uID

	d.db.Query(&p, projectQuery, userId, userId)

	return p, true, nil
}
