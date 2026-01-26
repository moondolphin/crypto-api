package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	Secret []byte
	TTL    time.Duration
	Issuer string
	Now    func() time.Time
}

func NewJWTService(secret string, ttl time.Duration, issuer string) JWTService {
	return JWTService{
		Secret: []byte(secret),
		TTL:    ttl,
		Issuer: issuer,
	}
}

func (s JWTService) Generate(userID int64, email string) (string, error) {
	now := s.Now
	if now == nil {
		now = time.Now
	}

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iss":   s.Issuer,
		"iat":   now().Unix(),
		"exp":   now().Add(s.TTL).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.Secret)
}
