package storage

import (
	pb "auth/genproto/user"
	"context"
)

type IStorage interface {
	User() IUserStorage
	Notifications() INotificationStorage
	Card() CardStorage
	Close()
}

type IUserStorage interface {
	CreateUser(context.Context, *pb.RegisterReq) (*pb.RegisterRes, error)
	Login(context.Context, *pb.LoginReq) (*pb.LoginRes, error)
	GetUserByEmail(context.Context, *pb.GetUSerByEmailReq) (*pb.GetUserResponse, error)
	GetUserById(context.Context, *pb.UserId) (*pb.GetUserResponse, error)
	UpdatePassword(context.Context, *pb.UpdatePasswordReq) error
	UpdateUser(context.Context, *pb.UpdateUserRequest) error
	DeleteUser(context.Context, *pb.UserId) error
	ResetPassword(context.Context, *pb.ResetPasswordReq) error
	IsUserExist(context.Context, *pb.UserId) error
}

type CardStorage interface {
	CreateCard(context.Context, *pb.CreateCardReq) (*pb.CreateCardRes, error)
	GetCardsOfUser(context.Context, *pb.GetCardsOfUserReq) (*pb.GetCardsOfUserRes, error)
	GetCardAmount(context.Context, *pb.GetCardAmountReq) (*pb.GetCardAmountRes, error)
	UpdateCardAmount(context.Context, *pb.UpdateCardAmountReq) (*pb.UpdateCardAmountRes, error)
	DeleteCard(context.Context, *pb.DeleteCardReq) error
}
type INotificationStorage interface {
	CreateNotifications(context.Context, *pb.CreateNotificationsReq) (*pb.CreateNotificationsRes, error)
	GetAllNotifications(context.Context, *pb.GetNotificationsReq) (*pb.GetNotificationsResponse, error)
	GetAndMarkNotificationAsRead(context.Context, *pb.GetAndMarkNotificationAsReadReq) (*pb.GetAndMarkNotificationAsReadRes, error)
}
