package dto

type ErrResponse struct {
	Msg string `json:"error"`
}

func NewErrResponse(msg string) *ErrResponse {
	return &ErrResponse{Msg: msg}
}
