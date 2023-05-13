package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Repo *SQLiteRepository
	M    sync.Mutex
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

// DB initialization
func StartDatabase(loc string) {
	log.Println("Starting database...")
	db, err := sql.Open("sqlite3", loc)
	if err != nil {
		log.Fatal(err)
	}

	Repo = NewSQLiteRepository(db)

	if err := Repo.Initialize(); err != nil {
		log.Fatal(err)
	}
}
