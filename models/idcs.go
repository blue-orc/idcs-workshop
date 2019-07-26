package models

type IdcsSigningKey struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	Kty     string   `json:"kty"`
	X5tS256 string   `json:"x5t#S256"`
	E       string   `json:"e"`
	X5t     string   `json:"x5t"`
	Kid     string   `json:"kid"`
	X5c     []string `json:"x5c"`
	KeyOps  []string `json:"key_ops"`
	Alg     string   `json:"alg"`
	N       string   `json:"n"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
