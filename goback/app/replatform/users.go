package replatform

import (
	"net/http"

	"github.com/go-chi/render"
)

func (i *Implementation) Users(w http.ResponseWriter, r *http.Request) {
	users, err := i.store.GetUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}

	render.JSON(w, r, users)
}
