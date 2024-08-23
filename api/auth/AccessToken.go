package auth

import (
	pb "auth/genproto/user"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

const (
	signingkey = "visca barsa"
)

func GeneratedAccessJWTToken(req *pb.LoginRes) error {
	token := *jwt.New(jwt.SigningMethodHS256)

	//payload
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = req.Id
	claims["role"] = req.Role
	claims["username"] = req.Username
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(30 * time.Minute).Unix()

	newToken, err := token.SignedString([]byte(signingkey))
	if err != nil {
		log.Println(err)
		return err
	}

	req.Accestoken = newToken
	return nil
}
