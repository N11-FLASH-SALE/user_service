package postgres

import (
	pb "auth/genproto/user"
	"auth/pkg/logger"
	"auth/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	Db  *sql.DB
	Log *slog.Logger
}

func NewUserRepository(db *sql.DB) storage.IUserStorage {
	return &UserRepository{Db: db, Log: logger.NewLogger()}
}

func (u *UserRepository) CreateUser(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var userID string
	userQuery := `INSERT INTO users (fullname, username, email, password_hash, phone, role) 
                  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = tx.QueryRowContext(ctx, userQuery, req.Fullname, req.Username, req.Email, string(hashedPassword), req.Phone, req.Role).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	locationQuery := `INSERT INTO user_locations (user_id, address, city, country) 
                      VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, locationQuery, userID, req.Address, req.City, req.Country)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert user location: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &pb.RegisterRes{
		Id: userID,
	}, nil
}

func (u *UserRepository) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	query := `SELECT id, username, password_hash, role FROM users WHERE email = $1`

	var id, username, passwordHash, role string

	err := u.Db.QueryRowContext(ctx, query, req.Email).Scan(&id, &username, &passwordHash, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			query = `SELECT id, username, password_hash, role FROM users WHERE username = $1`
			err = u.Db.QueryRowContext(ctx, query, req.Email).Scan(&id, &username, &passwordHash, &role)
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
			return nil, errors.New("password is incorrect")
		}
		return nil, err
	}

	return &pb.LoginRes{
		Id:       id,
		Username: username,
		Role:     role,
	}, nil

}

