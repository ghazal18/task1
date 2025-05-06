package repository

import (
	"fmt"
	"strconv"
	"task1/entity"

	"github.com/go-pg/pg/v10"
)

type PostgresDB struct {
	DB *pg.DB
}

func New() *PostgresDB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "123456",
		Database: "postgres",
	})

	return &PostgresDB{DB: db}
}

func (d *PostgresDB) Register(u entity.User) (entity.User, error) {
	_, err := d.DB.Model(&u).Insert()
	fmt.Println(err)
	return u, nil

}

func (d *PostgresDB) GetUser(u entity.User) (entity.User, bool, error) {

	userQuery := `select * from users where email=? and password_hash=?;`
	userEmail := u.Email
	userPass := u.Password
	_, err := d.DB.Query(&u, userQuery, userEmail, userPass)
	if err != nil {
		fmt.Errorf("Something unexpected happend")
	}
	if u.ID == 0 {
		return u, false, nil
	} else {
		return u, true, nil
	}
}

func (d *PostgresDB) CreateProject(p entity.Project) (entity.Project, error) {

	projectQuery := `Insert into projects(owner_id,name,company,description,social_links) values (?,?,?,?,?) RETURNING id;`
	projectOwner := p.OwnerID
	projectName := p.Name
	projectCompany := p.Company
	projectDesc := p.Description
	projectSocial := p.SocialLinks

	_, err := d.DB.Query(&p, projectQuery, projectOwner, projectName, projectCompany, projectDesc, projectSocial)
	if err != nil {
		fmt.Errorf("something unexpected happend")
	}

	return p, nil
}

func (d *PostgresDB) AllProject(uID int) (p []entity.Project, b bool, e error) {

	projectQuery := `SELECT DISTINCT p.*
    FROM projects p
    LEFT JOIN project_members pm ON p.id = pm.project_id
    WHERE p.owner_id = ? OR pm.user_id = ?;`
	userId := uID

	_, err := d.DB.Query(&p, projectQuery, userId, userId)
	if err != nil {
		fmt.Errorf("something unexpected happend")
		return p, false, err
	}
	if len(p) == 0 {
		return p, false, nil
	}

	return p, true, nil
}

func (d *PostgresDB) AllOtherProject(uID int) (p []entity.Project, b bool, e error) {

	projectQuery := `SELECT p.*
	FROM projects p
    WHERE p.owner_id != ?
    AND p.id NOT IN (
    SELECT project_id
    FROM project_members
    WHERE user_id = ?
);
`
	userId := uID

	_, err := d.DB.Query(&p, projectQuery, userId, userId)
	if err != nil {
		fmt.Errorf("something unexpected happend")
		return p, false, err
	}
	if len(p) == 0 {
		return p, false, nil
	}

	return p, true, nil
}

func (d *PostgresDB) FindProjectByID(pID int) (p entity.Project, e error) {

	projectQuery := `SELECT * FROM projects WHERE id = ?;`

	_, err := d.DB.Query(&p, projectQuery, pID)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (d *PostgresDB) DeleteProjectByID(pID int) (p entity.Project, b bool, e error) {

	projectQuery := `DELETE FROM projects WHERE id = ?;`

	res , err :=d.DB.Query(&p, projectQuery, pID)
	if res.RowsAffected() == 0 {
		return p , false , nil
	}
	
	if err!= nil {
		return p, true, err
	}

	return p, true, nil
}

func (d *PostgresDB) UpdateProjectByID(pID string, p entity.Project) {

	columns := []string{}

	id, _ := strconv.Atoi(pID)

	project := entity.Project{
		ID:          id,
		Name:        p.Name,
		Company:     p.Company,
		Description: p.Description,
		SocialLinks: p.SocialLinks,
	}

	if p.Name != "" {
		columns = append(columns, "name")
	}
	if p.Company != "" {
		columns = append(columns, "company")
	}
	if p.Description != "" {
		columns = append(columns, "description")
	}
	if p.SocialLinks != "" {
		columns = append(columns, "social_links")
	}
	fmt.Println("p,columns", p, columns)
	d.DB.Model(&project).Column(columns...).WherePK().Returning("*").Update()
	fmt.Println(project)

}

func (d *PostgresDB) JoinProjectByID(pID, uID string) {
	var pm entity.ProjectMembers

	projectMemberQuery := ` insert into project_members(project_id,user_id) values (?,?)`

	d.DB.Query(&pm, projectMemberQuery, pID, uID)

}

func (d *PostgresDB) IsOwner(userID, projectID int) (b bool) {
	var p entity.Project

	//query := `SELECT owner_id FROM projects WHERE id = ? ;`
	query := `SELECT * FROM projects WHERE id = ? ;`
	d.DB.Query(&p, query, projectID)

	// return err == nil && ownerID == userID
	fmt.Println("ownerID == userID", p.OwnerID, userID, p.OwnerID == userID, projectID)
	return p.OwnerID == userID
}

func (d *PostgresDB) IsMember(userID, projectID int) (b bool) {
	var pm entity.ProjectMembers

	// query := `SELECT EXISTS (
	// 	SELECT 1 FROM project_members
	// 	WHERE project_id = ? AND user_id = ?
	//     );`
	query := `SELECT * FROM project_members 
	 	WHERE project_id = ? AND user_id = ?;`

	d.DB.Query(&pm, query, projectID, userID)

	// return err == nil && exists
	if pm.ID != 0 {
		fmt.Println(pm.ID)
		return true
	}
	return false
}
