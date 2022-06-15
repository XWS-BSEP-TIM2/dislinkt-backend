package domain

import (
	"fmt"
	"time"
)

type Log struct {
	Id                  string
	Time                time.Time
	ServiceName         string
	ServiceFunctionName string
	LogType             string
	UserID              string
	IpAddress           string
	Description         string
}

func NewLog(serviceName, serviceFunctionName, logType, userID, ipAddress, description string) *Log {
	return &Log{
		Id:                  "",
		Time:                time.Now().UTC(),
		ServiceName:         serviceName,
		ServiceFunctionName: serviceFunctionName,
		LogType:             logType,
		UserID:              userID,
		IpAddress:           ipAddress,
		Description:         description,
	}
}

func (l *Log) ToString() string {
	return fmt.Sprintf("%s | %s | %s | %s | %s | %s | %s | %s", l.Id, l.Time, l.ServiceName, l.ServiceFunctionName, l.LogType, l.UserID, l.IpAddress, l.Description)
}
