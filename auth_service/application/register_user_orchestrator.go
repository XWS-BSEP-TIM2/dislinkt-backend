package application

import (
	"fmt"
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
	// Happy
	case events.UserCredentialsCreated:
		fmt.Println("NEXT CreateUserProfile")
		return events.CreateUserProfile
	case events.UserProfileCreated:
		fmt.Println("NEXT CreateNodeInConnectionBase")
		return events.CreateNodeInConnectionBase
	case events.NodeInConnectionBaseCreated:
		fmt.Println("NEXT CreateNodeInJobOfferBase")
		return events.CreateNodeInJobOfferBase
	case events.NodeInJobOfferBaseCreated:
		fmt.Println("SUPER ODLICNO SAGA JE USPESNO PROSLA :)")
		return events.UnknownCommand

	// fails
	case events.NodeInJobOfferBaseNotCreated:
		fmt.Println("NEXT RollbackCreateNodeInConnectionBase")
		return events.RollbackCreateNodeInConnectionBase

	case events.NodeInConnectionBaseNotCreated:
		fmt.Println("NEXT RollbackCreateUserProfile")
		return events.RollbackCreateUserProfile

	case events.UserProfileNotCreated:
		fmt.Println("NEXT RollbackCreateUserCredentials")
		return events.RollbackCreateUserCredentials

	case events.DoneRollBackInConnection:
		fmt.Println("NEXT RollbackCreateUserProfile")
		return events.RollbackCreateUserProfile

	case events.DoneRollbackOfProfile:
		fmt.Println("NEXT RollbackCreateUserCredentials")
		return events.RollbackCreateUserCredentials

	default:
		return events.UnknownCommand
	}
}
