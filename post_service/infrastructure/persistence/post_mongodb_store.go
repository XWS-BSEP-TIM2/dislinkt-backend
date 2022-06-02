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

func (store *PostMongoDBStore) GetPostsFromUser(userId primitive.ObjectID) ([]*domain.Post, error) {
	filter := bson.M{"ownerId": userId}
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

func (store *PostMongoDBStore) InsertLike(postId primitive.ObjectID, like *domain.Like) error {
	post, err := store.Get(postId)
	if err != nil {
		return err
	}

	for _, existingLike := range post.Likes {
		if like.OwnerId == existingLike.OwnerId {
			panic(errors.NewInvalidArgumentError("User cannot give like multiple times"))
		}
	}

	like.Id = primitive.NewObjectID()
	post.Likes = append(post.Likes, like)

	err = store.Update(post)
	if err != nil {
		return err
	}

	for _, existingDislike := range post.Dislikes {
		if like.OwnerId == existingDislike.OwnerId {
			store.UndoDislike(postId, existingDislike.Id)
			break
		}
	}

	return nil
}

func (store *PostMongoDBStore) GetLike(postId primitive.ObjectID, likeId primitive.ObjectID) (*domain.Like, error) {
	// should be fetched by mongo syntax
	post, err := store.Get(postId)
	if err != nil {
		panic(fmt.Errorf("Invalid post"))
	}
	for _, like := range post.Likes {
		if like.Id == likeId {
			return like, nil
		}
	}
	panic(errors.NewEntityNotFoundError("Like with given id not found."))
}

func (store *PostMongoDBStore) GetLikesForPost(postId primitive.ObjectID) ([]*domain.Like, error) {
	post, err := store.Get(postId)
	if err != nil {
		panic(fmt.Errorf("Invalid post"))
	}
	return post.Likes, nil
}

func (store *PostMongoDBStore) UndoLike(postId primitive.ObjectID, reactionId primitive.ObjectID) error {
	post, err := store.Get(postId)
	if err != nil {
		return err
	}

	index := -1
	for i, existingLike := range post.Likes {
		if reactionId == existingLike.Id {
			index = i
			break

		}
	}
	if index == -1 {
		panic(errors.NewEntityNotFoundError("Like with given id not found."))
	}

	post.Likes = RemoveIndexLike(post.Likes, index)

	err = store.Update(post)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostMongoDBStore) InsertDislike(postId primitive.ObjectID, dislike *domain.Dislike) error {
	post, err := store.Get(postId)
	if err != nil {
		return err
	}

	for _, existingDislike := range post.Dislikes {
		if dislike.OwnerId == existingDislike.OwnerId {
			panic(errors.NewInvalidArgumentError("User cannot give dislike multiple times"))
		}
	}
	dislike.Id = primitive.NewObjectID()
	post.Dislikes = append(post.Dislikes, dislike)

	err = store.Update(post)
	if err != nil {
		return err
	}

	for _, existingLike := range post.Likes {
		if dislike.OwnerId == existingLike.OwnerId {
			store.UndoLike(postId, existingLike.Id)
			break
		}
	}

	return nil
}

func (store *PostMongoDBStore) GetDislike(postId primitive.ObjectID, dislikeId primitive.ObjectID) (*domain.Dislike, error) {
	post, err := store.Get(postId)
	if err != nil {
		panic(fmt.Errorf("Invalid post"))
	}
	for _, dislike := range post.Dislikes {
		if dislike.Id == dislikeId {
			return dislike, nil
		}
	}
	panic(errors.NewEntityNotFoundError("Dislike with given id not found."))
}

func (store *PostMongoDBStore) GetDislikesForPost(postId primitive.ObjectID) ([]*domain.Dislike, error) {
	post, err := store.Get(postId)
	if err != nil {
		panic(fmt.Errorf("Invalid post"))
	}
	return post.Dislikes, nil
}

func (store *PostMongoDBStore) UndoDislike(postId primitive.ObjectID, reactionId primitive.ObjectID) error {
	post, err := store.Get(postId)
	if err != nil {
		return err
	}

	index := -1
	for i, existingDislike := range post.Dislikes {
		if reactionId == existingDislike.Id {
			index = i
			break

		}
	}
	if index == -1 {
		panic(errors.NewEntityNotFoundError("Dislike with given id not found."))
	}

	post.Dislikes = RemoveIndexDislike(post.Dislikes, index)

	err = store.Update(post)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostMongoDBStore) GetReactions(postId primitive.ObjectID, userId primitive.ObjectID) *domain.Reactions {
	post, err := store.Get(postId)
	if err != nil {
		panic(errors.NewEntityNotFoundError("No post with given id"))
	}

	for _, existingLike := range post.Likes {
		if userId == existingLike.OwnerId {
			return &domain.Reactions{
				Liked:    true,
				Disliked: false,
			}
		}
	}

	for _, existingDislike := range post.Dislikes {
		if userId == existingDislike.OwnerId {
			return &domain.Reactions{
				Liked:    false,
				Disliked: true,
			}
		}
	}

	return &domain.Reactions{
		Liked:    false,
		Disliked: false,
	}
}

func (store *PostMongoDBStore) GetAllPostsFromIds(ids []*primitive.ObjectID) []*domain.Post {
	var posts []*domain.Post
	for _, id := range ids {
		userPosts, err := store.GetPostsFromUser(*id)
		if err != nil {
			fmt.Printf("Error geting posts from user with id %s", id)
		}
		posts = append(posts, userPosts...)
	}

	return posts
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

func RemoveIndexLike(s []*domain.Like, index int) []*domain.Like {
	return append(s[:index], s[index+1:]...)
}

func RemoveIndexDislike(s []*domain.Dislike, index int) []*domain.Dislike {
	return append(s[:index], s[index+1:]...)
}
