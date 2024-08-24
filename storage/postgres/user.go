package postgres

import (
	pb "auth/genproto/user"
	"auth/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) storage.IUserStorage {
	return &UserRepository{Db: db}
}

func (u UserRepository) CreateUser(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var userID string
	userQuery := `INSERT INTO users (email, password_hash, full_name, phone_number, address, role) 
                  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = tx.QueryRowContext(ctx, userQuery, req.Email, string(hashedPassword), req.Fullname, req.Phone, req.Address, req.Role).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &pb.RegisterRes{
		Id: userID,
	}, nil
}

func (u UserRepository) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	query := `SELECT id, password_hash , role FROM users WHERE email = $1 and deleted_at=0`

	var id, passwordHash, role string
	err := u.Db.QueryRowContext(ctx, query, req.Email).Scan(&id, &passwordHash, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			query = `SELECT id, password_hash, role FROM users WHERE phone_number = $1`
			err = u.Db.QueryRowContext(ctx, query, req.Email).Scan(&id, &passwordHash, &role)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, errors.New("user not found")
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("password is incorrect")
		}
		return nil, err
	}

	return &pb.LoginRes{
		Id:   id,
		Role: role,
	}, nil
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, req *pb.GetUSerByEmailReq) (*pb.GetUserResponse, error) {
	query := `SELECT id, full_name,  phone_number, address, photo ,role, created_at FROM users WHERE email = $1 AND deleted_at=0`

	user := pb.GetUserResponse{
		Email: req.Email,
	}

	var photo sql.NullString

	err := u.Db.QueryRowContext(ctx, query, req.Email).Scan(&user.Id, &user.Fullname, &user.Phone, &user.Address, &photo, &user.Role, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	if photo.Valid {
		user.Photo = photo.String
	} else {
		user.Photo = ""
	}

	return &user, nil
}

func (u *UserRepository) GetUserById(ctx context.Context, req *pb.UserId) (*pb.GetUserResponse, error) {
	query := `SELECT email, full_name,  phone_number, address, photo ,role, created_at FROM users WHERE id = $1 AND deleted_at=0`

	user := pb.GetUserResponse{
		Id: req.Id,
	}

	var photo sql.NullString

	err := u.Db.QueryRowContext(ctx, query, req.Id).Scan(&user.Email, &user.Fullname, &user.Phone, &user.Address, &photo, &user.Role, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	if photo.Valid {
		user.Photo = photo.String
	} else {
		user.Photo = ""
	}

	return &user, nil
}

func (u *UserRepository) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordReq) error {
	query := `update users set password_hash=$1 where id=$2 and deleted_at=0`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	result, err := u.Db.ExecContext(ctx, query, hashedPassword, req.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) error {
	query := `update users set `
	n := 1
	var arr []interface{}
	if len(req.Fullname) > 0 {
		query += fmt.Sprintf("full_name=$%d, ", n)
		arr = append(arr, req.Fullname)
		n++
	}
	if len(req.Address) > 0 {
		query += fmt.Sprintf("address=$%d, ", n)
		arr = append(arr, req.Address)
		n++
	}
	if len(req.Photo) > 0 {
		query += fmt.Sprintf("photo=$%d, ", n)
		arr = append(arr, req.Photo)
		n++
	}
	if len(req.Phone) > 0 {
		query += fmt.Sprintf("phone_number=$%d, ", n)
		arr = append(arr, req.Phone)
		n++
	}
	arr = append(arr, req.Id)
	query += fmt.Sprintf("updated_at=current_timestamp where id=$%d and deleted_at=0", n)
	result, err := u.Db.Exec(query, arr...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, req *pb.UserId) error {
	query := `UPDATE users SET deleted_at = date_part('epoch', current_timestamp)::INT 
	WHERE id = $1 and deleted_at=0`

	result, err := u.Db.ExecContext(ctx, query, req.Id)
	if err != nil {
		return fmt.Errorf("failed to update deleted_at: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (u *UserRepository) ResetPassword(ctx context.Context, req *pb.ResetPasswordReq) error {
	query := `SELECT password_hash FROM users WHERE id = $1 AND deleted_at=0`
	var passwordHash string
	err := u.Db.QueryRowContext(ctx, query, req.Id).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Oldpassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return fmt.Errorf("password is incorrect")
		}
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Newpassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	query = `UPDATE users SET password_hash=$1 WHERE id=$2 AND deleted_at=0`
	result, err := u.Db.ExecContext(ctx, query, hashedPassword, req.Id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (u *UserRepository) IsUserExist(ctx context.Context, req *pb.UserId) error {
	var exists bool
	err := u.Db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = $1
		)
	`, req.GetId()).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("user with id %s does not exist", req.GetId())
	}

	return nil
}
