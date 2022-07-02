package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *CreadentialsMongoDBStore) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetByUsername")
	defer span.Finish()

	filter := bson.M{"username": username}
	return store.filterOne(filter)
}

func (store *CreadentialsMongoDBStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetByEmail")
	defer span.Finish()

	filter := bson.M{"email": email}
	return store.filterOne(filter)
}

func (store *CreadentialsMongoDBStore) GetAll(ctx context.Context) ([]*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()

	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *CreadentialsMongoDBStore) Insert(ctx context.Context, product *domain.User) (error, string) {
	span := tracer.StartSpanFromContext(ctx, "Insert")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	result, err := store.users.InsertOne(ctx2, product)
	if err != nil {
		return err, ""
	}
	product.Id = result.InsertedID.(primitive.ObjectID)
	return nil, product.Id.Hex()
}

func (store *CreadentialsMongoDBStore) DeleteById(ctx context.Context, id primitive.ObjectID) (int64, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteById")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	result, err := store.users.DeleteOne(ctx2, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (store *CreadentialsMongoDBStore) Update(ctx context.Context, user *domain.User) error {
	span := tracer.StartSpanFromContext(ctx, "Update")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userToUpdate := bson.M{"_id": user.Id}
	updatedUser := bson.M{"$set": bson.M{
		"username":                 user.Username,
		"password":                 user.Password,
		"role":                     user.Role,
		"locked":                   user.Locked,
		"lockReason":               user.LockReason,
		"email":                    user.Email,
		"verified":                 user.Verified,
		"verificationCode":         user.VerificationCode,
		"verificationCodeTime":     user.VerificationCodeTime,
		"numOfErrTryLogin":         user.NumOfErrTryLogin,
		"lastErrTryLoginTime":      user.LastErrTryLoginTime,
		"recoveryPasswordCode":     user.RecoveryPasswordCode,
		"recoveryPasswordCodeTime": user.RecoveryPasswordCodeTime,
		"isTFAEnabled":             user.IsTFAEnabled,
		"TFASecret":                user.TFASecret,
	}}

	_, err := store.users.UpdateOne(ctx2, userToUpdate, updatedUser)

	if err != nil {
		return err
	}
	return nil
}

func (store *CreadentialsMongoDBStore) DeleteAll(ctx context.Context) {
	span := tracer.StartSpanFromContext(ctx, "DeleteAll")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	store.users.DeleteMany(ctx2, bson.D{{}})
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
