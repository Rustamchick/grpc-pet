package models

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	PassHash []byte `db:"password_hash"`
}

type App struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Token string `db:"token"`
}
