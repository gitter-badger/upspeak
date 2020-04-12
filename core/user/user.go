package user

import (
	"errors"
	"log"

	"github.com/upspeak/upspeak/core"
)

/////////////////
// Create user //
/////////////////

var createUserQuery = `
insert into users (username, password, email_primary, created_at, is_active, display_name, is_verified)
    values($1, $2, $3, now(), $4, $5, $6)
    returning id
`

// Create creates a new user and returns the User ID
func Create(user *core.User) (int64, error) {
	var userID int64

	passwordHash, err := core.GeneratePasswordHash(user.Password.String)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	err = core.DB.QueryRow(createUserQuery, &user.Username, &passwordHash, &user.EmailPrimary, &user.IsActive, &user.DisplayName, &user.IsVerified).Scan(&userID)

	if err != nil {
		log.Println(err)
		return userID, err
	}
	return userID, nil
}

//////////////
// Get user //
//////////////

var getUserQuery = `
select
    id, username, email_primary, display_name, created_at, is_verified, is_active
from users
    where username = $1;
`

// Get gets user metadata
func Get(username *string) (core.User, error) {
	var u core.User
	err := core.DB.QueryRow(getUserQuery, &username).Scan(
		&u.UserID, &u.Username, &u.EmailPrimary, &u.DisplayName, &u.CreatedAt,
		&u.IsVerified, &u.IsActive,
	)

	if err != nil {
		log.Println(err)
		return u, err
	}

	return u, nil
}

/////////////////
// Update user //
/////////////////

var updateUserQuery = `
update users
    set display_name = $1
where id = $2
    returning id;
`

// Update updates user metadata
func Update(user *core.User) (int64, error) {
	var userID int64
	err := core.DB.QueryRow(updateUserQuery, &user.DisplayName, &user.UserID).Scan(&userID)

	if err != nil {
		return userID, err
	}
	return userID, nil
}

///////////////////////
// Update user email //
///////////////////////

var updateUserEmailQuery = `
update users
    set primary_email = $1
where id = $2;
`

// UpdateEmail updates user's email
func UpdateEmail(email string, userID int64) (int64, error) {
	res, err := core.DB.Exec(updateUserEmailQuery, &email, &userID)

	if err != nil {
		return userID, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return userID, err
	}
	if count == 0 {
		return userID, errors.New("Zero rows affected")
	}

	return userID, nil
}

//////////////////////////
// Update user password //
//////////////////////////

var updateUserPasswordQuery = `
update users
    set password = $1
where users.id = $2
    returning id;
`

var getUserPasswordQuery = `
select
    id, password
from users
    where id = $1;
`

// UpdatePassword updates the user's password
func UpdatePassword(userID int64, oldPassword string, newPassword string) error {
	// Get old password
	var oldPasswordHash string
	err := core.DB.QueryRow(getUserPasswordQuery, &userID).Scan(&userID, &oldPasswordHash)
	if err != nil {
		return err
	}

	err = core.VerifyPasswordHash(oldPassword, oldPasswordHash)
	if err != nil {
		return err
	}

	// Generate hash for new password
	newPasswordHash, err := core.GeneratePasswordHash(newPassword)
	if err != nil {
		return err
	}

	res, err := core.DB.Exec(updateUserPasswordQuery, &newPasswordHash, &userID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Zero rows affected")
	}

	return nil
}

///////////////////////////////////
// Toggle user activation status //
///////////////////////////////////

var toggleUserActivationStatusQuery = `
update users
    set is_active = $1
where users.id = $2
    returning id;
`

// Activate changes user activation status
func Activate(userID int64, activationStatus bool) (int64, error) {
	err := core.DB.QueryRow(toggleUserActivationStatusQuery, &activationStatus, &userID).Scan(&userID)

	if err != nil {
		return userID, err

	}
	return userID, nil
}

/////////////////
// Verify user //
/////////////////

var verifyUserQuery = `
update users
    set is_verified = $1
where users.id = $2
    returning id;
`

// Verify updates a user's verification status
func Verify(userID int64, verificationStatus bool) (int64, error) {
	err := core.DB.QueryRow(verifyUserQuery, &verificationStatus, &userID).Scan(&userID)

	if err != nil {
		return userID, err

	}
	return userID, nil

}

///////////////////////
// Authenticate user //
///////////////////////

var authenticateUserQuery = `

`

// Authenticate attempts to authenticate a user. TODO.
func Authenticate() {

}
