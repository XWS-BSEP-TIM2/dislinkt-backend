package logging_service_adapter

import "context"

type ILoggingServiceAdapter interface {
	Log(ctx context.Context, logType, serviceFunctionName, userID, description string)
}
