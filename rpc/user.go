package rpc

import (
	"errors"
	"log"
	"net/http"

	"github.com/applait/upspeak/models"
)

type UserService struct{}

//////////////
// Get user //
//////////////

// UserGetArgs defines the input parameters for user.Get
type UserGetArgs struct {
	Username string `json:"username"`
}

// UserGetReply defines the result data type for user.Get
type UserGetReply struct {
	User models.User `json:"user"`
}

// Get fetches the contents of a single user given a username
func (n *UserService) Get(r *http.Request, args *UserGetArgs, reply *UserGetReply) error {
	user, err := models.GetUser(&args.Username)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.User = user
	return nil
}

/////////////////
// Create user //
/////////////////

// UserCreateArgs defines the input parameters for user.Create
type UserCreateArgs struct {
	Username     models.NullString `json:"username"`
	Password     models.NullString `json:"password"`
	EmailPrimary models.NullString `json:"email_primary"`
	DisplayName  models.NullString `json:"display_name"`
}

// UserCreateReply defines the result data type for user.Create
type UserCreateReply struct {
	UserID int64 `json:"user_id"`
}

// Create creates a single user
func (n *UserService) Create(r *http.Request, args *UserCreateArgs, reply *UserCreateReply) error {
	data := &models.User{
		Username:     args.Username,
		Password:     args.Password,
		EmailPrimary: args.EmailPrimary,
		IsActive:     true,
		DisplayName:  args.DisplayName,
		IsVerified:   false,
	}
	userID, err := models.CreateUser(data)
	if err != nil {
		log.Println(err)
	}

	// TODO: Trigger email verification

	// Generate response
	reply.UserID = userID
	return nil
}

//////////////////////////
// Update user metadata //
//////////////////////////

// UserUpdateArgs defines the input parameters for user.Update
type UserUpdateArgs struct {
	UserID      models.NullInt64  `json:"user_id"`
	DisplayName models.NullString `json:"display_name"`
}

// UserUpdateReply defines the result data type for user.Update
type UserUpdateReply struct {
	UserID int64 `json:"user_id"`
}

// Update upadtes a single user's metadata
func (n *UserService) Update(r *http.Request, args *UserUpdateArgs, reply *UserUpdateReply) error {
	data := &models.User{
		DisplayName: args.DisplayName,
		UserID:      args.UserID,
	}
	userID, err := models.UpdateUser(data)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.UserID = userID
	return nil
}

///////////////////////
// Update user email //
///////////////////////

// UserUpdateEmail defines the input parameters for user.UpdateEmail
type UserUpdateEmailArgs struct {
	UserID       int64  `json:"user_id"`
	EmailPrimary string `json:"email_primary"`
}

// UserUpdateEmailReply defines the result data type for user.UpdateEmail
type UserUpdateEmailReply struct {
	UserID int64 `json:"user_id"`
}

// Update updates a single user's email
func (n *UserService) UpdateEmail(r *http.Request, args *UserUpdateEmailArgs, reply *UserUpdateEmailReply) error {
	userID, err := models.UpdateUserEmail(args.EmailPrimary, args.UserID)
	if err != nil {
		log.Println(err)
	}

	// TODO: Trigger email verification

	// Generate response
	reply.UserID = userID
	return nil
}

//////////////////////////
// Update user password //
//////////////////////////

// UserUpdatePasswordArgs defines the input parameters for user.UpdatePassword
type UserUpdatePasswordArgs struct {
	UserID      int64  `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UserUpdatePasswordReply defines the result data type for user.UpdatePassword
type UserUpdatePasswordReply struct {
	UserID int64 `json:"user_id"`
}

// UpdatePassword upadtes a single user's password
func (n *UserService) UpdatePassword(r *http.Request, args *UserUpdatePasswordArgs, reply *UserUpdatePasswordReply) error {
	err := models.UpdateUserPassword(args.UserID, args.OldPassword, args.NewPassword)
	if err != nil {
		return errors.New("Error updating password")
	}

	// Generate response
	reply.UserID = args.UserID
	return nil
}

/////////////////
// Verify user //
/////////////////

// UserVerify defines the input parameters for user.Verify
type UserVerifyArgs struct {
	UserID            int64  `json:"user_id"`
	Email             string `json:"email"`
	VerificationToken string `json:"verification_token"`
}

// UserVerify defines the result data type for user.Verify
type UserVerifyReply struct {
	UserID int64 `json:"user_id"`
}

// Verify updates a user's verification status
func (n *UserService) Verify(r *http.Request, args *UserVerifyArgs, reply *UserVerifyReply) error {
	_, err := ParseToken(args.VerificationToken, Conf.JWT.Secret, Conf.JWT.Audience)
	if err != nil {
		log.Println(err)
		return nil
	}

	userID, err := models.VerifyUser(args.UserID, true)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.UserID = userID
	return nil
}
