package infrastructure

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Ensure MyJwtService implements the domain.JwtService interface
var _ domain.JwtService = (*MyJwtService)(nil)

type MyJwtService struct {
	secretKey string
}

func NewJwtService(secretKey string) *MyJwtService {
	return &MyJwtService{secretKey: secretKey}
}

func (s *MyJwtService) GetSignedToken(c context.Context, user *domain.User) (string, error) {
	// Prepare claims
	claims := domain.Claims{
		UserId:   user.Id.Hex(),
		Role:     user.Role,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "task-manager-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("jwt service: failed to sign token: %w", err)
	}
	return signedToken, nil
}

func (s *MyJwtService) ParseToken(c context.Context, tokenString string) (*domain.Claims, error) {
	claims := &domain.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil // Return the secret key for verification
	})

	if err != nil {
		var ve *jwt.ValidationError // Check for specific JWT validation errors
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token expired")
			}
			// Other validation errors (e.g., malformed, signature invalid)
			return nil, errors.New("invalid token")
		}

		return nil, fmt.Errorf("jwt service: failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
