package postgres

import (
	pb "auth/genproto/user"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCard(t *testing.T) {
	db, err := ConnectionDb()
	if err != nil {
		panic(err)
	}

	repo := NewCardRepository(db)

	req := pb.CreateCardReq{
		UserId:         "user123",
		CardNumber:     "1234567890123456",
		ExpirationDate: "12/25",
		SecurityCode:   "123",
	}
	ctx := context.Background()

	resp, err := repo.CreateCard(ctx, &req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println(resp)
	assert.NotEmpty(t, resp.Id)
}

func TestGetCardsOfUser(t *testing.T) {
	db, err := ConnectionDb()
	if err != nil {
		panic(err)
	}

	repo := NewCardRepository(db)

	req := pb.GetCardsOfUserReq{
		UserId: "user123",
	}
	ctx := context.Background()

	resp, err := repo.GetCardsOfUser(ctx, &req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println(resp)
	assert.NotEmpty(t, resp.Cards)
}

func TestGetCardAmount(t *testing.T) {
	db, err := ConnectionDb()
	if err != nil {
		panic(err)
	}

	repo := NewCardRepository(db)

	req := pb.GetCardAmountReq{
		CardId: "card123",
	}
	ctx := context.Background()

	resp, err := repo.GetCardAmount(ctx, &req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println(resp)
	assert.NotEmpty(t, resp.Amount)
}

func TestUpdateCardAmount(t *testing.T) {
	db, err := ConnectionDb()
	if err != nil {
		panic(err)
	}

	repo := NewCardRepository(db)

	req := pb.UpdateCardAmountReq{
		CardId: "card123",
		Amount: "100.00",
	}
	ctx := context.Background()

	resp, err := repo.UpdateCardAmount(ctx, &req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println(resp)
	assert.Equal(t, "", resp.Void)
}
