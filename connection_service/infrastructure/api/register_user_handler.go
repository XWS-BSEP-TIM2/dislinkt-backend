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
		actionResult, err := handler.connectionService.Register(command.Order.Id, command.Order.IsPrivate)
		if err != nil {
			fmt.Println("Kreiranje usera nije uspelo" + err.Error())
			reply.Type = events.NodeInConnectionBaseNotCreated
		} else if actionResult != nil && actionResult.Status == 201 {
			reply.Type = events.NodeInConnectionBaseCreated
		} else {
			reply.Type = events.NodeInConnectionBaseNotCreated
		}
		break
	case events.RollbackCreateNodeInConnectionBase:
		actionResult, err := handler.connectionService.DeleteUser(command.Order.Id)
		if err != nil {
			fmt.Println("Rollback brisanje usera nije uspelo " + err.Error())
			reply.Type = events.UnknownReply
		} else if actionResult != nil && actionResult.Status == 200 {
			fmt.Println("Rollback brisanje usera uspesno")
			reply.Type = events.DoneRollBackInConnection
		} else {
			reply.Type = events.UnknownReply
		}
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
