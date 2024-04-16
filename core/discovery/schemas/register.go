package schemas

type RegisterRequest struct {
	Channels []string `json:"channel"`
}