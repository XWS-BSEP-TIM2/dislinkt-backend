package profile_service_adapter

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/services"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileServiceAdapter struct {
	address string
}

func NewProfileServiceAdapter(address string) *ProfileServiceAdapter {
	return &ProfileServiceAdapter{address: address}
}

func (p *ProfileServiceAdapter) GetAllProfiles(ctx context.Context) []*domain.Owner {
	profileClient := services.NewProfileClient(p.address)
	response, profileErr := profileClient.GetAll(ctx, &pb.EmptyRequest{})

	if profileErr != nil {
		panic(fmt.Errorf("Error during getting all profiles: Profile Service"))
	}
	res, ok := funk.Map(response.Profiles, mapProfileToOwner).([]*domain.Owner)
	if !ok {
		panic(fmt.Errorf("Cannot cast list as []*domain.Owner"))
	}
	return res
}

func (p *ProfileServiceAdapter) GetSingleProfile(ctx context.Context, profileId primitive.ObjectID) *domain.Owner {
	profileClient := services.NewProfileClient(p.address)
	response, profileErr := profileClient.Get(ctx, &pb.GetRequest{Id: profileId.Hex()})

	if profileErr != nil {
		panic(fmt.Errorf("Error during getting all profiles: Profile Service"))
	}
	return mapProfileToOwner(response.Profile)
}

func (p *ProfileServiceAdapter) GetAllPublicProfilesIds(ctx context.Context) []*primitive.ObjectID {
	profileClient := services.NewProfileClient(p.address)
	response, profileErr := profileClient.GetAll(ctx, &pb.EmptyRequest{})

	if profileErr != nil {
		panic(fmt.Errorf("Error during getting all profiles: Profile Service"))
	}

	publicProfiles := funk.Filter(response.Profiles, func(profile *pb.Profile) bool {
		return !profile.IsPrivate
	}).([]*pb.Profile)
	res := funk.Map(publicProfiles, getProfileId).([]*primitive.ObjectID)

	return res
}

func mapProfileToOwner(profile *pb.Profile) *domain.Owner {
	profileId := getProfileId(profile)
	return &domain.Owner{
		UserId:   *profileId,
		Username: profile.Username,
		Name:     profile.Name,
		Surname:  profile.Surname,
	}
}

func getProfileId(profile *pb.Profile) *primitive.ObjectID {
	profileId, err := primitive.ObjectIDFromHex(profile.Id)
	if err != nil {
		panic(fmt.Errorf("Given profile id is invalid."))
	}
	return &profileId
}
