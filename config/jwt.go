package config

import (
	"errors"
	"strconv"
	"time"
)

func JWTSecret() (string, error) {
	sec := Getenv("JWT_SECRET", "")
	if sec == "" {
		return "", errors.New("JWT_SECRET required")
	}
	return sec, nil
}

func JWTTTL() time.Duration {
	// minutos
	raw := Getenv("JWT_TTL_MINUTES", "60")
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return 60 * time.Minute
	}
	return time.Duration(n) * time.Minute
}
