package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	// No code used, but we need the init() function to run, so that the driver can register
	// itself with the database/sql package
	"github.com/Bude8/lets-go/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	dsn := flag.String("dsn", "web:pass@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", "MySQL data source name")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})) // NewTextHandler also works

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("Starting server on", slog.Any("addr", cfg.addr))

	err = http.ListenAndServe(cfg.addr, app.routes())
	logger.Error(err.Error()) // no logger.Fatal()
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	// sql.DB is a pool of many connections
	db, err := sql.Open("mysql", dsn) // doesn't create connections, just initialises pool
	if err != nil {
		return nil, err
	}

	err = db.Ping() // Verifies set up is correct by creating connection and checking for errs
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
