package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileService struct {
	store persistence.ProfileStore
}

func NewProfileService(store persistence.ProfileStore) *ProfileService {
	return &ProfileService{
		store: store,
	}
}

func (service *ProfileService) Get(ctx context.Context, id primitive.ObjectID) (*domain.Profile, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.Get(ctx2, id)
}

func (service *ProfileService) GetAll(ctx context.Context) ([]*domain.Profile, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAll(ctx2)
}

func (service *ProfileService) Insert(ctx context.Context, profile *domain.Profile) {
	span := tracer.StartSpanFromContext(ctx, "Insert")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.store.Insert(ctx2, profile)
}

func (service *ProfileService) Update(ctx context.Context, profile *domain.Profile) {
	span := tracer.StartSpanFromContext(ctx, "Update")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	//TODO if username changes update it in auth service
	service.store.Update(ctx2, profile)
}

func (service *ProfileService) Search(ctx context.Context, search string) ([]*domain.Profile, error) {
	span := tracer.StartSpanFromContext(ctx, "Search")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	return service.store.Search(ctx2, search)
}

func (service *ProfileService) DeleteById(ctx context.Context, id primitive.ObjectID) {
	span := tracer.StartSpanFromContext(ctx, "DeleteById")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.store.DeleteById(ctx2, id)
}
