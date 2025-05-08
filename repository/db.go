package repository

import (
	"fmt"
	"task1/entity"
	errormsg "task1/error"

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
	if err != nil {
		pgErr, ok := err.(pg.Error)
		pgErr.IntegrityViolation()
		if ok && pgErr.IntegrityViolation() {
			return u, fmt.Errorf("%w", errormsg.ErrUserAlreadyExists)
		}
		return u, fmt.Errorf("%w", errormsg.ErrInternal)
	}
	return u, nil

}

func (d *PostgresDB) GetUser(u entity.User) (entity.User, bool, error) {

	userQuery := `select * from users where email=? and password_hash=?;`
	userEmail := u.Email
	userPass := u.Password
	_, err := d.DB.Query(&u, userQuery, userEmail, userPass)
	if err != nil {
		return u, false, fmt.Errorf("%w", errormsg.ErrInternal)
	}
	if u.ID == 0 {
		return u, false, nil
	}
	return u, true, nil

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
		pgErr, ok := err.(pg.Error)
		pgErr.IntegrityViolation()
		if ok && pgErr.IntegrityViolation() {
			return p, fmt.Errorf("%w", errormsg.ErrProjectAlreadyExists)
		}
		return p, fmt.Errorf("%w", errormsg.ErrInternal)
	}
	return p, nil

}

func (d *PostgresDB) AllProject(uID int) (p []entity.Project, e error) {

	projectQuery := `SELECT DISTINCT p.*
    FROM projects p
    LEFT JOIN project_members pm ON p.id = pm.project_id
    WHERE p.owner_id = ? OR pm.user_id = ?;`
	userId := uID

	_, err := d.DB.Query(&p, projectQuery, userId, userId)
	if err != nil {
		return p, fmt.Errorf("%w", errormsg.ErrInternal)
	}

	return p, nil
}

func (d *PostgresDB) AllOtherProject(uID int) (p []entity.Project, e error) {

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
		return p, fmt.Errorf("%w", errormsg.ErrInternal)
	}

	return p, nil
}

func (d *PostgresDB) FindProjectByID(pID int) (p entity.Project, e error) {

	projectQuery := `SELECT * FROM projects WHERE id = ?;`

	_, err := d.DB.Query(&p, projectQuery, pID)
	if err != nil {
		return p, fmt.Errorf("%w", errormsg.ErrInternal)
	}

	return p, nil
}

func (d *PostgresDB) DeleteProjectByID(pID int) (p entity.Project, b bool, e error) {

	projectQuery := `DELETE FROM projects WHERE id = ?;`

	res, err := d.DB.Query(&p, projectQuery, pID)
	if res.RowsAffected() == 0 {
		return p, false, nil
	}

	if err != nil {
		return p, true, fmt.Errorf("%w", errormsg.ErrInternal)
	}

	return p, true, nil
}

func (d *PostgresDB) UpdateProjectByID(p entity.Project) (entity.Project, bool, error) {

	columns := []string{}

	project := entity.Project{
		ID:          p.ID,
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
	if len(p.SocialLinks) != 0 {
		columns = append(columns, "social_links")
	}


	res, err := d.DB.Model(&project).Column(columns...).WherePK().Returning("*").Update()
	if err != nil {
		return p, false, fmt.Errorf("%w", errormsg.ErrInternal)
	}
	if res.RowsAffected() == 0 {
		return p, false, nil
	}
	return p, true, nil
}

func (d *PostgresDB) JoinProjectByID(pID, uID string) (bool, error) {
	var pm entity.ProjectMembers

	projectMemberQuery := ` insert into project_members(project_id,user_id) values (?,?)`
	_, err := d.DB.Query(&pm, projectMemberQuery, pID, uID)
	if err != nil {
		return false, fmt.Errorf("%w", errormsg.ErrInternal)
	}

	return true, nil

}

func (d *PostgresDB) IsOwner(userID, projectID int) (b bool) {
	var p entity.Project

	
	query := `SELECT * FROM projects WHERE id = ? ;`
	d.DB.Query(&p, query, projectID)

	
	fmt.Println("ownerID == userID", p.OwnerID, userID, p.OwnerID == userID, projectID)
	return p.OwnerID == userID
}

func (d *PostgresDB) IsMember(userID, projectID int) (b bool) {
	var pm entity.ProjectMembers

	
	query := `SELECT * FROM project_members 
	 	WHERE project_id = ? AND user_id = ?;`

	d.DB.Query(&pm, query, projectID, userID)

	
	if pm.ID != 0 {
		return true
	}
	return false
}
