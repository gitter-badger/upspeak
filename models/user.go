package models

type User struct {
	ID           int64
	Username     string
	Password     string
	EmailPrimary string
	CreatedAt    string
	IsVerified   bool
	IsActive     bool
	DisplayName  string
}
