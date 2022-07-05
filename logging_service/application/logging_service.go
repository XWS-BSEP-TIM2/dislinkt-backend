package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	ERROR   = "ERROR"
	INFO    = "INFO"
	WARNING = "WARNING"
	SUCCESS = "SUCCESS"
)

type LoggingService struct {
	store       persistence.LoggingStore
	eventsStore persistence.EventsStore
}

func NewLoggingService(store persistence.LoggingStore, eventsStore persistence.EventsStore) *LoggingService {
	return &LoggingService{
		store:       store,
		eventsStore: eventsStore,
	}
}

func (s LoggingService) LoggInfo(ctx context.Context, request *logging_service.LogRequest) (*logging_service.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "LoggInfo")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	newLog := domain.NewLog(request.ServiceName, request.ServiceFunctionName, INFO, request.UserID, request.IpAddress, request.Description)
	return s.store.SaveLog(ctx2, newLog)
}

func (s LoggingService) LoggError(ctx context.Context, request *logging_service.LogRequest) (*logging_service.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "LoggError")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	newLog := domain.NewLog(request.ServiceName, request.ServiceFunctionName, ERROR, request.UserID, request.IpAddress, request.Description)
	return s.store.SaveLog(ctx2, newLog)
}

func (s LoggingService) LoggWarning(ctx context.Context, request *logging_service.LogRequest) (*logging_service.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "LoggWarning")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	newLog := domain.NewLog(request.ServiceName, request.ServiceFunctionName, WARNING, request.UserID, request.IpAddress, request.Description)
	return s.store.SaveLog(ctx2, newLog)
}

func (s LoggingService) LoggSuccess(ctx context.Context, request *logging_service.LogRequest) (*logging_service.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "LoggSuccess")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	newLog := domain.NewLog(request.ServiceName, request.ServiceFunctionName, SUCCESS, request.UserID, request.IpAddress, request.Description)
	return s.store.SaveLog(ctx2, newLog)
}

func (s LoggingService) InsertEvent(ctx context.Context, in *logging_service.EventRequest) (*logging_service.Empty, error) {
	span := tracer.StartSpanFromContext(ctx, "LoggSuccess")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	ownerId, err := primitive.ObjectIDFromHex(in.UserId)

	if err == nil {
		var newEvent = domain.Event{
			UserId:      ownerId,
			Title:       in.Title,
			Description: in.Description,
			Date:        time.Now(),
		}
		s.eventsStore.Insert(ctx2, &newEvent)
	}

	return &logging_service.Empty{}, err
}

func (s LoggingService) GetEvents(ctx context.Context, in *logging_service.Empty) (*logging_service.GetEventsResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "LoggSuccess")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	var fetchedEvents []*pb.EventRequest
	events, err := s.eventsStore.GetAll(ctx2)

	if err == nil {
		for _, event := range events {
			var newEvent = pb.EventRequest{
				UserId:      event.UserId.Hex(),
				Title:       event.Title,
				Description: event.Description,
				Date:        &timestamppb.Timestamp{Seconds: event.Date.Unix()},
			}
			fetchedEvents = append(fetchedEvents, &newEvent)
		}
	}

	return &pb.GetEventsResponse{
		Events: fetchedEvents,
	}, nil
}
