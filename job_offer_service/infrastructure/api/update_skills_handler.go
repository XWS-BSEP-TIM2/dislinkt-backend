package api

import (
	"context"
	"fmt"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/update_skills"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/application"
)

type UpdateSkillsCommandHandler struct {
	jobOfferService   *application.JobOfferService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateSkillsCommandHandler(profileService *application.JobOfferService, publisher saga.Publisher, subscriber saga.Subscriber) (*UpdateSkillsCommandHandler, error) {
	o := &UpdateSkillsCommandHandler{
		jobOfferService:   profileService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *UpdateSkillsCommandHandler) handle(command *events.UpdateSkillsCommand) {

	reply := events.UpdateSkillsReply{UpdateSkillsDetail: command.UpdateSkillsDetail}

	switch command.Type {
	case events.CommandUpdateSkillsInJobOffer:

		actionResult, err := handler.jobOfferService.UpdateUserSkills(context.TODO(), command.UpdateSkillsDetail.UserID, command.UpdateSkillsDetail.NewSkillsForJobOffer)
		if err != nil {
			fmt.Println("error " + err.Error())
			reply.Type = events.ReplySkillsNOTUpdatedInJobOffer
		}
		if actionResult != nil && actionResult.Status == 200 {
			fmt.Println("Uspesno azurirani skills u useru SAGA poziv")
			reply.Type = events.ReplySkillsUpdatedInJobOffer
		} else {
			if actionResult != nil {
				fmt.Println("Doslo je do greske prilikom azuriranja skills - SAGA " + actionResult.Msg)
			} else {
				fmt.Println("Doslo je do greske prilikom azuriranja skills - SAGA ")
			}
			reply.Type = events.ReplySkillsNOTUpdatedInJobOffer
		}

		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
