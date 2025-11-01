package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/eniehack/kb1/internal/id"
	"github.com/go-chi/chi/v5"
)

type CreateBookmarkRequestPayload struct {
	Content   string   `json:"content"`
	Url       string   `json:"url"`
	Tags      []string `json:"tags"`
	ReadLater bool     `json:"read_later"`
	IsPublic  bool     `json:"is_public"`
}

type CreateBookmarkResponsePayload struct {
	Id string `json:"id"`
}

func (h *Handler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	payload := new(CreateBookmarkRequestPayload)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	tx, err := h.Db.BeginTx(r.Context(), nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tagStmt, err := tx.Prepare("INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	now := time.Now()
	noteId := id.GenerateID(payload.Content+" "+strconv.FormatInt(now.UnixMilli(), 16), 0)
	tx.Exec("INSERT INTO notes (id, content, url, is_public, created_at) VALUES (?, ?, ?, ?, ?)", noteId, payload.Content, payload.Url, payload.IsPublic, now)
	if 0 < len(payload.Tags) {
		for _, tag := range payload.Tags {
			tagStmt.ExecContext(r.Context(), noteId, tag)
		}
	}
	if payload.ReadLater {
		tx.Exec("INSERT INTO to_read (note_id) VALUES (?);", noteId)
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resPayload := new(CreateBookmarkResponsePayload)
	resPayload.Id = noteId
	if err := json.NewEncoder(w).Encode(resPayload); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type Bookmark struct {
	Id        string         `db:"id"`
	Content   string         `db:"content"`
	Url       sql.NullString `db:"url"`
	IsPublic  bool           `db:"is_public"`
	CreatedAt time.Time      `db:"created_at"`
}

type ReadBookmarkResponsePayload struct {
	Id        string  `json:"id"`
	Content   string  `json:"content"`
	Url       *string `json:"url"`
	IsPublic  bool    `json:"is_public"`
	CreatedAt string  `json:"created_at"`
}

func (h *Handler) ReadBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkId := chi.URLParam(r, "bookmarkId")
	if !id.VerifyID(bookmarkId) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bookmark := new(Bookmark)
	err := h.Db.GetContext(r.Context(), bookmark, "SELECT * FROM notes WHERE id = ?;", bookmarkId)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resPayload := &ReadBookmarkResponsePayload{
		Id:        bookmark.Id,
		Content:   bookmark.Content,
		IsPublic:  bookmark.IsPublic,
		CreatedAt: bookmark.CreatedAt.Format(time.RFC3339),
	}
	if bookmark.Url.Valid {
		resPayload.Url = &bookmark.Url.String
	} else {
		resPayload.Url = nil
	}
	if err := json.NewEncoder(w).Encode(resPayload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
