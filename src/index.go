package forum

import "net/http"

func showIndex(res http.ResponseWriter, _ *http.Request, user User) {
	data := ResponseStruct{}
	data.Init().SetUser(user).SetView("index_view").WriteResponse(res)
}
