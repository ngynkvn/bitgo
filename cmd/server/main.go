package main

import (
	"bitgo/cmd/server/api"
	"bitgo/cmd/server/unixconn"
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	defer os.Remove(api.SocketFilePath)

	// TODO(config)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting server")

	db, cleanup := api.StartAppDB()
	defer cleanup()

	// Init socket listener
	if socketExists(api.SocketFilePath) {
		log.Error().
			Str("socket", api.SocketFilePath).
			Msg("already exists, exiting")
		return
	}
	listener, err := net.Listen("unix", api.SocketFilePath)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("exiting")
	}
	defer listener.Close()
	log.Info().Msg("Server started")

	api := api.NewAPI(db)
	go loopAccept(ctx, api, listener)
	log.Info().Msg("Listening")

	<-ctx.Done()
}

func loopAccept(ctx context.Context, api api.API, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().
				Err(err).
				Msg("error accepting connection")
			return
		}
		log.Info().Msg("New client starting handler")
		go unixconn.Handle(ctx, api, conn)
	}
}

func socketExists(s string) bool {
	_, err := os.Stat(s)
	return !os.IsNotExist(err)
}
