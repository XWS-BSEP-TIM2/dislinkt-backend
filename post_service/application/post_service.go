package application

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
	lsa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/logging_service_adapter"
	nsa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/notification_service_adapter"
	psa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/profile_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/util"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/ecoding"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostService struct {
	store                      domain.PostStore
	authServiceAdapter         asa.IAuthServiceAdapter
	connectionServiceAdapter   csa.IConnectionServiceAdapter
	notificationServiceAdapter nsa.INotificationServiceAdapter
	profileServiceAdapter      psa.IProfileServiceAdapter
	loggingServiceAdapter      lsa.ILoggingServiceAdapter
	postAccessValidator        *util.PostAccessValidator
	ownerFinder                *util.OwnerFinder
	feedCreator                *util.FeedCreator
}

func NewPostService(
	store domain.PostStore,
	authServiceAddress,
	connectionServiceAddress,
	profileServiceAddress,
	notificationServiceAddress,
	loggingServiceAddress string) *PostService {

	loggingServiceAdapter := lsa.NewLoggingServiceAdapter(loggingServiceAddress)
	authServiceAdapter := asa.NewAuthServiceAdapter(authServiceAddress, loggingServiceAdapter)
	connectionServiceAdapter := csa.NewConnectionServiceAdapter(connectionServiceAddress, loggingServiceAdapter)
	notificationServiceAdapter := nsa.NewNotificationServiceAdapter(notificationServiceAddress)
	profileServiceAdapter := psa.NewProfileServiceAdapter(profileServiceAddress, loggingServiceAdapter)
	postAccessValidator := util.NewPostAccessValidator(store, authServiceAdapter, connectionServiceAdapter, loggingServiceAdapter)
	ownerFinder := util.NewOwnerFinder(profileServiceAdapter)
	feedCreator := util.NewFeedCreator(store, connectionServiceAdapter, profileServiceAdapter)
	return &PostService{
		store:                      store,
		authServiceAdapter:         authServiceAdapter,
		connectionServiceAdapter:   connectionServiceAdapter,
		notificationServiceAdapter: notificationServiceAdapter,
		profileServiceAdapter:      profileServiceAdapter,
		loggingServiceAdapter:      loggingServiceAdapter,
		postAccessValidator:        postAccessValidator,
		ownerFinder:                ownerFinder,
		feedCreator:                feedCreator,
	}
}

func (service *PostService) GetPost(ctx context.Context, id primitive.ObjectID) *domain.PostDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GetPost")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, id)
	post, postNotFoundErr := service.store.Get(id)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)

	if postNotFoundErr != nil {
		message := fmt.Sprintf("Post with id: %s not found", id.Hex())
		service.loggingServiceAdapter.Log(ctx2, "WARNING", "GetPost", requesterId.Hex(), message)
		panic(errors.NewEntityNotFoundError("Post with given id does not exist."))
	}
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GetPost", requesterId.Hex(), "User fetched post.")
	return service.getPostDetailsMapper(ctx2)(post)
}

func (service *PostService) CreatePost(ctx context.Context, post *domain.Post) *domain.PostDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "CreatePost")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.authServiceAdapter.ValidateCurrentUser(ctx2, post.OwnerId)
	err := service.store.Insert(post)
	if err != nil {
		message := "Error during post creation."
		service.loggingServiceAdapter.Log(ctx2, "ERROR", "CreatePost", post.OwnerId.Hex(), message)
		panic(fmt.Errorf(message))
	}

	postOwner := service.profileServiceAdapter.GetSingleProfile(ctx2, post.OwnerId)

	connectionIds := service.connectionServiceAdapter.GetAllUserConnections(ctx2, postOwner.UserId)

	for _, id := range connectionIds {
		var notification pb.Notification
		notification.OwnerId = id.Hex()
		notification.ForwardUrl = "posts/" + post.Id.Hex()
		notification.Text = "posted on their profile"
		notification.UserFullName = postOwner.Name + " " + postOwner.Surname
		service.notificationServiceAdapter.InsertNotification(ctx2, &pb.InsertNotificationRequest{Notification: &notification})
	}

	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "CreatePost", post.OwnerId.Hex(), "User created a new post.")
	return service.getPostDetailsMapper(ctx2)(post)
}

func (service *PostService) GetPosts(ctx context.Context) []*domain.PostDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GetPosts")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	currentUserId := service.authServiceAdapter.GetRequesterId(ctx2)
	posts := service.feedCreator.CreateFeedForUser(ctx2, currentUserId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GetPosts", requesterId.Hex(), "User fetched posts for feed.")
	return service.getMultiplePostsDetails(ctx2, posts)
}

func (service *PostService) GetPostsFromUser(ctx context.Context, userId primitive.ObjectID) []*domain.PostDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GetPostsFromUser")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	currentUserId := service.authServiceAdapter.GetRequesterId(ctx2)
	if currentUserId != userId {
		result := service.connectionServiceAdapter.CanUserAccessPostFromOwner(ctx2, currentUserId, userId)
		if !result {
			message := fmt.Sprintf("Current user (id: %s) cannot access posts from user with id: %s.", currentUserId.Hex(), userId.Hex())
			service.loggingServiceAdapter.Log(ctx2, "WARNING", "GetPostsFromUser", currentUserId.Hex(), message)
			panic(errors.NewEntityForbiddenError(message))
		}
	}
	posts, err := service.store.GetPostsFromUser(userId)
	if err != nil {
		message := fmt.Sprintf("Posts from user id: %s unavailabe", userId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "WARNING", "GetPostsFromUser", currentUserId.Hex(), message)
		panic(errors.NewEntityNotFoundError(message))
	}

	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GetPostsFromUser", currentUserId.Hex(), "User fetched posts from other profile.")
	return service.getMultiplePostsDetails(ctx2, posts)
}

func (service *PostService) getMultiplePostsDetails(ctx context.Context, posts []*domain.Post) []*domain.PostDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "getMultiplePostsDetails")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	postsDetails, ok := funk.Map(posts, service.getPostDetailsMapper(ctx2)).([]*domain.PostDetailsDTO)
	if !ok {
		log("Error in conversion of posts to postDetails")
		panic(fmt.Errorf("posts unavailable"))
	}
	return postsDetails
}

func (service *PostService) getPostDetailsMapper(ctx context.Context) func(post *domain.Post) *domain.PostDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "getPostDetailsMapper")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	currentUserId := service.authServiceAdapter.GetRequesterId(ctx2)
	getOwner := service.ownerFinder.GetOwnerFinderFunction(ctx2)
	return func(post *domain.Post) *domain.PostDetailsDTO {
		var reactions *domain.Reactions
		if currentUserId == primitive.NilObjectID {
			reactions = &domain.Reactions{
				Liked:    false,
				Disliked: false,
			}
		} else {
			reactions = service.store.GetReactions(post.Id, currentUserId)
		}
		return &domain.PostDetailsDTO{
			Owner:       getOwner(post.OwnerId),
			Post:        post,
			ImageBase64: ecoding.NewBase64Coder().Encode(post.Image),
			Stats: &domain.Stats{
				CommentsNumber: len(post.Comments),
				LikesNumber:    len(post.Likes),
				DislikesNumber: len(post.Dislikes),
			},
			Reactions: reactions,
		}
	}
}

func log(message string) {
	fmt.Printf("[%v] [Post Service]: %s\n", time.Now(), message)
}
