package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/converter"
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/create_order"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUserCommandHandler struct {
	profileService    *application.ProfileService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewRegisterUserCommandHandler(profileService *application.ProfileService, publisher saga.Publisher, subscriber saga.Subscriber) (*RegisterUserCommandHandler, error) {
	o := &RegisterUserCommandHandler{
		profileService:    profileService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *RegisterUserCommandHandler) handle(command *events.RegisterUserCommand) {
	id, err := primitive.ObjectIDFromHex(command.Order.Id)
	if err != nil {
		return
	}
	order := &domain.Profile{Id: id,
		IsPrivate:   command.Order.IsPrivate,
		PhoneNumber: command.Order.PhoneNumber,
		Email:       command.Order.Email,
		Gender:      enums.ToEnumGender(command.Order.Gender),
		Username:    command.Order.Username,
		Surname:     command.Order.Surname,
		BirthDate:   command.Order.Birthday,
		Name:        command.Order.Name,
		IsTwoFactor: false,
	}

	reply := events.RegisterUserReply{Order: command.Order}

	switch command.Type {

	case events.CreateUserProfile:
		fmt.Println("Profil service:Kreiranje profila")
		handler.profileService.Insert(context.TODO(), order)
		reply.Type = events.UserProfileCreated
	case events.RollbackCreateUserProfile:
		fmt.Println("Profil service:Rollback profila")
		handler.profileService.DeleteById(context.TODO(), converter.GetObjectId(command.Order.Id))
		reply.Type = events.DoneRollbackOfProfile
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
