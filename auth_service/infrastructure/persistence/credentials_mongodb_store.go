package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "user_db"
	COLLECTION = "users"
)

type CreadentialsMongoDBStore struct {
	users *mongo.Collection
}

func NewCredentialsMongoDBStore(client *mongo.Client) *CreadentialsMongoDBStore {
	users := client.Database(DATABASE).Collection(COLLECTION)
	return &CreadentialsMongoDBStore{
		users: users,
	}
}

func (store *CreadentialsMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *CreadentialsMongoDBStore) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	filter := bson.M{"username": username}
	return store.filterOne(filter)
}

func (store *CreadentialsMongoDBStore) GetAll(ctx context.Context) ([]*domain.User, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *CreadentialsMongoDBStore) Insert(ctx context.Context, product *domain.User) (error, string) {
	result, err := store.users.InsertOne(context.TODO(), product)
	if err != nil {
		return err, ""
	}
	product.Id = result.InsertedID.(primitive.ObjectID)
	return nil, product.Id.Hex()
}

func (store *CreadentialsMongoDBStore) DeleteAll(ctx context.Context) {
	store.users.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *CreadentialsMongoDBStore) filter(filter interface{}) ([]*domain.User, error) {
	cursor, err := store.users.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *CreadentialsMongoDBStore) filterOne(filter interface{}) (product *domain.User, err error) {
	result := store.users.FindOne(context.TODO(), filter)
	err = result.Decode(&product)
	return
}

func decode(cursor *mongo.Cursor) (users []*domain.User, err error) {
	for cursor.Next(context.TODO()) {
		var product domain.User
		err = cursor.Decode(&product)
		if err != nil {
			return
		}
		users = append(users, &product)
	}
	err = cursor.Err()
	return
}
