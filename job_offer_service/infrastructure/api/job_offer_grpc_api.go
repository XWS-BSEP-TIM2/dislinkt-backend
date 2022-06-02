package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	jobOffer, err := handler.service.Get(ctx, objectId)
	if err != nil {
		return nil, err
	}
	return &pb.GetJobOfferResponse{
		JobOffer: mapJobOffer(jobOffer),
	}, nil
}

func (handler *JobOfferHandler) DeleteJobOffer(ctx context.Context, request *pb.GetJobOfferRequest) (*pb.EmptyResponse, error) {

	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	_, err = handler.service.Delete(ctx, objectId)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, err
}

func (handler *JobOfferHandler) GetAllJobOffers(ctx context.Context, request *pb.EmptyJobOfferRequest) (*pb.GetAllJobOffersResponse, error) {
	jobOffers, err := handler.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllJobOffersResponse{
		JobOffers: []*pb.JobOffer{},
	}
	for _, profile := range jobOffers {
		current := mapJobOffer(profile)
		response.JobOffers = append(response.JobOffers, current)
	}
	return response, nil
}

func (handler *JobOfferHandler) CreateJobOffer(ctx context.Context, request *pb.CreateJobOfferRequest) (*pb.CreateJobOfferResponse, error) {

	profile := MapJobOffer(request.JobOffer)
	handler.service.Insert(ctx, &profile)
	return &pb.CreateJobOfferResponse{
		JobOffer: mapJobOffer(&profile),
	}, nil
}

func (handler *JobOfferHandler) UpdateJobOffer(ctx context.Context, request *pb.CreateJobOfferRequest) (*pb.CreateJobOfferResponse, error) {
	profile := MapJobOffer(request.GetJobOffer())
	handler.service.Update(ctx, &profile)
	return &pb.CreateJobOfferResponse{JobOffer: mapJobOffer(&profile)}, nil
}

func (handler *JobOfferHandler) SearchJobOffer(ctx context.Context, request *pb.SearchJobOfferRequest) (*pb.GetAllJobOffersResponse, error) {

	jobOffers, err := handler.service.Search(ctx, request.GetParam())
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

	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	jobOffers, err := handler.service.GetUserJobOffers(ctx, objectId)
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
