package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/update_skills"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileHandler struct {
	pb.UnimplementedProfileServiceServer
	service                  *application.ProfileService
	updateSkillsOrchestrator *application.UpdateSkillsOrchestrator
}

func NewProfileHandler(service *application.ProfileService, updateSkillsOrchestrator *application.UpdateSkillsOrchestrator) *ProfileHandler {
	return &ProfileHandler{
		service:                  service,
		updateSkillsOrchestrator: updateSkillsOrchestrator,
	}
}

func (handler *ProfileHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	profile, err := handler.service.Get(ctx2, objectId)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{
		Profile: mapProfile(profile),
	}, nil
}

func (handler *ProfileHandler) GetAll(ctx context.Context, request *pb.EmptyRequest) (*pb.GetAllResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	profiles, err := handler.service.GetAll(ctx2)
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
	span := tracer.StartSpanFromContext(ctx, "CreateProfile")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	profile := MapProfile(request.Profile)
	handler.service.Insert(ctx2, &profile)
	return &pb.CreateProfileResponse{
		Profile: mapProfile(&profile),
	}, nil
}

func (handler *ProfileHandler) UpdateProfile(ctx context.Context, request *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateProfile")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	profile := MapProfile(request.GetProfile())
	handler.service.Update(ctx2, &profile)

	return &pb.CreateProfileResponse{Profile: mapProfile(&profile)}, nil
}

func (handler *ProfileHandler) UpdateProfileSkills(ctx context.Context, request *pb.UpdateProfileSkillsRequest) (*pb.UpdateProfileSkillsResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateProfileSkills")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	profile := MapProfileOnlySkills(request.Id, request.Skills)

	oldProfile, errGetOldProfile := handler.service.Get(ctx2, profile.Id)
	if errGetOldProfile != nil {
		return &pb.UpdateProfileSkillsResponse{Status: 400, Msg: errGetOldProfile.Error()}, errGetOldProfile
	}

	err := handler.service.UpdateSkills(ctx2, &profile)

	if err != nil {
		return &pb.UpdateProfileSkillsResponse{Status: 400, Msg: err.Error()}, err
	}

	//poziv saga
	newSkillsForJobOffer := handler.getNewSkills(&profile)
	handler.updateSkillsOrchestrator.Start(events.UpdateSkillsDetails{UserID: request.Id, OldSkills: oldProfile.Skills, NewSkillsForJobOffer: newSkillsForJobOffer})

	return &pb.UpdateProfileSkillsResponse{Status: 200, Msg: "ok"}, nil
}

func (handler *ProfileHandler) getNewSkills(profile *domain.Profile) []string {

	var technologies []string
	addTechnologie := true
	for _, s1 := range profile.Skills {
		addTechnologie = true
		for _, t := range technologies {
			if s1.Name == t {
				addTechnologie = false
				break
			}
		}
		if addTechnologie {
			technologies = append(technologies, s1.Name)
		}
	}

	return technologies
}

func (handler *ProfileHandler) SearchProfile(ctx context.Context, request *pb.SearchProfileRequest) (*pb.GetAllResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "SearchProfile")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	profiles, err := handler.service.Search(ctx2, request.GetParam())
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
