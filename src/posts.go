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

func returnMockPost(post_id int) Posts {
	return Posts{
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
	}
}

func showPost(res http.ResponseWriter, id string, user User) {
	post_id, err := strconv.Atoi(id)
	data := &ResponseStruct{}
	data.Init()
	data.SetUser(user)
	if err != nil {
		e := &Error{}
		errC := e.Consume(err)
		errC.LogError()
		data.SetView("error_view")
		data.SetError(errC)
		data.WriteResponse(res)
		return
	}
	data.SetView("post_view")
	data.SetPosts(returnMockPost(post_id))
	data.WriteResponse(res)
}

func showPosts(res http.ResponseWriter, req *http.Request, user User) {
	query := req.URL.Query()
	_, ok := query["id"]
	if ok {
		log.Printf("%v", query["id"])
		showPost(res, query["id"][0], user)
		return
	}
	data := ReturnMockResponse()
	data.SetUser(user)
	data.SetView("posts_view")
	data.WriteResponse(res)
}