func (u *UserRepository) GetUserByEmail(ctx context.Context, req *pb.GetUSerByEmailReq) (*pb.FilterUsers, error) {
	query := `
        SELECT 
            u.id, u.fullname, u.username, u.email, u.phone, u.image, u.role, 
            ul.city, ul.country, ul.address, u.created_at, u.updated_at
        FROM users u
        LEFT JOIN user_locations ul ON u.id = ul.user_id
        WHERE u.email = $1`

	var user pb.FilterUsers
	var image sql.NullString

	err := u.Db.QueryRowContext(ctx, query, req.Email).Scan(
		&user.Id, &user.Fullname, &user.Username, &user.Email, &user.PhoneNumber,
		&image, &user.Role, &user.City, &user.Country, &user.Address,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if !image.Valid {
		user.ImageUrl = ""
	} else {
		user.ImageUrl = image.String
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
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

func (u *UserRepository) GetUserById(ctx context.Context, req *pb.UserId) (*pb.GetUserResponse, error) {
	query := `SELECT id, fullname, username, email, phone, image, role, city, country, address, created_at, updated_at 
              FROM users 
              LEFT JOIN user_locations ON users.id = user_locations.user_id 
              WHERE users.id = $1`

	var user pb.FilterUsers
	var image sql.NullString

	err := u.Db.QueryRowContext(ctx, query, req.Id).Scan(
		&user.Id,
		&user.Fullname,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&image,
		&user.Role,
		&user.City,
		&user.Country,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if !image.Valid {
		user.ImageUrl = ""
	} else {
		user.ImageUrl = image.String
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Prepare the response
	response := &pb.GetUserResponse{
		User: &user,
	}

	return response, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) error {
	query := `update users set `
	n := 1
	var arr []interface{}
	if len(req.Fullname) > 0 {
		query += fmt.Sprintf("fullname=$%d, ", n)
		arr = append(arr, req.Fullname)
		n++
	}
	if len(req.Username) > 0 {
		query += fmt.Sprintf("username=$%d, ", n)
		arr = append(arr, req.Username)
		n++
	}
	if len(req.PhoneNumber) > 0 {
		query += fmt.Sprintf("phone=$%d, ", n)
		arr = append(arr, req.PhoneNumber)
		n++
	}
	if len(req.ImageUrl) > 0 {
		query += fmt.Sprintf("image=$%d, ", n)
		arr = append(arr, req.ImageUrl)
		n++
	}
	arr = append(arr, req.Id)
	query += fmt.Sprintf("updated_at=current_timestamp where id=$%d and deleted_at=0", n)
	_, err := u.Db.Exec(query, arr...)
	if err != nil {
		return err
	}

	query = `update user_locations set `
	n = 1
	var arr1 []interface{}
	if len(req.City) > 0 {
		query += fmt.Sprintf("city=$%d, ", n)
		arr1 = append(arr1, req.City)
		n++
	}
	if len(req.Country) > 0 {
		query += fmt.Sprintf("country=$%d, ", n)
		arr1 = append(arr1, req.Country)
		n++
	}
	if len(req.Address) > 0 {
		query += fmt.Sprintf("address=$%d, ", n)
		arr1 = append(arr1, req.Address)
		n++
	}
	if len(req.PostalCode) > 0 {
		query += fmt.Sprintf("postal_code=$%d, ", n)
		arr1 = append(arr1, req.PostalCode)
		n++
	}

	if len(arr1) == 0 {
		return nil
	} else {
		arr1 = append(arr1, req.Id)
		query = query[:len(query)-2]
		query += fmt.Sprintf(" where user_id=$%d", n)
		_, err = u.Db.Exec(query, arr1...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *UserRepository) GetUsers(ctx context.Context, req *pb.UsersListRequest) (*pb.UsersResponse, error) {
	query := `SELECT id, fullname, username, email, phone, image, role, city, country, address, created_at, updated_at 
	FROM users 
	LEFT JOIN user_locations ON users.id = user_locations.user_id
	where `
	n := 1
	var arr []interface{}
	if len(req.Users.Id) > 0 {
		query += fmt.Sprintf("id=$%d and ", n)
		arr = append(arr, req.Users.Id)
		n++
	}
	if len(req.Users.Fullname) > 0 {
		query += fmt.Sprintf("fullname=$%d and ", n)
		arr = append(arr, req.Users.Fullname)
		n++
	}
	if len(req.Users.Username) > 0 {
		query += fmt.Sprintf("username=$%d and ", n)
		arr = append(arr, req.Users.Username)
		n++
	}
	if len(req.Users.Email) > 0 {
		query += fmt.Sprintf("email=$%d and ", n)
		arr = append(arr, req.Users.Email)
		n++
	}
	if len(req.Users.PhoneNumber) > 0 {
		query += fmt.Sprintf("phone=$%d and ", n)
		arr = append(arr, req.Users.PhoneNumber)
		n++
	}
	if len(req.Users.Role) > 0 {
		query += fmt.Sprintf("role=$%d and ", n)
		arr = append(arr, req.Users.Role)
		n++
	}
	if len(req.Users.City) > 0 {
		query += fmt.Sprintf("city=$%d and ", n)
		arr = append(arr, req.Users.City)
		n++
	}
	if len(req.Users.Country) > 0 {
		query += fmt.Sprintf("country=$%d and ", n)
		arr = append(arr, req.Users.Country)
		n++
	}
	if len(req.Users.Address) > 0 {
		query += fmt.Sprintf("address=$%d and ", n)
		arr = append(arr, req.Users.Address)
		n++
	}
	query += "deleted_at=0 "
	query += fmt.Sprintf("Limit %d ", req.Limit)

	if req.Offset > 0 {
		query += fmt.Sprintf("OFFSET %d ", req.Offset)
	}

	rows, err := u.Db.Query(query, arr...)
	fmt.Print("\n\n", query, "\n")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := pb.UsersResponse{}
	for rows.Next() {
		us := pb.FilterUsers{}
		var image sql.NullString
		err = rows.Scan(&us.Id, &us.Fullname, &us.Username, &us.Email, &us.PhoneNumber, &image, &us.Role, &us.City, &us.Country, &us.Address, &us.CreatedAt, &us.UpdatedAt)
		if !image.Valid {
			us.ImageUrl = ""
		} else {
			us.ImageUrl = image.String
		}
		if err != nil {
			return nil, err
		}
		res.Users = append(res.Users, &us)
	}
	res.Offset = req.Offset
	res.Limit = req.Limit
	return &res, nil
}

func (u *UserRepository) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordReq) error {
	fmt.Println(req)
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
