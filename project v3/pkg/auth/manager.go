package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/erkkke/golang-start/project/internal/models"
	"math/rand"
	"time"
)


type TokenManger interface {
	NewJWT(userInfo *models.AuthorizedUserInfo, ttl time.Duration) (string, error)
	Parse(accessToken string) (*models.AuthorizedUserInfo, error)
	NewRefreshToken() (string, error)
}

type (
	Manager struct {
		signingKey string
	}

	tokenClaims struct {
		jwt.StandardClaims
		User *models.AuthorizedUserInfo `json:"user"`
	}
)

func NewManager(signingKey string) (TokenManger, error) {
	if signingKey == "" {
		return nil, errors.New("error: empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userInfo *models.AuthorizedUserInfo, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		User: userInfo,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (*models.AuthorizedUserInfo, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from token")
	}

	return claims.User, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
