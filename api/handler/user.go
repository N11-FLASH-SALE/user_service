package handler

import (
	"auth/api/auth"
	"auth/api/email"
	pb "auth/genproto/user"
	"auth/storage/redis"
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
	if !email.IsValidEmail(req.Email) {
		h.Log.Error("Invalid email")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
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
		c.JSON(500, gin.H{"error": err.Error()})
	}
	err = auth.GeneratedRefreshJWTToken(res)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
	}

	h.Log.Info("login is succesfully ended")
	c.JSON(http.StatusOK, gin.H{
		"accesToken":   res.Accestoken,
		"refreshToken": res.Refreshtoken,
	})
}

// Refresh godoc
// @Summary Refresh token
// @Description it generates new access token
// @Tags auth
// @Param token body user.Tokens true "enough"
// @Success 200 {object} string "tokens"
// @Failure 400 {object} string "Invalid date"
// @Failure 500 {object} string "error while reading from server"
// @Router /auth/refresh [post]
func (h Handler) Refresh(c *gin.Context) {
	h.Log.Info("Refresh is working")
	tok := pb.Tokens{}
	if err := c.BindJSON(&tok); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	req := pb.LoginRes{Refreshtoken: tok.Refreshtoken}

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
	h.Log.Info("Refresh is succesfully ended")
	c.JSON(http.StatusOK, gin.H{
		"accesToken":   req.Accestoken,
		"refreshToken": req.Refreshtoken,
	})
}

// ForgotPassword godoc
// @Summary Forgot Password
// @Description it send code to your email address
// @Tags auth
// @Param token body user.GetUSerByEmailReq true "enough"
// @Success 200 {object} string "message"
// @Failure 400 {object} string "Invalid date"
// @Failure 500 {object} string "error while reading from server"
// @Router /auth/forgot-password [post]
func (h Handler) ForgotPassword(c *gin.Context) {
	h.Log.Info("ForgotPassword is working")
	var req pb.GetUSerByEmailReq
	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := email.EmailCode(req.Email)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending email"})
		return
	}
	err = redis.StoreCodes(c, res, req.Email)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing codes in Redis"})
		return
	}
	h.Log.Info("ForgotPassword succeeded")
	c.JSON(200, gin.H{"message": "Password reset code sent to your email"})

}

// ResetPassword godoc
// @Summary Reset Password
// @Description it Reset your Password
// @Tags auth
// @Param token body user.ResetPassReq true "enough"
// @Success 200 {object} string "message"
// @Failure 400 {object} string "Invalid date"
// @Failure 500 {object} string "error while reading from server"
// @Router /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	h.Log.Info("ResetPassword is working")
	var req pb.ResetPassReq
	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code, err := redis.GetCodes(c, req.Email)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	if code != req.Code {
		h.Log.Error("Invalid code")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
		return
	}
	res, err := h.User.GetUSerByEmail(c, &pb.GetUSerByEmailReq{Email: req.Email})
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = h.User.UpdatePassword(c, &pb.UpdatePasswordReq{Id: res.Id, Password: req.Password})
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password"})
		return
	}
	c.JSON(200, gin.H{"message": "Password reset successfully"})
}

// logout godoc
// @Security ApiKeyAuth
// @Summary logout user
// @Description logout
// @Tags user
// @Success 200 {object} string "message"
// @Router /user/logout [post]
func (h Handler) Logout(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Successfully logged out"})
}

// GetUserProfile godoc
// @Security ApiKeyAuth
// @Summary Get User Profile
// @Description Get User Profile by token
// @Tags user
// @Success 200 {object} user.GetUserResponse
// @Failure 400 {object} string "Invalid date"
// @Failure 500 {object} string "error while reading from server"
// @Router /user/profile [get]
func (h Handler) GetUserProfile(c *gin.Context) {
	h.Log.Info("GetUserProfile is working")
	var req pb.LoginRes
	req.Accestoken = c.GetHeader("Authorization")
	err := auth.GetUserIdFromAccesToken(&req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	res, err := h.User.GetUserById(c, &pb.UserId{Id: req.Id})
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user"})
		return
	}
	h.Log.Info("GetUserProfile successful finished")
	c.JSON(200, res)
}

// UpdateUserProfile godoc
// @Security ApiKeyAuth
// @Summary Update User Profile
// @Description Update User Profile by token
// @Tags user
// @Param userinfo body user.UpdateUserRequest true "all"
// @Success 200 {object} string "message"
// @Failure 400 {object} string "Invalid date"
// @Failure 500 {object} string "error while reading from server"
// @Router /user/profile [put]
func (h Handler) UpdateUserProfile(c *gin.Context) {
	h.Log.Info("UpdateUserProfile is working")
	var req pb.LoginRes
	req.Accestoken = c.GetHeader("Authorization")
	err := auth.GetUserIdFromAccesToken(&req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var user pb.UpdateUserRequest
	if err := c.BindJSON(&user); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Id = req.Id
	_, err = h.User.UpdateUser(c, &user)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}
	h.Log.Info("User updated successfully finished")
	c.JSON(200, gin.H{"message": "User updated successfully"})
}

// ChangePassword godoc
// @Security ApiKeyAuth
// @Summary Update User Profile
// @Description Update User Profile by token
// @Tags user
// @Param userinfo body user.ResetPasswordReq true "all"
// @Success 200 {object} string "message"
// @Failure 400 {object} string "Invalid date"
// @Failure 500 {object} string "error while reading from server"
// @Router /user/change-password [post]
func (h Handler) ChangePassword(c *gin.Context) {
	h.Log.Info("ChangePassword is working")
	var req pb.LoginRes
	req.Accestoken = c.GetHeader("Authorization")
	err := auth.GetUserIdFromAccesToken(&req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var user pb.ResetPasswordReq
	if err := c.BindJSON(&user); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Id = req.Id
	_, err = h.User.ResetPassword(c, &user)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error resetting password"})
		return
	}
	h.Log.Info("Password changed successfully finished")
	c.JSON(200, gin.H{"message": "Password changed successfully"})
}
