package forum

import "time"

var (
	ResponseStructMock = ResponseStruct{
		WebsiteName: "Forum",
		User: GuestUser,
		Posts: Posts{
			{
				Id:        1,
				Title:     "something",
				Body:      "mpla mpla",
				Timestamp: time.Now().Unix(),
				Likes:     2,
				Dislikes:  1,
				Category: Category{
					Id:   1,
					Name: "various",
				},
				Comments: ReturnMockComments(),
			},
			{
				Id:        2,
				Title:     "something else",
				Body:      "even more mpla mpla",
				Timestamp: time.Now().Unix(),
				Likes:     1,
				Liked:     true,
				Dislikes:  0,
				Disliked:  false,
				Category: Category{
					Id:   1,
					Name: "various",
				},
				Comments: ReturnMockComments(),
			},
		},
		Categories: ReturnMockCategories(),
		Error: Error{
			Has:     false,
			Message: "mpla",
		},
	}
)

func ReturnMockResponse() ResponseStruct {
	return ResponseStructMock
}

func ReturnMockCategories() Categories {
	return Categories{
		{
			Id:   1,
			Name: "various",
		},
		{
			Id:   2,
			Name: "general",
		},
	}
}
