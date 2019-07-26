package config

import (
	"os"
)

func Get(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(key + " missing from .env file")
	}
	return value
}

func CheckConfig() {
	_, ok := os.LookupEnv("IDCS_TOKEN_URL")
	if !ok {
		panic("IDCS_TOKEN_URL missing from .env file")
	}

	_, ok = os.LookupEnv("IDCS_SIGNING_KEY_URL")
	if !ok {
		panic("IDCS_SIGNING_KEY_URL missing from .env file")
	}

	_, ok = os.LookupEnv("IDCS_AUTH_SECRET")
	if !ok {
		panic("IDCS_AUTH_SECRET missing from .env file")
	}
}
