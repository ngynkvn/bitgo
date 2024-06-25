package unixconn

import (
	"bitgo/cmd/server/api"
	"bitgo/cmd/server/messages"
	"context"
	"encoding/json"
	"net"

	"github.com/rs/zerolog/log"
)

func Handle(ctx context.Context, api api.API, conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	for {
		jr := messages.JsonRequest{}
		err := decoder.Decode(&jr)
		if err != nil {
			log.Err(err).Msg("decode error, exiting")
			return
		}
		log.Info().Any("request", jr).Msg("request received")
		resp := api.Receive(jr)
		err = encoder.Encode(resp)
		if err != nil {
			log.Err(err).Msg("response error")
		}
	}
}
