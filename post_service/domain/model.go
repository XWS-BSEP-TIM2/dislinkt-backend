package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Ownable interface {
	GetOwnerId() primitive.ObjectID
}

type Post struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	OwnerId      primitive.ObjectID `bson:"ownerId"`
	CreationTime time.Time          `bson:"creationTime"`
	Content      string             `bson:"content"`
	Image        primitive.Binary   `bson:"image"`
	Links        []string           `bson:"links"`
	Comments     []*Comment         `bson:"comments"`
	Likes        []*Like            `bson:"likes"`
	Dislikes     []*Dislike         `bson:"dislikes"`
}

func (p *Post) GetOwnerId() primitive.ObjectID {
	return p.OwnerId
}

type Comment struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	OwnerId      primitive.ObjectID `bson:"ownerId"`
	CreationTime time.Time          `bson:"creationTime"`
	Content      string             `bson:"content"`
}

func (c *Comment) GetOwnerId() primitive.ObjectID {
	return c.OwnerId
}

type Like struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	OwnerId      primitive.ObjectID `bson:"ownerId"`
	CreationTime time.Time          `bson:"creationTime"`
}

func (l *Like) GetOwnerId() primitive.ObjectID {
	return l.OwnerId
}

type Dislike struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	OwnerId      primitive.ObjectID `bson:"ownerId"`
	CreationTime time.Time          `bson:"creationTime"`
}

func (d *Dislike) GetOwnerId() primitive.ObjectID {
	return d.OwnerId
}

// Owner not serializing
type Owner struct {
	UserId string
	//Username string
	//Name     string
	//Surname  string
}
