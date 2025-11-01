package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/eniehack/kb1/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	port = 3000
)

func main() {
	db, err := sql.Open("sqlite3", "./test.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	dbx := sqlx.NewDb(db, "sqlite3")
	h := &handler.Handler{
		Db: dbx,
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Post("/api/v1/bookmarks/new", h.CreateBookmark)
	r.Get("/api/v1/bookmarks/{bookmarkId}", h.ReadBookmark)
	log.Printf("server started (port: %d)\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatalf("server failed: %s\n", err)
	}
}
