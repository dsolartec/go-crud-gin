package authenticator

type AuthenticatorToken struct {
	UserID      int
	Permissions []string
}

type Authenticator interface {
	GetToken(data AuthenticatorToken) (string, error)
	Authenticate(token string, permissions []string) (*AuthenticatorToken, error)
}
