package profile_service_adapter

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters"
	lsa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/logging_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileServiceAdapter struct {
	address               string
	loggingServiceAdapter lsa.ILoggingServiceAdapter
}

func NewProfileServiceAdapter(address string, loggingServiceAdapter lsa.ILoggingServiceAdapter) *ProfileServiceAdapter {
	return &ProfileServiceAdapter{address: address, loggingServiceAdapter: loggingServiceAdapter}
}

func (p *ProfileServiceAdapter) GetAllProfiles(ctx context.Context) []*domain.Owner {
	span := tracer.StartSpanFromContext(ctx, "GetAllProfiles")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	profileClient := adapters.NewProfileClient(p.address)
	response, profileErr := profileClient.GetAll(ctx2, &pb.EmptyRequest{})

	if profileErr != nil {
		message := "Error during getting all profiles: Profile Service"
		p.loggingServiceAdapter.Log(ctx2, "ERROR", "GetAllProfiles", "N/A", message)
		panic(fmt.Errorf(message))
	}
	res, ok := funk.Map(response.Profiles, mapProfileToOwner).([]*domain.Owner)
	if !ok {
		panic(fmt.Errorf("Cannot cast list as []*domain.Owner"))
	}
	return res
}

func (p *ProfileServiceAdapter) GetSingleProfile(ctx context.Context, profileId primitive.ObjectID) *domain.Owner {
	span := tracer.StartSpanFromContext(ctx, "GetSingleProfile")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	profileClient := adapters.NewProfileClient(p.address)
	response, profileErr := profileClient.Get(ctx2, &pb.GetRequest{Id: profileId.Hex()})

	if profileErr != nil {
		message := "Error during getting all profiles: Profile Service"
		p.loggingServiceAdapter.Log(ctx2, "ERROR", "GetSingleProfile", "N/A", message)
		panic(fmt.Errorf(message))
	}
	return mapProfileToOwner(response.Profile)
}

func (p *ProfileServiceAdapter) GetAllPublicProfilesIds(ctx context.Context) []*primitive.ObjectID {
	span := tracer.StartSpanFromContext(ctx, "GetAllPublicProfilesIds")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	profileClient := adapters.NewProfileClient(p.address)
	response, profileErr := profileClient.GetAll(ctx2, &pb.EmptyRequest{})

	if profileErr != nil {
		message := "Error during getting all profiles: Profile Service"
		p.loggingServiceAdapter.Log(ctx2, "ERROR", "GetAllPublicProfilesIds", "N/A", message)
		panic(fmt.Errorf(message))
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
