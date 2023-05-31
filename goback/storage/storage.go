package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/asiman161/re-platform/app/models"
)

type Storager interface {
	GetOpenRooms(ctx context.Context) ([]models.Room, error)
	GetRoom(ctx context.Context, roomID string) (models.Room, error)
	CreateRoom(ctx context.Context, name, author string) (models.Room, error)
	CloseRoom(ctx context.Context, roomID, author string) error
	ChangeRoomUserVisibility(ctx context.Context, activity models.RoomUserActivity) error
	GetCurrentRoomUsers(ctx context.Context, roomID string) ([]models.RoomUserActivity, error)

	GetUsers(ctx context.Context) ([]models.User, error)

	WriteChatMessage(ctx context.Context, message models.ChatMessage) (models.ChatMessage, error)
	GetMessages(ctx context.Context, chatID string) ([]models.ChatMessage, error)
	SubscribeMessages(ctx context.Context, chatID string) (chan models.RdMessage, error)

	CreateQuiz(ctx context.Context, quiz models.Quiz) (models.Quiz, error)
	CloseQuiz(ctx context.Context, id int, author string) error
	AnswerQuiz(ctx context.Context, quizID int, answer models.Answer) error
	GetQuizzes(ctx context.Context, roomID string, onlyOpen bool) ([]models.Quiz, error)
}

type Storage struct {
	db *sqlx.DB
	rd *redis.Client
}

func (s Storage) ChangeRoomUserVisibility(ctx context.Context, activity models.RoomUserActivity) error {
	q, args, _ := sq.Insert(activeRoomUsersTable).Columns(models.ActiveRoomUsersColumns[1:]...).
		Values(activity.RoomID, activity.Email, activity.Connected, activity.Active, time.Now()).
		PlaceholderFormat(sq.Dollar).ToSql()
	_, err := s.db.ExecContext(ctx, q, args...)
	if err != nil {
		return errors.Wrap(err, "can't create user activity")
	}

	bts := models.MakeRdMessageStr("change_visibility", "room_id")
	_, err = s.rd.Publish(ctx, redisRoomID(activity.RoomID), bts).Result()
	return err
}

func (s Storage) GetCurrentRoomUsers(ctx context.Context, roomID string) ([]models.RoomUserActivity, error) {
	users := make([]models.RoomUserActivity, 0)

	err := s.db.SelectContext(ctx, &users, sqlGetRoomUsers, roomID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get users for provided room")
	}

	return users, nil
}

func (s Storage) GetRoom(ctx context.Context, roomID string) (models.Room, error) {
	room := models.Room{}
	q, args, _ := sq.Select(models.RoomColumns...).From(roomTable).Where(sq.Eq{"id": roomID}).PlaceholderFormat(sq.Dollar).ToSql()

	err := s.db.GetContext(ctx, &room, q, args...)
	if err != nil {
		return models.Room{}, errors.Wrap(err, "can't get room")
	}

	return room, nil
}

func (s Storage) GetOpenRooms(ctx context.Context) ([]models.Room, error) {
	rooms := make([]models.Room, 0)
	q, args, _ := sq.Select(models.RoomColumns...).From(roomTable).
		Where(sq.Eq{"is_open": true}).
		OrderBy("id").
		PlaceholderFormat(sq.Dollar).ToSql()

	err := s.db.SelectContext(ctx, &rooms, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "can't select rooms")
	}

	return rooms, nil
}

func (s Storage) CreateRoom(ctx context.Context, name, author string) (models.Room, error) {
	room := models.Room{}
	q, args, _ := sq.Insert(roomTable).Columns(models.RoomColumns[1:]...).
		Values(name, author, true, time.Now(), time.Now()).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(models.RoomColumns, ","))).
		PlaceholderFormat(sq.Dollar).ToSql()
	err := s.db.GetContext(ctx, &room, q, args...)
	if err != nil {
		return models.Room{}, errors.Wrap(err, "can't create room")
	}

	return room, nil
}

func (s Storage) CloseRoom(ctx context.Context, roomID, author string) error {
	//id, _ := strconv.Atoi(roomID)
	q, args, _ := sq.Update(roomTable).
		Where(sq.Eq{"id": roomID, "author": author}).
		Set("is_open", false).
		PlaceholderFormat(sq.Dollar).ToSql()
	res, err := s.db.Exec(q, args...)
	if err != nil {
		return errors.Wrap(err, "can't close room")
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	bts := models.MakeRdMessageStr("close_room", "room_id")
	_, err = s.rd.Publish(ctx, redisRoomID(roomID), bts).Result()

	return nil
}

func New(db *sqlx.DB, rd *redis.Client) Storager {
	return &Storage{db: db, rd: rd}
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
