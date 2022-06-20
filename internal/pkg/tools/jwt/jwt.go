package jwt

import (
	"mogong/global"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

// MySecret 定义JWT TOKEN的加密盐
var MySecret = []byte(global.MySignedKey)

type MyClaims struct {
	UserName string `json:"user_name"`
	UserId   int64  `json:"user_id"`
	jwt.StandardClaims
}

func CreateToken(username string, userId int64) (string, int64, error) {
	expiresTime := time.Now().Add(7 * 24 * time.Hour).Unix()
	claims := &MyClaims{
		username,
		userId,
		jwt.StandardClaims{ExpiresAt: expiresTime},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(MySecret)
	if err != nil {
		return "", expiresTime, err
	}
	return token, expiresTime, nil
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
