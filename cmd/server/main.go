package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/eniehack/kb1/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Port     int    `toml:"port"`
	Database string `toml:"database"`
}

func main() {
	configFilePath := flag.String("config", "./config.toml", "")
	f, err := os.Open(*configFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	config := new(Config)
	if _, err := toml.NewDecoder(f).Decode(config); err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("sqlite3", config.Database)
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
	r.Get("/api/v1/bookmarks/{bookmarkId}", h.ReadBookmark)
	log.Printf("server started (port: %d)\n", config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r); err != nil {
		log.Fatalf("server failed: %s\n", err)
	}
}
