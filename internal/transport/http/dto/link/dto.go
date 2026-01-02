package link

type CreateRequest struct {
	Url string `json:"url"`
}

type GetRequest struct {
	Code string `json:"code"`
}
