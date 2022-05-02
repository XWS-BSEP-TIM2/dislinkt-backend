package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileHandler struct {
	pb.UnimplementedProfileServiceServer
	service *application.ProfileService
}

func NewProfileHandler(service *application.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		service: service,
	}
}

func (handler *ProfileHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	profile, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{
		Profile: mapProfile(profile),
	}, nil
}

func (handler *ProfileHandler) GetAll(ctx context.Context, request *pb.EmptyRequest) (*pb.GetAllResponse, error) {
	profiles, err := handler.service.GetAll()
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Profiles: []*pb.Profile{},
	}
	for _, profile := range profiles {
		current := mapProfile(profile)
		response.Profiles = append(response.Profiles, current)
	}
	return response, nil
}

func (handler *ProfileHandler) CreateProfile(ctx context.Context, request *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {

	profile := MapProfile(request.GetProfile())
	handler.service.Insert(&profile)
	return &pb.CreateProfileResponse{
		Profile: mapProfile(&profile),
	}, nil
}

func (handler *ProfileHandler) UpdateProfile(ctx context.Context, request *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	profile := MapProfile(request.GetProfile())
	handler.service.Update(&profile)
	return &pb.CreateProfileResponse{Profile: mapProfile(&profile)}, nil
}
