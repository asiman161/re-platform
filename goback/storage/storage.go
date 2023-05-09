package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	CreateRoom(ctx context.Context, name, author string) (models.Room, error)
	CloseRoom(ctx context.Context, roomID, author string) error

	GetUsers(ctx context.Context) ([]models.User, error)
	WriteChatMessage(ctx context.Context, message models.ChatMessage) (models.ChatMessage, error)
	GetMessages(ctx context.Context, chatID string) ([]models.ChatMessage, error)
	SubscribeMessages(ctx context.Context, chatID string) (chan models.ChatMessage, error)
}

type Storage struct {
	db *sqlx.DB
	rd *redis.Client
}

func (s Storage) GetOpenRooms(ctx context.Context) ([]models.Room, error) {
	rooms := make([]models.Room, 0)
	q, args, _ := sq.Select(models.RoomColumns...).From(roomTable).Where(sq.Eq{"is_open": true}).PlaceholderFormat(sq.Dollar).ToSql()

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

	return nil
}

func New(db *sqlx.DB, rd *redis.Client) Storager {
	return &Storage{db: db, rd: rd}
}

func (s Storage) WriteChatMessage(ctx context.Context, message models.ChatMessage) (models.ChatMessage, error) {
	message.CreatedAt = time.Now()
	q, args, _ := sq.Insert(chatTable).Columns(models.ChatColumns[1:]...).
		Values(message.ChatID, message.Content, message.Author, message.CreatedAt).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(models.ChatColumns, ","))).
		PlaceholderFormat(sq.Dollar).ToSql()
	err := s.db.GetContext(ctx, &message, q, args...)
	if err != nil {
		return models.ChatMessage{}, errors.Wrap(err, "can't save chat message")
	}

	bts, _ := json.Marshal(message)
	_, err = s.rd.Publish(ctx, redisChatID(message.ChatID), string(bts)).Result()
	if err != nil {
		return models.ChatMessage{}, errors.Wrap(err, "can't publish chat message to redis")
	}
	return message, nil
}

func (s Storage) GetMessages(ctx context.Context, chatID string) ([]models.ChatMessage, error) {
	messages := make([]models.ChatMessage, 0)
	q, args, _ := sq.Select(models.ChatColumns...).
		From(chatTable).
		Where(sq.Eq{"chat_id": chatID}).PlaceholderFormat(sq.Dollar).ToSql()

	err := s.db.SelectContext(ctx, &messages, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "can't select messages")
	}

	return messages, nil
}

func (s Storage) SubscribeMessages(ctx context.Context, chatID string) (chan models.ChatMessage, error) {
	pubsub := s.rd.Subscribe(ctx, redisChatID(chatID))
	rdch := pubsub.Channel()

	msgCh := make(chan models.ChatMessage)

	go func() {
		defer func() {
			pubsub.Close()
			defer close(msgCh)
		}()
		for rdMesg := range rdch {

			msg := models.ChatMessage{}

			err := json.Unmarshal([]byte(rdMesg.Payload), &msg)
			if err != nil {
				log.Println("can't unmarshal message from rd:", err)
				return
			}

			msgCh <- msg
		}
	}()

	return msgCh, nil
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
