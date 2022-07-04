package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/application"
)

type JobOfferHandler struct {
	pb.UnimplementedJobOfferServiceServer
	service *application.JobOfferService
}

func NewJobOfferHandler(service *application.JobOfferService) *JobOfferHandler {
	return &JobOfferHandler{
		service: service,
	}
}

func (handler *JobOfferHandler) GetJobOffer(ctx context.Context, request *pb.GetJobOfferRequest) (*pb.GetJobOfferResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetJobOffer")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	id := request.Id
	jobOffer, err := handler.service.Get(ctx2, id)
	if err != nil {
		return nil, err
	}
	return &pb.GetJobOfferResponse{
		JobOffer: mapJobOffer(jobOffer),
	}, nil
}

func (handler *JobOfferHandler) DeleteJobOffer(ctx context.Context, request *pb.GetJobOfferRequest) (*pb.EmptyResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteJobOffer")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	_, err := handler.service.Delete(ctx2, request.Id)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, err
}

func (handler *JobOfferHandler) GetAllJobOffers(ctx context.Context, request *pb.EmptyJobOfferRequest) (*pb.GetAllJobOffersResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllJobOffers")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	jobOffers, err := handler.service.GetAll(ctx2)
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllJobOffersResponse{
		JobOffers: []*pb.JobOffer{},
	}
	for _, jobs := range jobOffers {
		current := mapJobOffer(jobs)
		response.JobOffers = append(response.JobOffers, current)
	}
	return response, nil
}

func (handler *JobOfferHandler) CreateJobOffer(ctx context.Context, request *pb.CreateJobOfferRequest) (*pb.CreateJobOfferResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateJobOffer")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	jobOffer := MapJobOffer(request.JobOffer)
	handler.service.Insert(ctx2, &jobOffer)
	return &pb.CreateJobOfferResponse{
		JobOffer: mapJobOffer(&jobOffer),
	}, nil
}

func (handler *JobOfferHandler) UpdateJobOffer(ctx context.Context, request *pb.CreateJobOfferRequest) (*pb.CreateJobOfferResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateJobOffer")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	jobOffer := MapJobOffer(request.GetJobOffer())
	handler.service.Update(ctx2, &jobOffer)
	return &pb.CreateJobOfferResponse{JobOffer: mapJobOffer(&jobOffer)}, nil
}

func (handler *JobOfferHandler) SearchJobOffer(ctx context.Context, request *pb.SearchJobOfferRequest) (*pb.GetAllJobOffersResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "SearchJobOffer")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	jobOffers, err := handler.service.Search(ctx2, request.GetParam())
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllJobOffersResponse{
		JobOffers: []*pb.JobOffer{},
	}
	for _, jobOffer := range jobOffers {
		current := mapJobOffer(jobOffer)
		response.JobOffers = append(response.JobOffers, current)
	}
	return response, nil
}

func (handler *JobOfferHandler) GetUserJobOffers(ctx context.Context, request *pb.GetJobOfferRequest) (*pb.GetAllJobOffersResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetUserJobOffers")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userID := request.Id
	jobOffers, err := handler.service.GetUserJobOffers(ctx2, userID)
	if err != nil {
		return nil, err
	}

	response := &pb.GetAllJobOffersResponse{
		JobOffers: []*pb.JobOffer{},
	}
	for _, jobOffer := range jobOffers {
		current := mapJobOffer(jobOffer)
		response.JobOffers = append(response.JobOffers, current)
	}
	return response, nil
}

func (handler *JobOfferHandler) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateUser")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.CreateUser(ctx2, request.UserID)
}

func (handler *JobOfferHandler) UpdateUserSkills(ctx context.Context, request *pb.UpdateUserSkillsRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserSkills")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.UpdateUserSkills(ctx2, request.UserID, request.Technologies)
}

func (handler *JobOfferHandler) GetRecommendationJobOffer(ctx context.Context, request *pb.GetRecommendationJobOfferRequest) (*pb.GetAllJobOffersResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetRecommendationJobOffer")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	jobOffers, err := handler.service.GetRecommendationJobOffer(ctx2, request.UserID)
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllJobOffersResponse{
		JobOffers: []*pb.JobOffer{},
	}
	for _, jobs := range jobOffers {
		current := mapJobOffer(jobs)
		response.JobOffers = append(response.JobOffers, current)
	}
	return response, nil
}
