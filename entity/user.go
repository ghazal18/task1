package entity

type User struct {
	ID       int    `pg:"id"`
	Email    string `pg:"email"`
	Password string `pg:"password_hash"`
}
