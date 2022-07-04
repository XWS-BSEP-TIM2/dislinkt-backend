package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "events_db"
	COLLECTION = "events"
)

type EventsMongoDbStore struct {
	events *mongo.Collection
}

func (store *EventsMongoDbStore) DeleteAll(todo context.Context) {
	span := tracer.StartSpanFromContext(todo, "DeleteAll")
	defer span.Finish()

	store.events.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *EventsMongoDbStore) filter(filter interface{}) ([]*domain.Event, error) {
	cursor, err := store.events.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *EventsMongoDbStore) filterOne(filter interface{}) (notification *domain.Event, err error) {
	result := store.events.FindOne(context.TODO(), filter)
	err = result.Decode(&notification)
	return
}
func decode(cursor *mongo.Cursor) (events []*domain.Event, err error) {
	for cursor.Next(context.TODO()) {
		var event domain.Event
		err = cursor.Decode(&event)
		if err != nil {
			return
		}
		events = append(events, &event)
	}
	err = cursor.Err()
	return
}

func (store EventsMongoDbStore) GetAll(ctx context.Context) ([]*domain.Event, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()

	filter := bson.D{{}}
	return store.filter(filter)
}

func (store EventsMongoDbStore) Insert(ctx context.Context, event *domain.Event) error {
	span := tracer.StartSpanFromContext(ctx, "Insert")
	defer span.Finish()

	_, err := store.events.InsertOne(context.TODO(), event)
	if err != nil {
		return err
	}
	return nil
}

func NewEventsMongoDbStore(client *mongo.Client) EventsStore {
	eventsDb := client.Database(DATABASE).Collection(COLLECTION)
	return &EventsMongoDbStore{
		events: eventsDb,
	}
}
