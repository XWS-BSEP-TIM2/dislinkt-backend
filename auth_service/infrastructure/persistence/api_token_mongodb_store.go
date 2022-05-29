package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	API_TOKEN_COLLECTION = "api_tokens"
)

type ApiTokenMongoDBStore struct {
	apiTokens *mongo.Collection
}

func NewApiTokenMongoDBStore(client *mongo.Client) ApiTokenMongoDBStore {
	passwordless_token := client.Database(DATABASE).Collection(API_TOKEN_COLLECTION)
	return ApiTokenMongoDBStore{
		apiTokens: passwordless_token,
	}
}

func (store *ApiTokenMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.ApiToken, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *ApiTokenMongoDBStore) GetByTokenCode(ctx context.Context, tokenCode string) (*domain.ApiToken, error) {
	filter := bson.M{"token_code": tokenCode}
	return store.filterOne(filter)
}

func (store *ApiTokenMongoDBStore) Insert(ctx context.Context, token *domain.ApiToken) (error, string) {
	result, err := store.apiTokens.InsertOne(ctx, token)
	if err != nil {
		return err, ""
	}
	token.Id = result.InsertedID.(primitive.ObjectID)
	return nil, token.Id.Hex()
}

func (store ApiTokenMongoDBStore) DeleteById(ctx context.Context, id primitive.ObjectID) (int64, error) {
	result, err := store.apiTokens.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (store ApiTokenMongoDBStore) DeleteAllUserTokens(ctx context.Context, id primitive.ObjectID) error {
	_, err := store.apiTokens.DeleteMany(ctx, bson.M{"user_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (store *ApiTokenMongoDBStore) filter(filter interface{}) ([]*domain.ApiToken, error) {
	cursor, err := store.apiTokens.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return store.decode(cursor)
}

func (store *ApiTokenMongoDBStore) filterOne(filter interface{}) (product *domain.ApiToken, err error) {
	result := store.apiTokens.FindOne(context.TODO(), filter)
	err = result.Decode(&product)
	return
}

func (store *ApiTokenMongoDBStore) decode(cursor *mongo.Cursor) (users []*domain.ApiToken, err error) {
	for cursor.Next(context.TODO()) {
		var product domain.ApiToken
		err = cursor.Decode(&product)
		if err != nil {
			return
		}
		users = append(users, &product)
	}
	err = cursor.Err()
	return
}
