package profile_service_adapter

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProfileServiceAdapter interface {
	GetAllProfiles(ctx context.Context) []*domain.Owner
	GetSingleProfile(ctx context.Context, profileId primitive.ObjectID) *domain.Owner
}
