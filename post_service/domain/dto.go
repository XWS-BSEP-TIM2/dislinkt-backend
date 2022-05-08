package domain

type PostDetailsDTO struct {
	Owner       *Owner
	Post        *Post
	ImageBase64 string
	Stats       *Stats
}

type Stats struct {
	CommentsNumber int
	LikesNumber    int
	DislikesNumber int
}
