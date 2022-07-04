package application

import (
	"fmt"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/update_skills"
)

type UpdateSkillsOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewUpdateSkillsOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*UpdateSkillsOrchestrator, error) {
	o := &UpdateSkillsOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *UpdateSkillsOrchestrator) Start(updateSkillsDetails events.UpdateSkillsDetails) error {
	event := &events.UpdateSkillsCommand{
		Type:               events.CommandUpdateSkillsInJobOffer,
		UpdateSkillsDetail: updateSkillsDetails,
	}

	return o.commandPublisher.Publish(event)
}

func (o *UpdateSkillsOrchestrator) handle(reply *events.UpdateSkillsReply) {
	command := events.UpdateSkillsCommand{UpdateSkillsDetail: reply.UpdateSkillsDetail}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *UpdateSkillsOrchestrator) nextCommandType(reply events.UpdateSkillsReplyType) events.UpdateSkillsCommandType {
	switch reply {
	// Happy
	case events.ReplySkillsUpdatedInProfile:
		return events.CommandUpdateSkillsInJobOffer
	case events.ReplySkillsUpdatedInJobOffer:
		fmt.Println("USPESNO AZURIRANI SKILLLS U SAGAAA")
		return events.UnknownCommand

	//fail
	case events.ReplySkillsNOTUpdatedInJobOffer:
		return events.CommandRollbackUpdateSkillsInProfile

	default:
		return events.UnknownCommand
	}
}
