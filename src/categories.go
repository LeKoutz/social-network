package forum

import (
	"net/http"
)

type Category struct {
	Id   int
	Name string
}

type Categories []Category

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

func showCategories(res http.ResponseWriter, _ *http.Request, user User) {
	categories, err := getAllCategories()
	if err != nil {
		var e Error
		e = e.Consume(err)
		e.LogError()
		e.RespondError(res)
		return
	}
	data := ResponseStruct{
		WebsiteName: "Forum",
		User:        user,
		Categories:  categories,
	}
	respondView(res, "categories_view", data)
}
