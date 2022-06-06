package data

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"path/filepath"
	"time"
)

func InitializePostStore(store domain.PostStore) {
	store.DeleteAll()
	insertPosts(store)
}

func insertPosts(store domain.PostStore) {
	for _, post := range posts {
		err := store.Insert(post)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}

func getIdFromHex(stringId string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(stringId)
	return id
}

func loadImageBytes(imageName string) primitive.Binary {
	absPath, absErr := filepath.Abs(fmt.Sprintf("startup/data/post_images/%s", imageName))
	if absErr != nil {
		log.Fatal(absErr)
	}
	bytes, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatal(err)
	}
	return primitive.Binary{Data: bytes}
}

const ( // iota is reset to 0
	tara  = iota
	sveta = iota
	djole = iota
	rasti = iota
	zare  = iota
)

var userIdMap = map[int]primitive.ObjectID{
	tara:  getIdFromHex("62752bf27407f54ce1839cb7"),
	sveta: getIdFromHex("62752bf27407f54ce1839cb5"),
	djole: getIdFromHex("62752bf27407f54ce1839cb8"),
	rasti: getIdFromHex("62752bf27407f54ce1839cb9"),
	zare:  getIdFromHex("62752bf27407f54ce1839cb6"),
}

const ( // iota is reset to 0
	postId    = "IDPOST"
	commentId = "IDCOMM"
	likeId    = "IDLIKE"
	dislikeId = "IDDISL"
)

func getIdCreator(id string) func() primitive.ObjectID {
	idxMap := map[string]int{
		postId:    0,
		commentId: 0,
		likeId:    0,
		dislikeId: 0,
	}

	return func() primitive.ObjectID {
		idxMap[id] = idxMap[id] + 1
		return getObjectId(fmt.Sprintf("%s%05d", id, idxMap[id]))
	}
}

var getPostId = getIdCreator(postId)
var getCommentId = getIdCreator(commentId)
var getLikeId = getIdCreator(likeId)
var getDislikeId = getIdCreator(dislikeId)

