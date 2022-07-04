package notification_service_adapter

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters"
)

type NotificationServiceAdapter struct {
	address string
}

func NewNotificationServiceAdapter(address string) *NotificationServiceAdapter {
	return &NotificationServiceAdapter{address: address}
}

func (conn *NotificationServiceAdapter) InsertNotification(ctx context.Context, request *pb.InsertNotificationRequest) {
	span := tracer.StartSpanFromContext(ctx, "InsertNotification")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	connectionClient := adapters.NewNotificationClient(conn.address)
	connectionClient.InsertNotification(ctx2, request)
}
