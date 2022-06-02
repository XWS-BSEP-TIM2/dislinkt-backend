package api

import (
	converter "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/converter"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
)

func mapJobOffer(jobOffer *domain.JobOffer) *pb.JobOffer {
	jobOfferPb := &pb.JobOffer{
		Id:                 jobOffer.Id.Hex(),
		Position:           jobOffer.Position,
		Seniority:          jobOffer.Seniority,
		Description:        jobOffer.Description,
		UserId:             jobOffer.UserId.Hex(),
		Technologies:       jobOffer.Technologies,
		CompanyName:        jobOffer.CompanyName,
		JobOfferUniqueCode: jobOffer.JobOfferUniqueCode,
	}
	return jobOfferPb
}

func MapJobOffer(jobOffer *pb.JobOffer) domain.JobOffer {
	domainProfile := domain.JobOffer{
		Id:                 converter.GetObjectId(jobOffer.GetId()),
		Position:           jobOffer.Position,
		Seniority:          jobOffer.Seniority,
		Description:        jobOffer.Description,
		Technologies:       jobOffer.Technologies,
		UserId:             converter.GetObjectId(jobOffer.UserId),
		CompanyName:        jobOffer.CompanyName,
		JobOfferUniqueCode: jobOffer.GetJobOfferUniqueCode(),
	}
	return domainProfile
}
