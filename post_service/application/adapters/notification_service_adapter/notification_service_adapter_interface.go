package notification_service_adapter

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
)

type INotificationServiceAdapter interface {
	InsertNotification(ctx context.Context, request *pb.InsertNotificationRequest)
}
