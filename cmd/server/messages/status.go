package messages

type JsonRequest struct {
    Method string `json:"method"`
	Params string `json:"params"`
	ID string `json:"id"`
}

type JsonResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
	ID string `json:"id"`
}
