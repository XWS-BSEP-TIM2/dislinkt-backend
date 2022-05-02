package domain

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Skill struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Type enums.SkillType    `bson:"type"`
}

func NewSkill(name string, skillType enums.SkillType) Skill {
	return Skill{Name: name, Type: skillType}
}
