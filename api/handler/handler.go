package handler

import (
	"auth/genproto/user"
	"log/slog"
)

type Handler struct {
	User user.UsersClient
	Log  *slog.Logger
}
