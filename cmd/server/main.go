package main

import (
	"bitgo/cmd/server/api"
	"bitgo/cmd/server/unixconn"
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const SocketFilePath = "/tmp/bitgo.sock"

// TODO: probably not the best place to persist application state
const DBFilePath = "/tmp/bitgo.db"

func socketExists(s string) bool {
	_, err := os.Stat(s)
	return !os.IsNotExist(err)
}

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	defer os.Remove(SocketFilePath)
	log.Info().Msg("Starting server")

	db := sqlx.MustOpen("sqlite3", DBFilePath)
	defer func() { os.Remove(DBFilePath); db.Close() }()
	log.Info().Msg("DB connected")
	dbInitialization(db)

	// Init socket listener
	if socketExists(SocketFilePath) {
		log.Error().
			Str("socket", SocketFilePath).
			Msg("already exists, exiting")
		return
	}
	listener, err := net.Listen("unix", SocketFilePath)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("exiting")
	}
	defer listener.Close()
	log.Info().Msg("Server started")
	api := api.NewAPI(db)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Error().
					Err(err).
					Msg("error accepting connection")
				return
			}
			go unixconn.Handle(ctx, api, conn)
		}
	}()
	<-ctx.Done()
}

func dbInitialization(db *sqlx.DB) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS torrents(
        file text UNIQUE,
        path text,
        progress float
    );`)
	if err != nil {
		panic(err.Error())
	}
}
