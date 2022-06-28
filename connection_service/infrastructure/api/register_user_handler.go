package api

import (
	"fmt"
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/create_order"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/application"
)

type RegisterUserCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewRegisterUserCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) (*RegisterUserCommandHandler, error) {
	o := &RegisterUserCommandHandler{
		connectionService: connectionService,
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

	reply := events.RegisterUserReply{Order: command.Order}

	switch command.Type {

	case events.CreateNodeInConnectionBase:
		err, _ := handler.connectionService.Register(command.Order.Id, command.Order.IsPrivate)
		if err.Status != 201 {
			fmt.Println("Connection service:Ne mogu da kreiram node")
			reply.Type = events.NodeInConnectionBaseNotCreated
		} else {
			reply.Type = events.NodeInConnectionBaseCreated
		}
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
