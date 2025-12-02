package forum

type Comment struct {
	Id       int
	Body     string
	Likes    int
	Liked    bool
	Dislikes int
	Disliked bool
}

type Comments []Comment

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
