package domain

type UserConn struct {
	UserID    string
	IsPrivate bool
}

type UserConnDetail struct {
	MyUserID  string
	UserID    string
	IsPrivate bool
	Relation  string
	MsgID     string
}
