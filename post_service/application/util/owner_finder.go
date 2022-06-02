package util

import (
	"context"
	psa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/profile_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OwnerFinder struct {
	profileServiceAdapter psa.IProfileServiceAdapter
}

func NewOwnerFinder(adapter psa.IProfileServiceAdapter) *OwnerFinder {
	return &OwnerFinder{
		profileServiceAdapter: adapter,
	}
}

func (finder *OwnerFinder) GetOwnerFinderFunction(ctx context.Context) func(id primitive.ObjectID) *domain.Owner {
	m := finder.getOwnerMap(ctx)
	return func(id primitive.ObjectID) *domain.Owner {
		return m[id]
	}
}

func (finder *OwnerFinder) getOwnerMap(ctx context.Context) map[primitive.ObjectID]*domain.Owner {
	owners := finder.profileServiceAdapter.GetAllProfiles(ctx)
	m := make(map[primitive.ObjectID]*domain.Owner)
	for _, owner := range owners {
		m[owner.UserId] = owner
	}
	return m
}
