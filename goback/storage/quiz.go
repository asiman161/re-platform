package storage

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/asiman161/re-platform/app/models"
	"github.com/pkg/errors"
)

func (s Storage) CreateQuiz(ctx context.Context, quiz models.Quiz) (models.Quiz, error) {
	now := time.Now()
	quiz.UpdatedAt = now
	quiz.CreatedAt = now
	quiz.IsOpen = true
	insertQuery, args, _ := sq.Insert(quizzesTable).Columns(models.QuizColumns[1:]...).
		Values(quiz.RoomID, quiz.Author, quiz.Name, quiz.Content, quiz.Variants, quiz.Answers, quiz.IsOpen, quiz.CreatedAt, quiz.UpdatedAt).
		Suffix(suffixReturning(models.QuizColumns)).
		PlaceholderFormat(sq.Dollar).ToSql()

	newQuiz := models.Quiz{}
	err := s.db.GetContext(ctx, &newQuiz, insertQuery, args...)
	if err != nil {
		return models.Quiz{}, errors.Wrap(err, "[store.CreateQuiz] can't insert new quiz")
	}

	bts := models.MakeRdMessage("quiz", &newQuiz)
	_, err = s.rd.Publish(ctx, redisRoomID(newQuiz.RoomID), string(bts)).Result()
	if err != nil {
		return models.Quiz{}, errors.Wrap(err, "[store.CreateQuiz] can't publish chat message to redis")
	}

	return newQuiz, nil
}

func (s Storage) CloseQuiz(ctx context.Context, id int, author string) error {
	where := sq.Eq{"id": id, "author": author}
	q, args, _ := sq.Update(quizzesTable).
		Set("is_open", false).Set("updated_at", time.Now()).
		Where(where).
		Suffix(suffixReturning(models.QuizColumns)).
		PlaceholderFormat(sq.Dollar).ToSql()

	quiz := models.Quiz{}
	err := s.db.GetContext(ctx, &quiz, q, args...)
	if err != nil {
		return errors.Wrap(err, "can't close quiz")
	}

	bts := models.MakeRdMessageStr("close_quiz", "someone")
	_, err = s.rd.Publish(ctx, redisRoomID(quiz.RoomID), string(bts)).Result()
	if err != nil {
		return errors.Wrap(err, "[store.CreateQuiz] can't publish chat message to redis")
	}

	return nil
}

func (s Storage) AnswerQuiz(ctx context.Context, quizID int, answer models.Answer) error {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return errors.Wrap(err, "[storage.AnswerQuiz] can't init tx")
	}
	defer tx.Rollback()

	where := sq.Eq{"id": quizID, "is_open": true}

	getQuery, args, _ := sq.Select(models.QuizColumns...).From(quizzesTable).Where(where).PlaceholderFormat(sq.Dollar).ToSql()

	quiz := models.Quiz{}
	err = tx.GetContext(ctx, &quiz, getQuery, args...)
	if err != nil {
		return errors.Wrap(err, "[storage.AnswerQuiz] can't get answer")
	}

	quiz.Answers = append(quiz.Answers, answer)
	quiz.UpdatedAt = time.Now()

	updateQuery, args, _ := sq.Update(quizzesTable).
		Set("answers", quiz.Answers).
		Set("updated_at", quiz.UpdatedAt).
		Where(where).
		PlaceholderFormat(sq.Dollar).ToSql()

	res, err := tx.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return errors.Wrap(err, "[storage.AnswerQuiz] can't update answer")
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return errors.Wrap(sql.ErrNoRows, "[storage.AnswerQuiz] not found quiz during update")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "[storage.AnswerQuiz] can't commit tx")
	}

	bts := models.MakeRdMessageStr("new_answer", "someone")
	_, err = s.rd.Publish(ctx, redisRoomID(quiz.RoomID), string(bts)).Result()
	if err != nil {
		return errors.Wrap(err, "[store.CreateQuiz] can't publish chat message to redis")
	}

	return nil
}

func (s Storage) GetQuizzes(ctx context.Context, roomID string, onlyOpen bool) ([]models.Quiz, error) {
	where := sq.Eq{"room_id": roomID}
	if onlyOpen {
		where["is_open"] = true
	}
	q, args, _ := sq.Select(models.QuizColumns...).From(quizzesTable).
		Where(where).
		OrderBy("created_at").
		PlaceholderFormat(sq.Dollar).ToSql()

	quizzes := make([]models.Quiz, 0)

	err := s.db.SelectContext(ctx, &quizzes, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "can't get quizzes")
	}

	return quizzes, nil
}
