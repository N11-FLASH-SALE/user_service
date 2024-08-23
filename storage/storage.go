package storage

import (
	pb "auth/genproto/user"
	"context"
)

type IStorage interface {
	User() IUserStorage
	Close()
}

type IUserStorage interface {
	CreateUser(context.Context, *pb.RegisterReq) (*pb.RegisterRes, error)
	Login(context.Context, *pb.LoginReq) (*pb.LoginRes, error)
	GetUserByEmail(context.Context, *pb.GetUSerByEmailReq) (*pb.FilterUsers, error)
	DeleteUser(context.Context, *pb.UserId) error
	GetUserById(context.Context, *pb.UserId) (*pb.GetUserResponse, error)
	UpdateUser(context.Context, *pb.UpdateUserRequest) error
	GetUsers(context.Context, *pb.UsersListRequest) (*pb.UsersResponse, error)
	UpdatePassword(context.Context, *pb.UpdatePasswordReq) error
}
