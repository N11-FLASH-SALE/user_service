package service

import (
	pb "auth/genproto/user"
	"auth/storage"
	"auth/storage/postgres"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type NotificationsService struct {
	pb.UnimplementedNotificationsServer
	Storage storage.IStorage
	Logger  *slog.Logger
}

func NewNotificationsService(db *sql.DB, Logger *slog.Logger) *NotificationsService {
	return &NotificationsService{
		Storage: postgres.NewPostgresStorage(db),
		Logger:  Logger,
	}
}

func (s *NotificationsService) CreateNotifications(ctx context.Context, req *pb.CreateNotificationsReq) (*pb.CreateNotificationsRes, error) {
	s.Logger.Info("CreateNotifications rpc method is working")
	resp, err := s.Storage.Notifications().CreateNotifications(ctx, req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error creating notification: %v", err))
		return nil, err
	}
	s.Logger.Info("CreateNotifications rpc method finished")
	return resp, nil
}

func (s *NotificationsService) GetAllNotifications(ctx context.Context, req *pb.GetNotificationsReq) (*pb.GetNotificationsResponse, error) {
	s.Logger.Info("GetAllNotifications rpc method is working")
	resp, err := s.Storage.Notifications().GetAllNotifications(ctx, req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error getting notifications: %v", err))
		return nil, err
	}
	s.Logger.Info("GetAllNotifications rpc method finished")
	return resp, nil
}

func (s *NotificationsService) GetAndMarkNotificationAsRead(ctx context.Context, req *pb.GetAndMarkNotificationAsReadReq) (*pb.GetAndMarkNotificationAsReadRes, error) {
	s.Logger.Info("GetAndMarkNotificationAsRead rpc method is working")
	resp, err := s.Storage.Notifications().GetAndMarkNotificationAsRead(ctx, req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error getting and marking notification as read: %v", err))
		return nil, err
	}
	s.Logger.Info("GetAndMarkNotificationAsRead rpc method finished")
	return resp, nil
}
