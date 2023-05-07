package storage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/asiman161/re-platform/app/models"
)

type Storager interface {
	GetUsers(ctx context.Context) ([]models.User, error)
}

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Storager {
	return &Storage{db: db}
}

func (s Storage) GetUsers(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0)
	q, args, _ := sq.Select(models.UserColumns...).
		From(usersTable).
		PlaceholderFormat(sq.Dollar).ToSql()
	err := s.db.SelectContext(ctx, &users, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "can't ger users")
	}

	return users, nil
}
