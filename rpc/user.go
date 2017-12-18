package rpc

import "net/http"

type UserService struct{}

////////////////////
// User construct //
////////////////////
type UserReply struct {
}

/////////////////
// Create user //
/////////////////
type UserCreateArgs struct {
	// TBD
}

// Create creates a user
func (u *UserService) Create(r *http.Request, args *UserCreateArgs, reply *UserReply) error {
	return nil
}

/////////////////
// Verify user //
/////////////////
type UserVerifyArgs struct {
	UserId      int64  `json:"user_id"`
	VerifyToken string `json:"verify_token"`
}

func (u *UserService) Verify(r *http.Request, args *UserVerifyArgs, reply *UserReply) error {
	return nil
}
