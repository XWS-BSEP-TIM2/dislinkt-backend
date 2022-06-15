package domain

import (
	"time"
)

type Chat struct {
	Id                  string
	Time                time.Time
	ServiceName         string
	ServiceFunctionName string
	UserID              string
	IpAddress           string
	Description         string
}
