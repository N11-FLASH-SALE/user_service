package service

import (
	pb "auth/genproto/user"
	"auth/storage"
	"auth/storage/postgres"
	"context"
	"database/sql"
	"log/slog"
)

type UserService struct {
	pb.UnimplementedUsersServer
	Repo storage.IStorage
	Log  *slog.Logger
}

func NewUserService(db *sql.DB, log *slog.Logger) *UserService {
	return &UserService{
		Repo: postgres.NewPostgresStorage(db, log),
		Log:  log,
	}
}

func (u *UserService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	u.Log.Info("Register rpc method started")
	res, err := u.Repo.User().CreateUser(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Register rpc method finished")
	return res, nil
}

func (u *UserService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	u.Log.Info("Login rpc method started")
	res, err := u.Repo.User().Login(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Login rpc method finished")
	return res, nil
}

func (u *UserService) GetUSerByEmail(ctx context.Context, req *pb.GetUSerByEmailReq) (*pb.FilterUsers, error) {
	u.Log.Info("Get user by email rpc method started")
	res, err := u.Repo.User().GetUserByEmail(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Get user by email rpc method finished")
	return res, nil
}

func (u *UserService) GetUsers(ctx context.Context, req *pb.UsersListRequest) (*pb.UsersResponse, error) {
	u.Log.Info("Get all users rpc method started")
	res, err := u.Repo.User().GetUsers(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Get all users rpc method finished")
	return res, nil
}

func (u *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.Void, error) {
	u.Log.Info("Update user rpc method started")
	err := u.Repo.User().UpdateUser(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Update user rpc method finished")
	return &pb.Void{}, nil
}

func (u *UserService) DeleteUser(ctx context.Context, req *pb.UserId) (*pb.Void, error) {
	u.Log.Info("Delete user rpc method started")
	err := u.Repo.User().DeleteUser(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Delete user rpc method finished")
	return &pb.Void{}, nil
}

func (u *UserService) GetUserById(ctx context.Context, req *pb.UserId) (*pb.GetUserResponse, error) {
	u.Log.Info("Get user by id rpc method started")
	res, err := u.Repo.User().GetUserById(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Get user by id rpc method finished")
	return res, nil
}

func (u *UserService) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordReq) (*pb.Void, error) {
	u.Log.Info("Change password rpc method started")
	err := u.Repo.User().UpdatePassword(ctx, req)
	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}
	u.Log.Info("Change password rpc method finished")
	return &pb.Void{}, nil
}
