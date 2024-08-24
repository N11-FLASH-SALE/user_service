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

type CardsService struct {
	pb.UnimplementedCardsServer
	Cards  storage.IStorage
	Logger *slog.Logger
}

func NewCardsService(db *sql.DB, Logger *slog.Logger) *CardsService {
	return &CardsService{
		Cards:  postgres.NewPostgresStorage(db),
		Logger: Logger,
	}
}

func (c *CardsService) CreateCard(ctx context.Context, req *pb.CreateCardReq) (*pb.CreateCardRes, error) {
	c.Logger.Info("CreateCard rpc is working")
	resp, err := c.Cards.Card().CreateCard(ctx, req)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("create card error: %v", err))
		return nil, err
	}
	c.Logger.Info("CreateCard rpc method finished")
	return resp,nil
}
func (c *CardsService) GetCardsOfUser(ctx context.Context, req *pb.GetCardsOfUserReq) (*pb.GetCardsOfUserRes, error) {
	c.Logger.Info("GetCardOfUser rpc is working")
	resp, err := c.Cards.Card().GetCardsOfUser(ctx, req)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("get card of user error: %v", err))
		return nil, err
	}
	c.Logger.Info("GetCardOfUser rpc method finished")
	return resp,nil
}
func (c *CardsService) GetCardAmount(ctx context.Context, req *pb.GetCardAmountReq) (*pb.GetCardAmountRes, error) {
	c.Logger.Info("GetCardAmount rpc is working")
	resp, err := c.Cards.Card().GetCardAmount(ctx, req)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("get card amount error: %v", err))
		return nil, err
	}
	c.Logger.Info("GetCardAmount rpc method finished")
	return resp,nil
}
func (c *CardsService) UpdateCardAmount(ctx context.Context, req *pb.UpdateCardAmountReq) (*pb.UpdateCardAmountRes, error) {
	c.Logger.Info("UpdateCardAmount rpc is working")
	resp, err := c.Cards.Card().UpdateCardAmount(ctx, req)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("update card amount error: %v", err))
		return nil, err
	}
	c.Logger.Info("UpdateCardAmount rpc method finished")
	return resp,nil
}
