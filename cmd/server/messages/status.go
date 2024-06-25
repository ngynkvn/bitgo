package messages

type JsonRequest struct {
	Method string         `json:"method"`
	Params map[string]any `json:"params"`
	ID     string         `json:"id"`
}

type JsonResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
	ID     string `json:"id"`
}

type ParamsAddTorrent struct {
	File       string `json:"file"`
	OutputPath string `json:"path"`
}
