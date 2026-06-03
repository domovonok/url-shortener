package url

type createRequest struct {
	Url string `json:"url"`
}

type createResponse struct {
	Code string `json:"code"`
}

type getResponse struct {
	Url string `json:"url"`
}

type errorResponse struct {
	Error string `json:"error"`
}
