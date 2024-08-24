package auth

import (
	"auth/config"
	pb "auth/genproto/user"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GeneratedAccessJWTToken(req *pb.LoginRes) error {
	conf := config.Load()
	token := *jwt.New(jwt.SigningMethodHS256)

	//payload
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = req.Id
	claims["role"] = req.Role
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(30 * time.Minute).Unix()

	newToken, err := token.SignedString([]byte(conf.Token.ACCES_KEY))
	if err != nil {
		log.Println(err)
		return err
	}

	req.Accestoken = newToken
	return nil
}

func ValidateAccesToken(tokenStr string) (bool, error) {
	_, err := ExtractAccesClaim(tokenStr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractAccesClaim(tokenStr string) (*jwt.MapClaims, error) {
	conf := config.Load()
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(conf.Token.ACCES_KEY), nil
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

func GetUserIdFromAccesToken(req *pb.LoginRes) error {
	conf := config.Load()
	AccesToken, err := jwt.Parse(req.Accestoken, func(token *jwt.Token) (interface{}, error) { return []byte(conf.Token.ACCES_KEY), nil })
	if err != nil || !AccesToken.Valid {
		return err
	}
	claims, ok := AccesToken.Claims.(jwt.MapClaims)
	if !ok {
		return err
	}
	req.Id = claims["user_id"].(string)
	req.Role = claims["role"].(string)

	return nil
}
