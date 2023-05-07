package replatform

import "github.com/asiman161/re-platform/storage"

type Implementation struct {
	store storage.Storager
}

func New(store storage.Storager) *Implementation {
	return &Implementation{store: store}
}
