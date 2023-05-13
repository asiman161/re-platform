package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/asiman161/re-platform/app/models"
	"github.com/pkg/errors"
)

func (s Storage) CreatePool(ctx context.Context, pool models.Pool) (models.Pool, error) {
	now := time.Now()
	pool.UpdatedAt = now
	pool.CreatedAt = now
	pool.IsOpen = true
	insertQuery, args, _ := sq.Insert(poolsTable).Columns(models.PoolColumns[1:]...).
		Values(pool.RoomID, pool.Author, pool.Content, pool.Variants, pool.Answers, pool.IsOpen, pool.CreatedAt, pool.UpdatedAt).
		Suffix(suffixReturning(models.PoolColumns)).
		PlaceholderFormat(sq.Dollar).ToSql()

	newPool := models.Pool{}
	err := s.db.SelectContext(ctx, &newPool, insertQuery, args...)
	if err != nil {
		return models.Pool{}, errors.Wrap(err, "[store.CreatePool] can't insert new pool")
	}

	//bts := models.MakeRdMessage("pool", newPool)
	bts, _ := json.Marshal(newPool)
	_, err = s.rd.Publish(ctx, redisRoomID(newPool.RoomID), string(bts)).Result()
	if err != nil {
		return models.Pool{}, errors.Wrap(err, "[store.CreatePool] can't publish chat message to redis")
	}

	return newPool, nil
}

func (s Storage) ClosePool(ctx context.Context, id int) error {
	q, args, _ := sq.Update(poolsTable).Set("is_open", false).Where("id", id).
		PlaceholderFormat(sq.Dollar).ToSql()

	result, err := s.db.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.Wrap(sql.ErrNoRows, "can't close pool")
	}

	return nil
}

func (s Storage) AnswerPool(ctx context.Context, poolID int, answer models.Answer) error {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return errors.Wrap(err, "[storage.AnswerPool] can't init tx")
	}
	defer tx.Rollback()

	where := sq.Eq{"id": poolID}

	getQuery, args, _ := sq.Select(models.PoolColumns...).From(poolsTable).Where(where).PlaceholderFormat(sq.Dollar).ToSql()

	pool := models.Pool{}
	err = tx.GetContext(ctx, &pool, getQuery, args...)
	if err != nil {
		return errors.Wrap(err, "[storage.AnswerPool] can't get answer")
	}

	pool.Answers = append(pool.Answers, answer)
	pool.UpdatedAt = time.Now()

	updateQuery, args, _ := sq.Update(poolsTable).
		Set("answers", pool.Answers).
		Set("updated_at", pool.UpdatedAt).
		Where(where).
		PlaceholderFormat(sq.Dollar).ToSql()

	res, err := tx.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return errors.Wrap(err, "[storage.AnswerPool] can't update answer")
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return errors.Wrap(sql.ErrNoRows, "[storage.AnswerPool] not found pool during update")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "[storage.AnswerPool] can't commit tx")
	}

	return nil
}

func (s Storage) GetPools(ctx context.Context, roomID string, onlyOpen bool) ([]models.Pool, error) {
	where := sq.Eq{"room_id": roomID}
	if onlyOpen {
		where["is_open"] = true
	}
	q, args, _ := sq.Select(models.PoolColumns...).From(poolsTable).Where(where).PlaceholderFormat(sq.Dollar).ToSql()

	pools := make([]models.Pool, 0)

	err := s.db.SelectContext(ctx, &pools, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "can't get pools")
	}

	return pools, nil
}
