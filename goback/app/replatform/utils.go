package replatform

import (
	"net/http"

	"github.com/go-chi/render"
)

func writeError(w http.ResponseWriter, r *http.Request, httpStatus int, text string) {
	w.WriteHeader(httpStatus)
	render.PlainText(w, r, text)
}
