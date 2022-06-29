package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/create_order"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type CreateRegisterCommandHandler struct {
	registerService   *application.AuthService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
	tracer            opentracing.Tracer
	closer            io.Closer
}

func NewRegisterUserCommandHandler(registerService *application.AuthService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateRegisterCommandHandler, error) {
	o := &CreateRegisterCommandHandler{
		registerService:   registerService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateRegisterCommandHandler) handle(command *events.RegisterUserCommand) {
	id, err := primitive.ObjectIDFromHex(command.Order.Id)
	if err != nil {
		return
	}
	order := &domain.User{Id: id}

	reply := events.RegisterUserReply{Order: command.Order}

	switch command.Type {

	case events.RollbackCreateUserProfile:
		fmt.Println("Auth service:Rollback kredencijala")
		err := handler.registerService.DeleteById(context.TODO(), order.Id)
		if err != nil {
			return
		}
		reply.Type = events.UserNotRegistered
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
