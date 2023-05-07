package bootstrap

import (
	"github.com/asiman161/re-platform/app/replatform"
	"github.com/asiman161/re-platform/storage"
	"github.com/pkg/errors"
)

type AppConfig struct {
	DatabaseDSN string `env:"DATABASE_DSN,required"`
	Host        string `env:"HOST" envDefault:"0.0.0.0"`
	Port        int    `env:"PORT" envDefault:"80"`
}

func NewApp(cfg AppConfig) (*replatform.Implementation, error) {
	db, err := InitDB(cfg.DatabaseDSN)
	if err != nil {
		return nil, errors.Wrap(err, "can't init db")
	}

	store := storage.New(db)
	return replatform.New(store), nil
}
