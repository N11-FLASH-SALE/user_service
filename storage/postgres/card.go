package postgres

import (
	pb "auth/genproto/user"
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type CardRepository struct {
	Db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{Db: db}
}

func (r *CardRepository) CreateCard(ctx context.Context, req *pb.CreateCardReq) (*pb.CreateCardRes, error) {
	// Hash the security code
	hashedSecurityCode, err := bcrypt.GenerateFromPassword([]byte(req.SecurityCode), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash security code: %w", err)
	}

	query := `
        INSERT INTO cards (user_id, card_number, expiration_date, security_code_hash)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var cardID string
	err = r.Db.QueryRowContext(ctx, query, req.UserId, req.CardNumber, req.ExpirationDate, hashedSecurityCode).Scan(&cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	return &pb.CreateCardRes{Id: cardID}, nil
}

func (r *CardRepository) GetCardsOfUser(ctx context.Context, req *pb.GetCardsOfUserReq) (*pb.GetCardsOfUserRes, error) {
	query := `
        SELECT id, user_id, card_number, expiration_date
        FROM cards
        WHERE user_id = $1
    `
	rows, err := r.Db.QueryContext(ctx, query, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}
	defer rows.Close()

	var cards []*pb.Card
	for rows.Next() {
		var card pb.Card
		if err := rows.Scan(&card.Id, &card.UserId, &card.CardNumber, &card.ExpirationDate); err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, &card)
	}

	return &pb.GetCardsOfUserRes{Cards: cards}, nil
}

func (r *CardRepository) GetCardAmount(ctx context.Context, req *pb.GetCardAmountReq) (*pb.GetCardAmountRes, error) {
	query := `
        SELECT amount
        FROM cards
        WHERE card_number = $1
    `
	var amount float64
	err := r.Db.QueryRowContext(ctx, query, req.CardNumber).Scan(&amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("card not found")
		}
		return nil, fmt.Errorf("failed to get card amount: %w", err)
	}

	return &pb.GetCardAmountRes{Amount: amount}, nil
}

func (r *CardRepository) UpdateCardAmount(ctx context.Context, req *pb.UpdateCardAmountReq) (*pb.UpdateCardAmountRes, error) {
	query := `
        UPDATE cards
        SET amount = $1, updated_at = current_timestamp
        WHERE card_number = $2
    `
	result, err := r.Db.ExecContext(ctx, query, req.Amount, req.CardNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to update card amount: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("card not found")
	}

	return &pb.UpdateCardAmountRes{Void: ""}, nil
}
