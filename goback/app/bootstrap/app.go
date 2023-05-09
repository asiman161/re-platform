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
	RedisDSN    string `env:"REDIS_DSN,required"`
}

func NewApp(cfg AppConfig) (*replatform.Implementation, error) {
	db, err := InitDB(cfg.DatabaseDSN)
	if err != nil {
		return nil, errors.Wrap(err, "can't init db")
	}

	rd, err := InitRD(cfg.RedisDSN)
	if err != nil {
		return nil, errors.Wrap(err, "can't init rd")
	}

	store := storage.New(db, rd)
	return replatform.New(store), nil
}
