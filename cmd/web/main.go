package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger *slog.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})) // NewTextHandler also works

	app := &application{
		logger: logger,
	}

	logger.Info("Starting server on", slog.Any("addr", cfg.addr))

	err := http.ListenAndServe(cfg.addr, app.routes())
	logger.Error(err.Error()) // no logger.Fatal()
	os.Exit(1)
}
