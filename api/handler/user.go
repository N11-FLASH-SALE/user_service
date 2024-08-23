package handler

import (
	"auth/api/auth"
	"auth/api/email"
	pb "auth/genproto/user"
	"auth/storage/redis"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register user
// @Description create new users
// @Tags auth
// @Param info body user.RegisterReq true "User info"
// @Success 200 {object} user.RegisterRes
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error"
// @Router /auth/register [post]
func (h Handler) Register(c *gin.Context) {
	h.Log.Info("Register is starting")
	req := pb.RegisterReq{}
	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.User.Register(c, &req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	h.Log.Info("Register ended")
	c.JSON(http.StatusOK, res)
}

// Login godoc
// @Summary login user
// @Description it generates new access and refresh tokens
// @Tags auth
// @Param userinfo body user.LoginReq true "username and password"
// @Success 200 {object} string "tokens"
// @Failure 400 {object} string "Invalid date"
// @Failure 500 {object} string "error while reading from server"
// @Router /auth/login [post]
func (h Handler) Login(c *gin.Context) {
	h.Log.Info("Login is working")
	req := pb.LoginReq{}

	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.User.Login(c, &req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = auth.GeneratedAccessJWTToken(res)

	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error3": err.Error()})
	}
	err = auth.GeneratedRefreshJWTToken(res)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error4": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"accesToken":   res.Accestoken,
		"refreshToken": res.Refreshtoken,
	})
	h.Log.Info("login is succesfully ended")
}

// UpdatePassword godoc
// @Summary change password
// @Description it change your password
// @Tags auth
// @Param userinfo body user.Password true "username and password"
// @Success 200 {object} string "message"
// @Failure 400 {object} string "Invalid date"
// @Failure 500 {object} string "error while reading from server"
// @Router /auth/change/password [post]
func (h Handler) UpdatePassword(c *gin.Context) {
	h.Log.Info("reset password is working")
	var req pb.Password
	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Code == "" && req.Email != "" {
		code, err := email.Email(req.Email)
		if err != nil {
			h.Log.Error(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
		}
		err = redis.StoreCodes(c, code, req.Email)
		if err != nil {
			h.Log.Error(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "code sent to your email"})
		return
	} else {
		coderes, err := redis.GetCodes(c, req.Email)
		if err != nil {
			h.Log.Error(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if req.Code != coderes {
			h.Log.Info("code is incorrect")
			c.JSON(500, gin.H{"message": "code is not correct"})
			return
		}
		res, err := h.User.GetUSerByEmail(c, &pb.GetUSerByEmailReq{Email: req.Email})
		if err != nil {
			h.Log.Error(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(req.NewPassword)
		_, err = h.User.UpdatePassword(c, &pb.UpdatePasswordReq{Id: res.Id, Password: req.NewPassword})
		if err != nil {
			h.Log.Error(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{"message": "succesfully changed"})
	h.Log.Info("reset password is succesfully ended")
}

// Refresh godoc
// @Summary Refresh token
// @Description it changes your access token
// @Tags auth
// @Param userinfo body user.LoginRes true "all"
// @Success 200 {object} string
// @Failure 400 {object} string "Invalid date"
// @Failure 401 {object} string "Invalid token"
// @Failure 500 {object} string "error while reading from server"
// @Router /auth/refresh [post]
func (h Handler) Refresh(c *gin.Context) {
	h.Log.Info("Refresh is working")
	req := pb.LoginRes{}
	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	_, err := auth.ValidateRefreshToken(req.Refreshtoken)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = auth.GetUserIdFromRefreshToken(&req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	err = auth.GeneratedAccessJWTToken(&req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	c.JSON(http.StatusOK, gin.H{
		"accesToken":   req.Accestoken,
		"refreshToken": req.Refreshtoken,
	})
}
