package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	PASSWORDLESS_TOKEN_COLLECTION = "passwordless_tokens"
)

type PasswordlessTokenMongoDBStore struct {
	passwordlessTokens *mongo.Collection
}

func NewPasswordlessTokenMongoDBStore(client *mongo.Client) PasswordlessTokenMongoDBStore {
	passwordless_token := client.Database(DATABASE).Collection(PASSWORDLESS_TOKEN_COLLECTION)
	return PasswordlessTokenMongoDBStore{
		passwordlessTokens: passwordless_token,
	}
}

func (store *PasswordlessTokenMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.PasswordlessToken, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *PasswordlessTokenMongoDBStore) GetByTokenCode(ctx context.Context, tokenCode string) (*domain.PasswordlessToken, error) {
	filter := bson.M{"token_code": tokenCode}
	return store.filterOne(filter)
}

func (store *PasswordlessTokenMongoDBStore) Insert(ctx context.Context, product *domain.PasswordlessToken) (error, string) {
	result, err := store.passwordlessTokens.InsertOne(ctx, product)
	if err != nil {
		return err, ""
	}
	product.Id = result.InsertedID.(primitive.ObjectID)
	return nil, product.Id.Hex()
}

func (store PasswordlessTokenMongoDBStore) DeleteById(ctx context.Context, id primitive.ObjectID) (int64, error) {
	result, err := store.passwordlessTokens.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (store *PasswordlessTokenMongoDBStore) filter(filter interface{}) ([]*domain.PasswordlessToken, error) {
	cursor, err := store.passwordlessTokens.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return store.decode(cursor)
}

func (store *PasswordlessTokenMongoDBStore) filterOne(filter interface{}) (product *domain.PasswordlessToken, err error) {
	result := store.passwordlessTokens.FindOne(context.TODO(), filter)
	err = result.Decode(&product)
	return
}

func (store *PasswordlessTokenMongoDBStore) decode(cursor *mongo.Cursor) (users []*domain.PasswordlessToken, err error) {
	for cursor.Next(context.TODO()) {
		var product domain.PasswordlessToken
		err = cursor.Decode(&product)
		if err != nil {
			return
		}
		users = append(users, &product)
	}
	err = cursor.Err()
	return
}
