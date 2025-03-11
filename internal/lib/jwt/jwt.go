package jwt

import (
	"time"

	"grpc-service-ref/internal/domain/models"
  
  "github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    
    claims := token.Claims.(jwt.MapClaims)
    claims["uid"] = user.ID
    claims["exp"] = time.Now().Add(duration).Unix()

    // TODO: secret
    tokenString, err := token.SignedString([]byte("secret"))
    if err != nil {
        return "", err
    }
    
    return tokenString, nil
}
