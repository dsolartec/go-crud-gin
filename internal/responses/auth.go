package responses

type LogInResponse struct {
	AccessToken string `json:"access_token"`
}

type SignUpResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int    `json:"user_id"`
}
