package replatform

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/asiman161/re-platform/app/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (i *Implementation) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	quiz := models.Quiz{}

	err := json.NewDecoder(r.Body).Decode(&quiz)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "can't decode quiz")
		return
	}
	author := extractAuthor(r)

	quiz.Author = author
	quiz, err = i.store.CreateQuiz(r.Context(), quiz)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "can't create quiz")
		return
	}

	render.JSON(w, r, quiz)
}

func (i *Implementation) AnswerQuiz(w http.ResponseWriter, r *http.Request) {
	poolID, err := strconv.Atoi(chi.URLParam(r, "quiz_ID"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "can't parse quiz ID")
		return
	}

	req := quizAnswerReq{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("can't decode answer, err: %v", err.Error()))
		return
	}
	author := extractAuthor(r)

	answer := models.Answer{
		VariantID: req.VariantID,
		Author:    author,
	}

	err = i.store.AnswerQuiz(r.Context(), poolID, answer)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "can't update quiz with answer")
		return
	}

	render.PlainText(w, r, "OK")
}

func (i *Implementation) CloseQuiz(w http.ResponseWriter, r *http.Request) {
	poolID, err := strconv.Atoi(chi.URLParam(r, "quiz_ID"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "can't parse quiz ID")
		return
	}

	author := extractAuthor(r)

	err = i.store.CloseQuiz(r.Context(), poolID, author)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("can't close quiz: %d, err: %v", poolID, err.Error()))
		return
	}

	render.PlainText(w, r, "OK")
}

func (i *Implementation) GetQuizzes(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_ID")
	r.URL.Query().Get("is_open")

	onlyOpen, _ := strconv.ParseBool(r.URL.Query().Get("is_open"))

	pools, err := i.store.GetQuizzes(r.Context(), roomID, onlyOpen)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Sprintf("can't get pools by room: %s", roomID))
		return
	}

	render.JSON(w, r, pools)
}
