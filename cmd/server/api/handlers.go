package api

import (
	"bitgo/cmd/server/messages"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/mapstructure"
)

type API struct {
	db *sqlx.DB
}

func NewAPI(db *sqlx.DB) API {
	return API{db}
}

func (api *API) Receive(jr messages.JsonRequest) messages.JsonResponse {
	switch jr.Method {
	case "version": return api.ResponseVersion(jr)
	case "status": return api.ResponseStatus(jr)
	case "add": return api.ResponseAddTorrent(jr)
	default: return MethodNotImplemented(jr)
	}
}


func MethodNotImplemented(m messages.JsonRequest) messages.JsonResponse {
	return messages.JsonResponse {
		Error: fmt.Sprintf(`method not implemented: "%s"`, m.Method),
	}
}

func MethodOK(m messages.JsonRequest) messages.JsonResponse {
	return messages.JsonResponse {
		Result: "OK",
	}
}


func (api *API) ResponseVersion(m messages.JsonRequest) messages.JsonResponse {
	return messages.JsonResponse{
		Result: "bitgo 0.1",
		ID: m.ID,
	}
}

func (api *API) ResponseStatus(m messages.JsonRequest) messages.JsonResponse {
	return messages.JsonResponse{
		Result: `{"torrents":[]}`,
		ID: m.ID,
	}
}

func (api *API) ResponseAddTorrent(m messages.JsonRequest) messages.JsonResponse {
	params := messages.ParamsAddTorrent{}
	err := mapstructure.Decode(m.Params, &params)
	if err != nil {
		return messages.JsonResponse{
			Error: err.Error(),
		}
	}
	err = api.AddTorrent(params)
	if err != nil {
		return messages.JsonResponse{
			Error: err.Error(),
		}
	}
	return MethodOK(m)
}