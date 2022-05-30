package application

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/ecoding"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostService struct {
	store                    domain.PostStore
	authServiceAddress       string
	connectionServiceAddress string
	profileServiceAddress    string
}

func NewPostService(
	store domain.PostStore,
	authServiceAddress,
	connectionServiceAddress,
	profileServiceAddress string) *PostService {

	return &PostService{
		store:                    store,
		authServiceAddress:       authServiceAddress,
		connectionServiceAddress: connectionServiceAddress,
		profileServiceAddress:    profileServiceAddress,
	}
}

func (service *PostService) GetPost(ctx context.Context, id primitive.ObjectID) *domain.PostDetailsDTO {
	post, postNotFoundErr := service.store.Get(id)
	if postNotFoundErr != nil {
		log(fmt.Sprintf("Post with id: %v not found", id))
		panic(errors.NewEntityNotFoundError("Post with given id does not exist."))
	}
	return service.getPostDetails(post)
}

func (service *PostService) CreatePost(ctx context.Context, post *domain.Post) *domain.PostDetailsDTO {
	err := service.store.Insert(post)
	if err != nil {
		log("Error during post creation")
		panic(fmt.Errorf("error during post creation"))
	}
	return service.getPostDetails(post)
}

func (service *PostService) GetPosts(ctx context.Context) []*domain.PostDetailsDTO {
	posts, err := service.store.GetAll()
	if err != nil {
		log("Error loading posts")
		panic(errors.NewEntityNotFoundError("Posts unavailable."))
	}

	return service.getMultiplePostsDetails(posts)
}

func (service *PostService) getMultiplePostsDetails(posts []*domain.Post) []*domain.PostDetailsDTO {
	//profiles, err := serviceClients.NewProfileClient(service.profileServiceAddress).GetAll(ctx, &profileService.EmptyRequest{})
	//if err != nil {
	//	log(fmt.Sprintf("Error loading profiles: %v", err))
	//	panic(fmt.Errorf("posts unavailable"))
	//}
	postsDetails, ok := funk.Map(posts, service.getPostDetails).([]*domain.PostDetailsDTO)
	if !ok {
		log("Error in conversion of posts to postDetails")
		panic(fmt.Errorf("posts unavailable"))
	}
	return postsDetails
}

func (service *PostService) getPostDetails(post *domain.Post) *domain.PostDetailsDTO {
	return &domain.PostDetailsDTO{
		//Owner:       mapProfileToOwner(service.getPostOwnerProfile(ctx, post.OwnerId)),
		Post:        post,
		ImageBase64: ecoding.NewBase64Coder().Encode(post.Image),
		Stats: &domain.Stats{
			CommentsNumber: len(post.Comments),
			LikesNumber:    len(post.Likes),
			DislikesNumber: len(post.Dislikes),
		},
	}
}

func (service *PostService) GetAllPosts() ([]*domain.Post, error) {
	return service.store.GetAll()
}

//func mapProfileToOwner(ownerProfile *profileService.Profile) *domain.Owner {
//	return &domain.Owner{
//		UserId:   ownerProfile.Id,
//		Username: ownerProfile.Username,
//		Name:     ownerProfile.Name,
//		Surname:  ownerProfile.Surname,
//	}
//}

//func (service *PostService) getPostOwnerProfile(ctx context.Context, ownerId primitive.ObjectID) *profileService.Profile {
//	profileClient := serviceClients.NewProfileClient(service.profileServiceAddress)
//	hexId := ownerId.Hex()
//	profileResponse, err := profileClient.Get(ctx, &profileService.GetRequest{Id: hexId})
//	if err != nil {
//		log(fmt.Sprintf("Error getting post owner with id: %s", hexId))
//		panic(fmt.Errorf("error getting post owner"))
//	}
//	return profileResponse.Profile
//}

func log(message string) {
	fmt.Printf("[%v] [Post Service]: %s\n", time.Now(), message)
}
