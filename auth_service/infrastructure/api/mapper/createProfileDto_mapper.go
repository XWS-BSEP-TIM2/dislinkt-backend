package mapper

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/infrastructure/api/dto"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
)

func ProtoToCreateProfileDto(request *pb.RegisterRequest) dto.CreateProfileDto {
	return dto.CreateProfileDto{
		Gender:    request.Data.GetGender(),
		BirthDate: request.Data.BirthDate,
		Email:     request.Data.GetEmail(),
		FirstName: request.Data.FirstName,
		LastName:  request.Data.LastName,
		IsPrivate: request.Data.IsPrivate,
	}
}
