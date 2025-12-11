package forum

import "net/http"

func showIndex(res http.ResponseWriter, _ *http.Request, user User) {
	data := ResponseStruct{}
	data.Init().SetUser(user)
	categories, err := getAllCategories()
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.SetCategories(categories)
	data.SetView("index_view").WriteResponse(res)
}
