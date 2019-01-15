package model

type Account struct {
	AccountID int    `db:"account_id"`
	Email     string `db:"email"`
	Name      string `db:"name"`
	Password  string `db:"password"`
}
