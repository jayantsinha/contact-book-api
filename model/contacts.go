package model

type Contact struct {
	ContactID int    `db:"contact_id"`
	AccoutnID int    `db:"account_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	PhoneNum  string `db:"phone_num"`
}
