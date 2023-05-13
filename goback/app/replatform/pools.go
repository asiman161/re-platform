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

func (i *Implementation) CreatePool(w http.ResponseWriter, r *http.Request) {
	pool := models.Pool{}

	err := json.NewDecoder(r.Body).Decode(&pool)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "can't decode pool")
		return
	}
	author := extractAuthor(r)

	pool.Author = author
	pool, err = i.store.CreatePool(r.Context(), pool)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "can't create pool")
		return
	}

	render.JSON(w, r, pool)
}

func (i *Implementation) AnswerPool(w http.ResponseWriter, r *http.Request) {
	poolID, err := strconv.Atoi(chi.URLParam(r, "pool_ID"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "can't parse pool ID")
		return
	}

	req := poolAnswerReq{}

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

	err = i.store.AnswerPool(r.Context(), poolID, answer)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "can't update pool with answer")
		return
	}

	render.PlainText(w, r, "OK")
}

func (i *Implementation) ClosePool(w http.ResponseWriter, r *http.Request) {
	poolID, err := strconv.Atoi(chi.URLParam(r, "pool_ID"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "can't parse pool ID")
		return
	}

	err = i.store.ClosePool(r.Context(), poolID)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("can't close pool: %d, err: %v", poolID, err.Error()))
		return
	}

	render.PlainText(w, r, "OK")
}

func (i *Implementation) GetPools(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_ID")
	r.URL.Query().Get("is_open")

	onlyOpen, _ := strconv.ParseBool(r.URL.Query().Get("is_open"))

	pools, err := i.store.GetPools(r.Context(), roomID, onlyOpen)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Sprintf("can't get pools by room: %s", roomID))
		return
	}

	render.JSON(w, r, pools)
}
