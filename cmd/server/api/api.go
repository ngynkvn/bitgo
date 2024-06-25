package api

import (
	"bitgo/cmd/server/messages"
	"fmt"
)

type API struct {
}

func NewAPI() API {
	return API{}
}

func (api *API) Receive(jr messages.JsonRequest) messages.JsonResponse {
	switch jr.Method {
	case "version": return api.VersionResponse(jr)
	case "status": return api.StatusResponse(jr)
	default: return MethodNotImplemented(jr.Method)
	}
}


func MethodNotImplemented(method string) messages.JsonResponse {
	return messages.JsonResponse {
		Error: fmt.Sprintf(`method not implemented: "%s"`, method),
	}

}

func (api *API) VersionResponse(m messages.JsonRequest) messages.JsonResponse {
	return messages.JsonResponse{
		Result: "bitgo 0.1",
		ID: m.ID,
	}
}

func (api *API) StatusResponse(m messages.JsonRequest) messages.JsonResponse {
	return messages.JsonResponse{
		Result: `{"torrents":[]}`,
		ID: m.ID,
	}
}