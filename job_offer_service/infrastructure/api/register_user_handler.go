package api

import (
	"context"
	"fmt"
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/create_order"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/application"
)

type RegisterUserCommandHandler struct {
	jobOfferService   *application.JobOfferService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewRegisterUserCommandHandler(jobOfferService *application.JobOfferService, publisher saga.Publisher, subscriber saga.Subscriber) (*RegisterUserCommandHandler, error) {
	o := &RegisterUserCommandHandler{
		jobOfferService:   jobOfferService,
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
	case events.CreateNodeInJobOfferBase:
		actionResult, err := handler.jobOfferService.CreateUser(context.TODO(), command.Order.Id)
		if err != nil {
			fmt.Println("NodeInJobOfferBaseNotCreated " + err.Error())
			reply.Type = events.NodeInJobOfferBaseNotCreated
		} else {
			if actionResult.Status == 200 {
				fmt.Println("SUPERRR NodeInJobOfferBaseCreated")
				reply.Type = events.NodeInJobOfferBaseCreated
			} else {
				fmt.Println("LOSEEEE NodeInJobOfferBaseNotCreated")
				reply.Type = events.NodeInJobOfferBaseNotCreated
			}
		}
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
