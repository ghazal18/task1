package entity

type User struct {
	Email    string `pg:"email"`
	Password string `pg:"password_hash"`
}
