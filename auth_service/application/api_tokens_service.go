package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ApiTokenService struct {
	store domain.ApiTokenStore
}

func NewApiTokenService(store domain.ApiTokenStore) *ApiTokenService {
	return &ApiTokenService{
		store: store,
	}
}

func (service *ApiTokenService) Create(ctx context.Context, userId primitive.ObjectID) (string, error) {
	service.store.DeleteAllUserTokens(ctx, userId)
	tokenCode, _ := utils.GenerateRandomString(30)
	token := domain.ApiToken{UserId: userId, ApiCode: tokenCode, CreationDate: time.Now()}
	err, _ := service.store.Insert(ctx, &token)
	if err != nil {
		return "", err
	}
	return tokenCode, err
}

func (service *ApiTokenService) Get(ctx context.Context, id primitive.ObjectID) (*domain.ApiToken, error) {
	return service.store.Get(ctx, id)
}

func (service *ApiTokenService) Delete(ctx context.Context, id primitive.ObjectID) {
	service.store.DeleteById(ctx, id)
}

func (service *ApiTokenService) GetByTokenCode(ctx context.Context, tokenCode string) (*domain.ApiToken, error) {
	return service.store.GetByTokenCode(ctx, tokenCode)
}
