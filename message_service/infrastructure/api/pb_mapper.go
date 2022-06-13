package api

import (
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapChat(chat *domain.Chat) *pb.Chat {
	pbChat := &pb.Chat{
		MsgID:         chat.Id.Hex(),
		UserIDa:       chat.UserIDa,
		UserIDb:       chat.UserIDb,
		UserASeenDate: &timestamppb.Timestamp{Seconds: chat.UserASeenDate.Unix()},
		UserBSeenDate: &timestamppb.Timestamp{Seconds: chat.UserBSeenDate.Unix()},
	}

	for _, msg := range chat.Messages {
		pbChat.Messages = append(pbChat.Messages, &pb.Message{AuthorUserID: msg.AuthorUserID, Text: msg.Text, Date: &timestamppb.Timestamp{Seconds: msg.Date.Unix()}})
	}

	return pbChat
}
