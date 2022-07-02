package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "Create")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	service.store.DeleteAllUserTokens(ctx2, userId)
	tokenCode, _ := utils.GenerateRandomString(30)
	token := domain.ApiToken{UserId: userId, ApiCode: tokenCode, CreationDate: time.Now()}
	err, _ := service.store.Insert(ctx2, &token)
	if err != nil {
		return "", err
	}
	return tokenCode, err
}

func (service *ApiTokenService) Get(ctx context.Context, id primitive.ObjectID) (*domain.ApiToken, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.Get(ctx2, id)
}

func (service *ApiTokenService) Delete(ctx context.Context, id primitive.ObjectID) {
	span := tracer.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	service.store.DeleteById(ctx2, id)
}

func (service *ApiTokenService) GetByTokenCode(ctx context.Context, tokenCode string) (*domain.ApiToken, error) {
	span := tracer.StartSpanFromContext(ctx, "GetByTokenCode")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetByTokenCode(ctx2, tokenCode)
}
