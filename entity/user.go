package entity

type User struct {
//	ID       uint `pg:"id"`
	Email    string `pg:"email"`
	Password string `pg:"password_hash"`
}
