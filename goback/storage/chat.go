package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/asiman161/re-platform/app/models"
	"github.com/pkg/errors"
)

func (s Storage) WriteChatMessage(ctx context.Context, message models.ChatMessage) (models.ChatMessage, error) {
	message.CreatedAt = time.Now()
	q, args, _ := sq.Insert(chatTable).Columns(models.ChatColumns[1:]...).
		Values(message.RoomID, message.Content, message.Author, message.CreatedAt).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(models.ChatColumns, ","))).
		PlaceholderFormat(sq.Dollar).ToSql()
	err := s.db.GetContext(ctx, &message, q, args...)
	if err != nil {
		return models.ChatMessage{}, errors.Wrap(err, "can't save chat message")
	}

	bts := models.MakeRdMessage("message", &message)
	_, err = s.rd.Publish(ctx, redisRoomID(message.RoomID), string(bts)).Result()
	if err != nil {
		return models.ChatMessage{}, errors.Wrap(err, "can't publish chat message to redis")
	}
	return message, nil
}

func (s Storage) GetMessages(ctx context.Context, chatID string) ([]models.ChatMessage, error) {
	messages := make([]models.ChatMessage, 0)
	q, args, _ := sq.Select(models.ChatColumns...).
		From(chatTable).
		Where(sq.Eq{"room_id": chatID}).OrderBy("created_at").PlaceholderFormat(sq.Dollar).ToSql()

	err := s.db.SelectContext(ctx, &messages, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "can't select messages")
	}

	return messages, nil
}

// SubscribeMessages TODO: handle all events
func (s Storage) SubscribeMessages(ctx context.Context, roomID string) (chan models.RdMessage, error) {
	pubsub := s.rd.Subscribe(ctx, redisRoomID(roomID))
	rdch := pubsub.Channel()

	msgCh := make(chan models.RdMessage)

	go func() {
		defer func() {
			_ = pubsub.Close()
			defer close(msgCh)
		}()
		for rdMesg := range rdch {

			msg := models.RdMessage{}

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
