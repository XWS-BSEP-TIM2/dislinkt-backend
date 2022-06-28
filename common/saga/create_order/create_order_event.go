package create_order

import "time"

type UserDetails struct {
	Id          string
	Name        string
	Surname     string
	Username    string
	Email       string
	Birthday    time.Time
	Gender      string
	PhoneNumber string
	IsPrivate   bool
}

type RegisterUserCommandType int8

const (
	CreateUserCredentials RegisterUserCommandType = iota
	RollbackCreateUserCredentials
	CreateUserProfile
	RollbackCreateUserProfile
	CreateNodeInConnectionBase
	RollbackCreateNodeInConnectionBase
	UserRegistered
	UnknownCommand
)

type RegisterUserCommand struct {
	Order UserDetails
	Type  RegisterUserCommandType
}

type RegisterUserReplyType int8

const (
	UserCredentialsCreated RegisterUserReplyType = iota
	UserProfileCreated
	UserProfileNotCreated
	NodeInConnectionBaseCreated
	NodeInConnectionBaseNotCreated
	DoneRollbackOfProfile
	UserNotRegistered
	UnknownReply
)

type RegisterUserReply struct {
	Order UserDetails
	Type  RegisterUserReplyType
}
