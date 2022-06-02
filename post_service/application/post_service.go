package application

import (
	"context"
	"fmt"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
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
	store                    domain.PostStore
	authServiceAdapter       asa.IAuthServiceAdapter
	connectionServiceAdapter csa.IConnectionServiceAdapter
	profileServiceAdapter    psa.IProfileServiceAdapter
	postAccessValidator      *util.PostAccessValidator
	ownerFinder              *util.OwnerFinder
	feedCreator              *util.FeedCreator
}

func NewPostService(
	store domain.PostStore,
	authServiceAddress,
	connectionServiceAddress,
	profileServiceAddress string) *PostService {

	authServiceAdapter := asa.NewAuthServiceAdapter(authServiceAddress)
	connectionServiceAdapter := csa.NewConnectionServiceAdapter(connectionServiceAddress)
	profileServiceAdapter := psa.NewProfileServiceAdapter(profileServiceAddress)
	postAccessValidator := util.NewPostAccessValidator(store, authServiceAdapter, connectionServiceAdapter)
	ownerFinder := util.NewOwnerFinder(profileServiceAdapter)
	feedCreator := util.NewFeedCreator(store, connectionServiceAdapter, profileServiceAdapter)
	return &PostService{
		store:                    store,
		authServiceAdapter:       authServiceAdapter,
		connectionServiceAdapter: connectionServiceAdapter,
		profileServiceAdapter:    profileServiceAdapter,
		postAccessValidator:      postAccessValidator,
		ownerFinder:              ownerFinder,
		feedCreator:              feedCreator,
	}
}

func (service *PostService) GetPost(ctx context.Context, id primitive.ObjectID) *domain.PostDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, id)
	post, postNotFoundErr := service.store.Get(id)
	if postNotFoundErr != nil {
		log(fmt.Sprintf("Post with id: %v not found", id))
		panic(errors.NewEntityNotFoundError("Post with given id does not exist."))
	}
	return service.getPostDetailsMapper(ctx)(post)
}

func (service *PostService) CreatePost(ctx context.Context, post *domain.Post) *domain.PostDetailsDTO {
	service.authServiceAdapter.ValidateCurrentUser(ctx, post.OwnerId)
	err := service.store.Insert(post)
	if err != nil {
		log("Error during post creation")
		panic(fmt.Errorf("error during post creation"))
	}
	return service.getPostDetailsMapper(ctx)(post)
}

func (service *PostService) GetPosts(ctx context.Context) []*domain.PostDetailsDTO {
	currentUserId := service.authServiceAdapter.GetRequesterId(ctx)
	posts := service.feedCreator.CreateFeedForUser(ctx, currentUserId)
	return service.getMultiplePostsDetails(ctx, posts)
}

func (service *PostService) GetPostsFromUser(ctx context.Context, userId primitive.ObjectID) []*domain.PostDetailsDTO {
	currentUserId := service.authServiceAdapter.GetRequesterId(ctx)
	if currentUserId != userId {
		result := service.connectionServiceAdapter.CanUserAccessPostFromOwner(ctx, currentUserId, userId)
		if !result {
			panic(errors.NewEntityForbiddenError("Current user cannot access posts from given user."))
		}
	}
	posts, err := service.store.GetPostsFromUser(userId)
	if err != nil {
		log("Error loading posts")
		panic(errors.NewEntityNotFoundError("Posts unavailable."))
	}

	return service.getMultiplePostsDetails(ctx, posts)
}

func (service *PostService) getMultiplePostsDetails(ctx context.Context, posts []*domain.Post) []*domain.PostDetailsDTO {
	postsDetails, ok := funk.Map(posts, service.getPostDetailsMapper(ctx)).([]*domain.PostDetailsDTO)
	if !ok {
		log("Error in conversion of posts to postDetails")
		panic(fmt.Errorf("posts unavailable"))
	}
	return postsDetails
}

func (service *PostService) getPostDetailsMapper(ctx context.Context) func(post *domain.Post) *domain.PostDetailsDTO {
	currentUserId := service.authServiceAdapter.GetRequesterId(ctx)
	getOwner := service.ownerFinder.GetOwnerFinderFunction(ctx)
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

func (service *PostService) GetAllPosts() ([]*domain.Post, error) {
	return service.store.GetAll()
}

//func mapProfileToOwner(ownerProfile *profileService.JobOffer) *domain.Owner {
//	return &domain.Owner{
//		UserId:   ownerProfile.Id,
//		Description: ownerProfile.Description,
//		Position:     ownerProfile.Position,
//		Seniority:  ownerProfile.Seniority,
//	}
//}

func (service *PostService) getPostOwnerProfile(ctx context.Context, ownerId primitive.ObjectID) *profileService.Profile {
	profileClient := serviceClients.NewProfileClient(service.profileServiceAddress)
	hexId := ownerId.Hex()
	profileResponse, err := profileClient.Get(ctx, &profileService.GetRequest{Id: hexId})
	if err != nil {
		log(fmt.Sprintf("Error getting post owner with id: %s", hexId))
		panic(fmt.Errorf("error getting post owner"))
	}
	return profileResponse.Profile
}

func log(message string) {
	fmt.Printf("[%v] [Post Service]: %s\n", time.Now(), message)
}
