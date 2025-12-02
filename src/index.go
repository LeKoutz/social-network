package forum

import "net/http"

func showIndex(res http.ResponseWriter, _ *http.Request, user User) {
	data := ResponseStructMock
	data.User = user
	respondView(res, "index_view", data)
}
