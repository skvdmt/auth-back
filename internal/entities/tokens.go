package entities

// Tokens Токены.
type Tokens struct {
	// Токен доступа.
	AccessToken string `json:"access_token"`
	// Токен обновления.
	refreshToken string
}

// NewTokens Конструктор токенов.
func NewTokens(access, refresh string) *Tokens {
	return &Tokens{
		AccessToken:  access,
		refreshToken: refresh,
	}
}

// RefreshToken Токен обновления.
func (t *Tokens) RefreshToken() string {
	return t.refreshToken
}
