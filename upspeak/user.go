package upspeak

import (
	"errors"
	"log"
)

/////////////////
// Create user //
/////////////////

var createUserQuery = `
insert into users (username, password, email_primary, created_at, is_active, display_name, is_verified)
    values($1, $2, $3, now(), $4, $5, $6)
    returning id
`

// CreateUser returns a User given a username
func CreateUser(user *User) (int64, error) {
	var userID int64

	passwordHash, err := GeneratePasswordHash(user.Password.String)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	err = db.QueryRow(createUserQuery, &user.Username, &passwordHash, &user.EmailPrimary, &user.IsActive, &user.DisplayName, &user.IsVerified).Scan(&userID)

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

// GetUser gets user metadata
func GetUser(username *string) (User, error) {
	var u User
	err := db.QueryRow(getUserQuery, &username).Scan(
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

// UpdateUser updates user metadata
func UpdateUser(user *User) (int64, error) {
	var userID int64
	err := db.QueryRow(updateUserQuery, &user.DisplayName, &user.UserID).Scan(&userID)

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

// UpdateUserEmail updates user's email
func UpdateUserEmail(email string, userID int64) (int64, error) {
	res, err := db.Exec(updateUserEmailQuery, &email, &userID)

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
func UpdateUserPassword(userID int64, oldPassword string, newPassword string) error {
	// Get old password
	var oldPasswordHash string
	err := db.QueryRow(getUserPasswordQuery, &userID).Scan(&userID, &oldPasswordHash)
	if err != nil {
		return err
	}

	err = VerifyPasswordHash(oldPassword, oldPasswordHash)
	if err != nil {
		return err
	}

	// Generate hash for new password
	newPasswordHash, err := GeneratePasswordHash(newPassword)
	if err != nil {
		return err
	}

	res, err := db.Exec(updateUserPasswordQuery, &newPasswordHash, &userID)
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

func ToggleUserActivationStatus(userID int64, activationStatus bool) (int64, error) {
	err := db.QueryRow(toggleUserActivationStatusQuery, &activationStatus, &userID).Scan(&userID)

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

func VerifyUser(userID int64, verificationStatus bool) (int64, error) {
	err := db.QueryRow(verifyUserQuery, &verificationStatus, &userID).Scan(&userID)

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

func AuthenticateUser() {

}
