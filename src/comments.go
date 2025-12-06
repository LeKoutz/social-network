package forum

import "time"

type Comment struct {
	Id       int
	PostId   int
	UserId   int
	Body     string
	Timestamp time.Time
	Likes    int
	Liked    bool
	Dislikes int
	Disliked bool
	Username string
}

type Comments []Comment

func (c *Comment) validateComment() error {
	if len(c.Body) == 0 {
		return ErrorCommentEmpty
	}
	if len(c.Body) > 1000 {
		return ErrorCommentTooLong
	}
	return nil
}

func ReturnMockComments() Comments {
	return Comments{
		{
			Id:       1,
			Body:     "mpla mpla",
			Likes:    2,
			Disliked: true,
			Dislikes: 1,
		},
	}
}
