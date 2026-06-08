// Package auth provides password hashing (bcrypt) and signed session tokens (JWT).
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Roles.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// ErrInvalidToken is returned when a session token fails verification.
var ErrInvalidToken = errors.New("invalid token")

// Service issues and verifies JWT session tokens.
type Service struct {
	secret []byte
	ttl    time.Duration
}

// New constructs an auth service. ttl<=0 defaults to 24h.
func New(secret string, ttl time.Duration) *Service {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &Service{secret: []byte(secret), ttl: ttl}
}

// TTL is the token lifetime (used to set the cookie max-age).
func (s *Service) TTL() time.Duration { return s.ttl }

// HashPassword returns a bcrypt hash of pw.
func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword reports whether pw matches the stored bcrypt hash.
func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

// Claims is the JWT payload.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Issue signs a token for the user, returning the token and its expiry.
func (s *Service) Issue(userID, username, role string) (string, time.Time, error) {
	exp := time.Now().Add(s.ttl)
	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

// Parse verifies a token and returns its claims.
func (s *Service) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
