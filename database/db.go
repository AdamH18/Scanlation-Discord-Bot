package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Repo          *SQLiteRepository
	M             sync.Mutex
	SeriesCh      chan func() (string, string)
	AssignmentsCh chan string
	ColorsCh      chan string
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

func RegisterChannels(serc chan func() (string, string), assc chan string, colc chan string) {
	SeriesCh = serc
	AssignmentsCh = assc
	ColorsCh = colc
}
