package main

import (
	"bitgo/cmd/server/unixconn"
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

const SocketFilePath = "/tmp/bitgo.sock"
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
    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                log.Error().
                    Err(err).
                    Msg("error accepting connection")
                return
            }
            go unixconn.Handle(ctx, conn)
        }
    }()
    <-ctx.Done()
}
