package auth

import (
	"auth/config"
	pb "auth/genproto/user"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GeneratedRefreshJWTToken(req *pb.LoginRes) error {
	conf := config.Load()
	token := *jwt.New(jwt.SigningMethodHS256)
	//payload
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = req.Id
	claims["role"] = req.Role
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().AddDate(0, 6, 0).Unix()

	newToken, err := token.SignedString([]byte(conf.Token.REFRESH_KEY))
	if err != nil {
		return err
	}

	req.Refreshtoken = newToken
	return nil
}

func ValidateRefreshToken(tokenStr string) (bool, error) {
	_, err := ExtractRefreshClaim(tokenStr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractRefreshClaim(tokenStr string) (*jwt.MapClaims, error) {
	conf := config.Load()
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(conf.Token.REFRESH_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, err
	}

	return &claims, nil
}

func GetUserIdFromRefreshToken(req *pb.LoginRes) error {
	conf := config.Load()
	refreshToken, err := jwt.Parse(req.Refreshtoken, func(token *jwt.Token) (interface{}, error) { return []byte(conf.Token.REFRESH_KEY), nil })
	if err != nil || !refreshToken.Valid {
		return err
	}
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return err
	}
	req.Id = claims["user_id"].(string)
	req.Role = claims["role"].(string)

	return nil
}
