package persistence

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "post_db"
	COLLECTION = "post"
)

type PostMongoDBStore struct {
	posts *mongo.Collection
}

func NewPostMongoDBStore(client *mongo.Client) domain.PostStore {
	posts := client.Database(DATABASE).Collection(COLLECTION)
	return &PostMongoDBStore{
		posts: posts,
	}
}

func (store *PostMongoDBStore) Get(id primitive.ObjectID) (*domain.Post, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *PostMongoDBStore) GetAll() ([]*domain.Post, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *PostMongoDBStore) Insert(post *domain.Post) error {
	result, err := store.posts.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	post.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *PostMongoDBStore) Update(post *domain.Post) error {
	_, err := store.posts.UpdateOne(context.TODO(), bson.M{"_id": post.Id}, bson.M{"$set": post})
	if err != nil {
		return err
	}
	return nil
}

func (store *PostMongoDBStore) InsertComment(postId primitive.ObjectID, comment *domain.Comment) error {
	post, err := store.Get(postId)
	if err != nil {
		return err
	}
	comment.Id = primitive.NewObjectID()
	post.Comments = append(post.Comments, comment)

	err = store.Update(post)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostMongoDBStore) GetComment(postId primitive.ObjectID, commentId primitive.ObjectID) (*domain.Comment, error) {
	// should be fetched by mongo syntax
	post, err := store.Get(postId)
	if err != nil {
		panic(fmt.Errorf("Invalid post"))
	}
	for _, comment := range post.Comments {
		if comment.Id == commentId {
			return comment, nil
		}
	}
	panic(errors.NewEntityNotFoundError("Comment with given id not found."))
}

func (store *PostMongoDBStore) GetCommentsForPost(postId primitive.ObjectID) ([]*domain.Comment, error) {
	post, err := store.Get(postId)
	if err != nil {
		panic(fmt.Errorf("Invalid post"))
	}
	return post.Comments, nil
}

func (store *PostMongoDBStore) DeleteAll() {
	store.posts.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *PostMongoDBStore) filter(filter interface{}) ([]*domain.Post, error) {
	cursor, err := store.posts.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *PostMongoDBStore) filterOne(filter interface{}) (post *domain.Post, err error) {
	result := store.posts.FindOne(context.TODO(), filter)
	err = result.Decode(&post)
	return
}

func decode(cursor *mongo.Cursor) (posts []*domain.Post, err error) {
	for cursor.Next(context.TODO()) {
		var post domain.Post
		err = cursor.Decode(&post)
		if err != nil {
			return
		}
		posts = append(posts, &post)
	}
	err = cursor.Err()
	return
}
