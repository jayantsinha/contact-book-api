package model

// Contact is used to bind query result
type Contact struct {
	ContactID int    `db:"contact_id"`
	AccountID int    `db:"account_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

// ContactResp is used in endpoint responses
type ContactResp struct {
	ContactID int    `json:"contact_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}