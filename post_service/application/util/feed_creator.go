package util

import (
	"context"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
	psa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/profile_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
)

type FeedCreator struct {
	store             domain.PostStore
	connectionAdapter csa.IConnectionServiceAdapter
	profileAdapter    psa.IProfileServiceAdapter
}

func NewFeedCreator(
	store domain.PostStore,
	connectionAdapter csa.IConnectionServiceAdapter,
	profileAdapter psa.IProfileServiceAdapter) *FeedCreator {
	return &FeedCreator{
		store:             store,
		connectionAdapter: connectionAdapter,
		profileAdapter:    profileAdapter,
	}
}

func (creator *FeedCreator) CreateFeedForUser(ctx context.Context, userId primitive.ObjectID) []*domain.Post {
	var posts []*domain.Post

	if userId != primitive.NilObjectID {
		connections := creator.connectionAdapter.GetAllUserConnections(ctx, userId)
		if len(connections) > 0 {
			posts = append(posts, creator.store.GetAllPostsFromIds(connections)...)
		}
	}

	publicProfilesIds := creator.profileAdapter.GetAllPublicProfilesIds(ctx)
	if len(publicProfilesIds) > 0 {
		posts = append(posts, creator.store.GetAllPostsFromIds(publicProfilesIds)...)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreationTime.After(posts[j].CreationTime)
	})

	return posts
}
