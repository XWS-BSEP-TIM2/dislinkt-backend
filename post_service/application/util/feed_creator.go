package util

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "CreateFeedForUser")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	var posts []*domain.Post

	var connections []*primitive.ObjectID
	if userId != primitive.NilObjectID {
		connections = creator.connectionAdapter.GetAllUserConnections(ctx2, userId)
		if len(connections) > 0 {
			posts = append(posts, creator.store.GetAllPostsFromIds(connections)...)
		}
	}

	publicProfilesIds := creator.profileAdapter.GetAllPublicProfilesIds(ctx2)
	var filtered []*primitive.ObjectID
	for _, ppId := range publicProfilesIds {
		shouldAppend := true
		for _, cnnId := range connections {
			if ppId.Hex() == cnnId.Hex() {
				shouldAppend = false
				break
			}
		}
		if shouldAppend {
			filtered = append(filtered, ppId)
		}
	}

	if len(filtered) > 0 {
		posts = append(posts, creator.store.GetAllPostsFromIds(filtered)...)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreationTime.After(posts[j].CreationTime)
	})

	return posts
}
