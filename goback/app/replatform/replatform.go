package replatform

import (
	"net/http"

	"github.com/asiman161/re-platform/storage"
)

type Implementation struct {
	store storage.Storager
}

func New(store storage.Storager) *Implementation {
	return &Implementation{store: store}
}

func extractAuthor(r *http.Request) string {
	return r.Header.Get("Email")
}
