package update_skills

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
)

type UpdateSkillsDetails struct {
	UserID               string
	OldSkills            []domain.Skill
	NewSkillsForJobOffer []string
}

type UpdateSkillsCommandType int8

const (
	CommandUpdateSkillsInProfile UpdateSkillsCommandType = iota
	CommandRollbackUpdateSkillsInProfile
	CommandUpdateSkillsInJobOffer
	UnknownCommand
)

type UpdateSkillsCommand struct {
	UpdateSkillsDetail UpdateSkillsDetails
	Type               UpdateSkillsCommandType
}

type UpdateSkillsReplyType int8

const (
	ReplySkillsUpdatedInProfile UpdateSkillsReplyType = iota
	ReplyRollbackSkillsUpdatedInProfile
	ReplySkillsUpdatedInJobOffer
	ReplySkillsNOTUpdatedInJobOffer
	UnknownReply
)

type UpdateSkillsReply struct {
	UpdateSkillsDetail UpdateSkillsDetails
	Type               UpdateSkillsReplyType
}
