package forum

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type Post struct {
	Id        int
	Title     string
	Body      string
	Timestamp time.Time
	Likes     int
	Liked     bool
	Dislikes  int
	Disliked  bool
	Category  Category
	Comments  Comments
}

type Posts []Post

func showPost(id string) ResponseStruct {
	post_id, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return ResponseStruct{
			Error: Error{
				True:    true,
				Message: "id value not numerical",
			},
		}
	}
	return ResponseStruct{
		WebsiteName: "Forum",
		Posts: Posts{
			{
				Id:        post_id,
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
		},
	}
}

func showPosts(res http.ResponseWriter, req *http.Request, user User) {
	query := req.URL.Query()
	_, ok := query["id"]
	if ok {
		log.Printf("%v", query["id"])
		respondView(res, "post_view", showPost(query["id"][0]))
		return
	}
	respondView(res, "posts_view", ValuesToClient())
}
