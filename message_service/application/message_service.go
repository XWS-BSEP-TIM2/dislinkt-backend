package application

import (
	"context"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/application/adapters"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/startup/config"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sort"
	"time"
)

type MessageService struct {
	store            persistence.MessageStore
	ConnectionClient connectionService.ConnectionServiceClient
	ProfileClient    profileService.ProfileServiceClient
}

func NewMessageService(store persistence.MessageStore, c *config.Config) *MessageService {
	return &MessageService{
		store:            store,
		ConnectionClient: adapters.NewConnectionClient(fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort)),
		ProfileClient:    adapters.NewProfileClient(fmt.Sprintf("%s:%s", c.ProfileHost, c.ProfilePort)),
	}
}

func (service *MessageService) GetMyContacts(ctx context.Context, request *pb.GetMyContactsRequest) (*pb.MyContactsResponse, error) {
	userID := request.UserID
	contactResponse, err := service.ConnectionClient.GetMyContacts(ctx, &connectionService.GetMyContactsRequest{UserID: userID})
	if err != nil {
		return nil, err
	}

	var chatsPreview []*pb.ChatPreview

	for i, contact := range contactResponse.Contacts {
		fmt.Println("Contact: ", i, contact.UserID, contact.MsgID)
		chat, errGetChat := service.store.GetChat(ctx, contact.MsgID)
		if errGetChat != nil {
			fmt.Println("Error: ", errGetChat.Error())
			continue
		}

		profile, errGetProfile := service.ProfileClient.Get(ctx, &profileService.GetRequest{Id: contact.UserID})
		if errGetProfile != nil {
			fmt.Println("Error", errGetProfile.Error())
			continue
		}

		fullNameUser := profile.Profile.Name + " " + profile.Profile.Surname
		lastMessage := domain.Message{AuthorUserID: "", Text: "", Date: chat.GetSeenDateByUserID(userID)}
		if len(chat.Messages) > 0 {
			lastMessage = chat.Messages[len(chat.Messages)-1]
		}
		date := &timestamppb.Timestamp{Seconds: lastMessage.Date.Unix()}
		numOfNewMessages := 0
		if lastMessage.AuthorUserID == contact.UserID {
			// zadnja poruka koja se nalazi u mesingeru nije nasa nego od osobe sa kojom delimo chat
			// treba da pogledamo koliko ima poruka koje nismo videli
			var lastSeen time.Time
			if userID == chat.UserIDa {
				lastSeen = chat.UserASeenDate
			} else if userID == chat.UserIDb {
				lastSeen = chat.UserBSeenDate
			}
			for j := len(chat.Messages) - 1; j >= 0; j-- {
				// ako je poruka od naseg sagovornika, i ako je ta poruka poslata kasnije nego sto smo mi zadnji put seenovali
				if chat.Messages[j].AuthorUserID == contact.UserID && chat.Messages[j].Date.After(lastSeen) {
					numOfNewMessages++
				} else {
					break
				}
			}
		}

		chatsPreview = append(chatsPreview, &pb.ChatPreview{UserID: contact.UserID, MsgID: contact.MsgID, FullNameUser: fullNameUser,
			LastMessage: &pb.Message{AuthorUserID: lastMessage.AuthorUserID, Text: lastMessage.Text, Date: date}, NumOfNewMessages: int32(numOfNewMessages)})
	}

	sort.Slice(chatsPreview, func(i, j int) bool {
		return chatsPreview[i].LastMessage.Date.Seconds > chatsPreview[j].LastMessage.Date.Seconds
	})

	return &pb.MyContactsResponse{Chats: chatsPreview}, nil
}

func (service *MessageService) GetChat(ctx context.Context, request *pb.GetChatRequest) (*pb.ChatResponse, error) {
	chat, err := service.store.GetChat(ctx, request.MsgID)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	if !chat.HaveUserID(request.UserID) {
		fmt.Println("Error: the user has no access")
		return nil, nil
	}

	pbChat := mapChat(chat, request.UserID)

	profileResponseA, errGetProfileA := service.ProfileClient.Get(ctx, &profileService.GetRequest{Id: pbChat.UserIDa})
	if errGetProfileA != nil {
		fmt.Println("Error", errGetProfileA.Error())
	} else {
		pbChat.FullNameUserA = profileResponseA.Profile.Name + " " + profileResponseA.Profile.Surname
	}

	profileResponseB, errGetProfileB := service.ProfileClient.Get(ctx, &profileService.GetRequest{Id: pbChat.UserIDb})
	if errGetProfileB != nil {
		fmt.Println("Error", errGetProfileB.Error())
	} else {
		pbChat.FullNameUserB = profileResponseB.Profile.Name + " " + profileResponseB.Profile.Surname
	}

	return &pb.ChatResponse{Chat: pbChat}, nil
}

func mapChat(chat *domain.Chat, myUserID string) *pb.Chat {

	var pbChat *pb.Chat

	if chat.UserIDa == myUserID {
		pbChat = &pb.Chat{
			MsgID:         chat.Id.Hex(),
			UserIDa:       chat.UserIDa,
			UserIDb:       chat.UserIDb,
			FullNameUserA: "",
			FullNameUserB: "",
			UserASeenDate: &timestamppb.Timestamp{Seconds: chat.UserASeenDate.Unix()},
			UserBSeenDate: &timestamppb.Timestamp{Seconds: chat.UserBSeenDate.Unix()},
			Messages:      []*pb.Message{},
		}
	} else if chat.UserIDb == myUserID {
		pbChat = &pb.Chat{
			MsgID:         chat.Id.Hex(),
			UserIDa:       chat.UserIDb,
			UserIDb:       chat.UserIDa,
			FullNameUserA: "",
			FullNameUserB: "",
			UserASeenDate: &timestamppb.Timestamp{Seconds: chat.UserBSeenDate.Unix()},
			UserBSeenDate: &timestamppb.Timestamp{Seconds: chat.UserASeenDate.Unix()},
			Messages:      []*pb.Message{},
		}
	}

	for _, msg := range chat.Messages {
		pbChat.Messages = append(pbChat.Messages, &pb.Message{AuthorUserID: msg.AuthorUserID, Text: msg.Text, Date: &timestamppb.Timestamp{Seconds: msg.Date.Unix()}})
	}

	return pbChat
}
