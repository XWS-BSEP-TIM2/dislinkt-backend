package api

import (
	"context"
	"fmt"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/update_skills"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateSkillsCommandHandler struct {
	profileService    *application.ProfileService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateSkillsCommandHandler(profileService *application.ProfileService, publisher saga.Publisher, subscriber saga.Subscriber) (*UpdateSkillsCommandHandler, error) {
	o := &UpdateSkillsCommandHandler{
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

func (handler *UpdateSkillsCommandHandler) handle(command *events.UpdateSkillsCommand) {

	reply := events.UpdateSkillsReply{UpdateSkillsDetail: command.UpdateSkillsDetail}

	switch command.Type {
	case events.CommandRollbackUpdateSkillsInProfile:
		fmt.Println("Profile service CommandRollbackUpdateSkillsInProfile")
		id, err := primitive.ObjectIDFromHex(command.UpdateSkillsDetail.UserID)
		if err != nil {
			fmt.Println("Error convert id")
			return
		}

		profile, errGetProfile := handler.profileService.Get(context.TODO(), id)
		if errGetProfile != nil {
			fmt.Println("Error user is not found id" + command.UpdateSkillsDetail.UserID)
			return
		}

		profile.Skills = command.UpdateSkillsDetail.OldSkills
		err1 := handler.profileService.UpdateSkills(context.TODO(), profile)
		if err1 != nil {
			fmt.Println("Rollback, updateSkills with old skills")
			return
		}
		reply.Type = events.ReplyRollbackSkillsUpdatedInProfile
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
