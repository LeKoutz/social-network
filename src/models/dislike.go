package models

type Dislike struct {
	PostId    int64
	UserId    int64
	CommentId int64
	Timestamp int64
	TimestampString string
}
