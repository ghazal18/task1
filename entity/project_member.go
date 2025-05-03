package entity
type ProjectMembers struct {
	ID          int    `pg:"id"`
	ProjectID     int    `pg:"project_id"`
	UserID     int    `pg:"user_id"`
	
}