var posts = []*domain.Post{
	{
		Id:           getPostId(),
		OwnerId:      userIdMap[tara],
		CreationTime: time.Date(2022, time.January, 5, 10, 0, 0, 10000000, time.UTC),

		Content: "Go is an amazing language. Can't wait to learn more of it on https://go.dev/tour/welcome/1",
		Image:   loadImageBytes("post1.jpg"),
		Links:   []string{},
		Comments: []*domain.Comment{
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.January, 5, 12, 0, 0, 10000000, time.UTC),
				Content:      "Wow that's awesome!",
			},
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[djole],
				CreationTime: time.Date(2022, time.January, 5, 11, 25, 0, 10000000, time.UTC),
				Content:      "I will check it out too...",
			},
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[zare],
				CreationTime: time.Date(2022, time.January, 6, 12, 25, 0, 10000000, time.UTC),
				Content:      "IKR Go is amazing!",
			},
		},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.January, 5, 12, 5, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[zare],
				CreationTime: time.Date(2022, time.January, 5, 11, 25, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[djole],
				CreationTime: time.Date(2022, time.January, 6, 12, 25, 0, 10000000, time.UTC),
			},
		},
		Dislikes: []*domain.Dislike{
			{
				Id:           getDislikeId(),
				OwnerId:      userIdMap[sveta],
				CreationTime: time.Date(2022, time.January, 5, 12, 45, 0, 10000000, time.UTC),
			},
		},
	},
	{
		Id:           getPostId(),
		OwnerId:      userIdMap[sveta],
		CreationTime: time.Date(2022, time.January, 10, 10, 0, 0, 10000000, time.UTC),

		Content: "I really love this website that teaches you git https://learngitbranching.js.org/",
		Image:   loadImageBytes("post2.png"),
		Links:   []string{},
		Comments: []*domain.Comment{
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.January, 12, 12, 0, 0, 10000000, time.UTC),
				Content:      "Sounds awesome, I'll check it out.",
			},
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.January, 12, 12, 25, 0, 10000000, time.UTC),
				Content:      "This is what I was looking for...",
			},
		},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.January, 12, 11, 25, 0, 10000000, time.UTC),
			},
		},
	},
	{
		Id:           getPostId(),
		OwnerId:      userIdMap[djole],
		CreationTime: time.Date(2022, time.January, 20, 10, 0, 0, 10000000, time.UTC),

		Content: "Follow more content on the topic of Go here https://go.dev/",
		Image:   loadImageBytes("post3.jpg"),
		Links:   []string{},
		Comments: []*domain.Comment{
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.January, 20, 12, 0, 0, 10000000, time.UTC),
				Content:      "Will do!",
			},
		},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[zare],
				CreationTime: time.Date(2022, time.January, 20, 11, 25, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.January, 21, 12, 25, 0, 10000000, time.UTC),
			},
		},
		Dislikes: []*domain.Dislike{
			{
				Id:           getDislikeId(),
				OwnerId:      userIdMap[sveta],
				CreationTime: time.Date(2022, time.January, 20, 12, 45, 0, 10000000, time.UTC),
			},
		},
	},
	{
		Id:           getPostId(),
		OwnerId:      userIdMap[rasti],
		CreationTime: time.Date(2022, time.January, 5, 10, 0, 0, 10000000, time.UTC),

		Content: "Don't you just love programming?",
		Image:   loadImageBytes("post4.jpg"),
		Links:   []string{},
		Comments: []*domain.Comment{
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[zare],
				CreationTime: time.Date(2022, time.February, 5, 12, 0, 0, 10000000, time.UTC),
				Content:      "😂😂😂",
			},
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[djole],
				CreationTime: time.Date(2022, time.February, 5, 11, 25, 0, 10000000, time.UTC),
				Content:      "Right?! Right!?",
			},
		},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.February, 5, 12, 5, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[zare],
				CreationTime: time.Date(2022, time.February, 5, 11, 25, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[djole],
				CreationTime: time.Date(2022, time.February, 6, 12, 25, 0, 10000000, time.UTC),
			},
		},
	},
	{
		Id:           getPostId(),
		OwnerId:      userIdMap[zare],
		CreationTime: time.Date(2022, time.February, 5, 10, 0, 0, 10000000, time.UTC),

		Content:  "These Go tutorials are amazing: https://www.youtube.com/c/NicJackson/playlists",
		Image:    loadImageBytes("post5.jpg"),
		Links:    []string{},
		Comments: []*domain.Comment{},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.February, 12, 12, 5, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[sveta],
				CreationTime: time.Date(2022, time.February, 12, 11, 25, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[djole],
				CreationTime: time.Date(2022, time.February, 14, 12, 25, 0, 10000000, time.UTC),
			},
		},
		Dislikes: []*domain.Dislike{
			{
				Id:           getDislikeId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.February, 12, 12, 45, 0, 10000000, time.UTC),
			},
		},
	},

	// ------------------------------------------------------------------------------------------------------------------------------------------

	{
		Id:           getPostId(),
		OwnerId:      userIdMap[tara],
		CreationTime: time.Date(2022, time.April, 5, 10, 0, 0, 10000000, time.UTC),

		Content: "Feeling proud of my colleagues 🚀",
		Image:   loadImageBytes("post6.jpg"),
		Links:   []string{},
		Comments: []*domain.Comment{
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.April, 5, 12, 0, 0, 10000000, time.UTC),
				Content:      "You Go!!!",
			},
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[djole],
				CreationTime: time.Date(2022, time.April, 5, 11, 25, 0, 10000000, time.UTC),
				Content:      "Awesome stuff...",
			},
		},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.April, 5, 12, 5, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[zare],
				CreationTime: time.Date(2022, time.April, 5, 11, 25, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[djole],
				CreationTime: time.Date(2022, time.April, 6, 12, 25, 0, 10000000, time.UTC),
			},
		},
		Dislikes: []*domain.Dislike{
			{
				Id:           getDislikeId(),
				OwnerId:      userIdMap[sveta],
				CreationTime: time.Date(2022, time.April, 5, 12, 45, 0, 10000000, time.UTC),
			},
		},
	},
	{
		Id:           getPostId(),
		OwnerId:      userIdMap[sveta],
		CreationTime: time.Date(2022, time.April, 10, 10, 0, 0, 10000000, time.UTC),

		Content: "I hate data wrangling",
		Image:   loadImageBytes("post7.png"),
		Links:   []string{},
		Comments: []*domain.Comment{
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.April, 12, 12, 0, 0, 10000000, time.UTC),
				Content:      "IKR...",
			},
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.April, 12, 12, 25, 0, 10000000, time.UTC),
				Content:      "It's the worst",
			},
		},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.April, 12, 11, 25, 0, 10000000, time.UTC),
			},
		},
	},
	{
		Id:           getPostId(),
		OwnerId:      userIdMap[djole],
		CreationTime: time.Date(2022, time.April, 20, 10, 0, 0, 10000000, time.UTC),

		Content: "I am amazed by this article: https://towardsdatascience.com/understanding-boxplots-5e2df7bcbd51",
		Image:   loadImageBytes("post8.png"),
		Links:   []string{},
		Comments: []*domain.Comment{
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.April, 20, 12, 0, 0, 10000000, time.UTC),
				Content:      "I will check it out!",
			},
		},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[zare],
				CreationTime: time.Date(2022, time.April, 20, 11, 25, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.April, 21, 12, 25, 0, 10000000, time.UTC),
			},
		},
		Dislikes: []*domain.Dislike{
			{
				Id:           getDislikeId(),
				OwnerId:      userIdMap[sveta],
				CreationTime: time.Date(2022, time.April, 20, 12, 45, 0, 10000000, time.UTC),
			},
		},
	},
	{
		Id:           getPostId(),
		OwnerId:      userIdMap[rasti],
		CreationTime: time.Date(2022, time.April, 5, 10, 0, 0, 10000000, time.UTC),

		Content: "I love Data Science",
		Image:   loadImageBytes("post9.jpg"),
		Links:   []string{},
		Comments: []*domain.Comment{
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[zare],
				CreationTime: time.Date(2022, time.May, 5, 12, 0, 0, 10000000, time.UTC),
				Content:      "Me 2",
			},
			{
				Id:           getCommentId(),
				OwnerId:      userIdMap[djole],
				CreationTime: time.Date(2022, time.May, 5, 11, 25, 0, 10000000, time.UTC),
				Content:      "It's the best career ever...",
			},
		},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.May, 5, 12, 5, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[zare],
				CreationTime: time.Date(2022, time.May, 5, 11, 25, 0, 10000000, time.UTC),
			},
		},
	},
	{
		Id:           getPostId(),
		OwnerId:      userIdMap[zare],
		CreationTime: time.Date(2022, time.May, 5, 10, 0, 0, 10000000, time.UTC),

		Content:  "I love this meme. Take care of your side effect people",
		Image:    loadImageBytes("post10.png"),
		Links:    []string{},
		Comments: []*domain.Comment{},
		Likes: []*domain.Like{
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[rasti],
				CreationTime: time.Date(2022, time.May, 12, 12, 5, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[sveta],
				CreationTime: time.Date(2022, time.May, 12, 11, 25, 0, 10000000, time.UTC),
			},
			{
				Id:           getLikeId(),
				OwnerId:      userIdMap[djole],
				CreationTime: time.Date(2022, time.May, 14, 12, 25, 0, 10000000, time.UTC),
			},
		},
		Dislikes: []*domain.Dislike{
			{
				Id:           getDislikeId(),
				OwnerId:      userIdMap[tara],
				CreationTime: time.Date(2022, time.May, 12, 12, 45, 0, 10000000, time.UTC),
			},
		},
	},
}
