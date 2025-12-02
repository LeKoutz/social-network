package forum

import "time"

var (
	ResponseStructMock = ResponseStruct{
		WebsiteName: "Forum",
		User: User{
			Name: "Guest",
		},
		Posts: []Post{
			{
				Id:        1,
				Title:     "something",
				Body:      "mpla mpla",
				Timestamp: time.Now().UTC(),
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
				Timestamp: time.Now().UTC(),
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
			True:    true,
			Message: "mpla",
		},
	}
)

func ValuesToClient() ResponseStruct {
	return ResponseStructMock
}
