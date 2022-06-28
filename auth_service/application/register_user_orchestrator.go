package application

import (
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/create_order"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
)

type RegisterUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewRegisterUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*RegisterUserOrchestrator, error) {
	o := &RegisterUserOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *RegisterUserOrchestrator) Start(userDetails events.UserDetails) error {
	event := &events.RegisterUserCommand{
		Type:  events.CreateUserProfile,
		Order: userDetails,
	}

	return o.commandPublisher.Publish(event)
}

func (o *RegisterUserOrchestrator) handle(reply *events.RegisterUserReply) {
	command := events.RegisterUserCommand{Order: reply.Order}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *RegisterUserOrchestrator) nextCommandType(reply events.RegisterUserReplyType) events.RegisterUserCommandType {
	switch reply {
	case events.UserCredentialsCreated:
		return events.CreateUserProfile
	case events.UserProfileCreated:
		return events.CreateNodeInConnectionBase
	case events.UserProfileNotCreated:
		return events.RollbackCreateUserCredentials
	case events.NodeInConnectionBaseNotCreated:
		return events.RollbackCreateUserProfile
	case events.DoneRollbackOfProfile:
		return events.RollbackCreateUserCredentials
	default:
		return events.UnknownCommand
	}
}
