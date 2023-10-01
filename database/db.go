package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Repo          *SQLiteRepository
	SeriesCh      chan func() (string, string)
	AssignmentsCh chan string
	ColorsCh      chan string
	ActionsCh     chan bool
	ErrorsCh      chan func() (string, []any, string)

	DBWriterCh chan ExecIn
	Quit       chan struct{}
)

type SQLiteRepository struct {
	db *sql.DB
}

type ExecIn struct {
	quer string
	vals []any
	ch   chan ExecOut
}

type ExecOut struct {
	res sql.Result
	err error
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

// Handle all database executions in series rather than parallel
func DBWriter() {
	for {
		select {
		case <-Quit:
			return
		case exe := <-DBWriterCh:
			out := ExecOut{}
			out.res, out.err = Repo.db.Exec(exe.quer, exe.vals...)
			exe.ch <- out
			close(exe.ch)
		}
	}
}

// DB initialization
func StartDatabase(loc string) {
	log.Println("Starting database...")
	// Database locking error fix from API spec
	db, err := sql.Open("sqlite3", "file:"+loc+"?_mutex=full&_busy_timeout=9999999")
	if err != nil {
		log.Fatal(err)
	}

	Repo = NewSQLiteRepository(db)

	if err := Repo.Initialize(); err != nil {
		log.Fatal(err)
	}

	DBWriterCh = make(chan ExecIn)
	go DBWriter()
}

func RegisterChannels(serc chan func() (string, string), assc chan string, colc chan string, actc chan bool, errc chan func() (string, []any, string), quit chan struct{}) {
	SeriesCh = serc
	AssignmentsCh = assc
	ColorsCh = colc
	ActionsCh = actc
	ErrorsCh = errc
	Quit = quit
}
