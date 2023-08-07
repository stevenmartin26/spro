package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	postgreSQLConflictErrCode pq.ErrorCode = "23505"
)

type UserRepositoryImplOptions struct {
	DB *sql.DB
}

type UserRepositoryImpl struct {
	opts *UserRepositoryImplOptions
}

func NewUserRepository(opts UserRepositoryImplOptions) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		opts: &opts,
	}
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user model.User) (uuid.UUID, *common.CustomError) {
	query := `INSERT INTO users (id, phone_number, full_name, password_hash) VALUES ($1, $2, $3, $4);`

	user.ID = uuid.New()

	if _, err := r.opts.DB.ExecContext(ctx, query, user.ID.String(), user.PhoneNumber, user.FullName, user.PasswordHash); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == postgreSQLConflictErrCode {
				return uuid.Nil, common.NewCustomError(common.ErrEntityAlreadyExists, "phone number is already used")
			}
		}
		return uuid.Nil, common.NewCustomError(common.ErrUnexpectedError, err.Error())
	}
	return user.ID, nil
}

func (r *UserRepositoryImpl) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*model.User, *common.CustomError) {
	query := `SELECT id, full_name, password_hash FROM users WHERE phone_number = $1;`

	user := model.User{
		PhoneNumber: phoneNumber,
	}

	if err := r.opts.DB.QueryRowContext(ctx, query, user.PhoneNumber).Scan(&user.ID, &user.FullName, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NewCustomError(common.ErrEntityNotFound, "user does not exist in database")
		}
		return nil, common.NewCustomError(common.ErrUnexpectedError, err.Error())
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.User, *common.CustomError) {
	query := `SELECT phone_number, full_name, password_hash FROM users WHERE id = $1;`

	user := model.User{
		ID: userID,
	}

	if err := r.opts.DB.QueryRowContext(ctx, query, user.ID.String()).Scan(&user.PhoneNumber, &user.FullName, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NewCustomError(common.ErrEntityNotFound, "user does not exist in database")
		}
		return nil, common.NewCustomError(common.ErrUnexpectedError, err.Error())
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user model.User) *common.CustomError {
	query, args := r.constructUpdateQueryAndArgs(user)

	if _, err := r.opts.DB.ExecContext(ctx, query, args...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == postgreSQLConflictErrCode {
				return common.NewCustomError(common.ErrEntityAlreadyExists, "phone number is already used")
			}
		}
		return common.NewCustomError(common.ErrUnexpectedError, err.Error())
	}
	return nil
}

func (r *UserRepositoryImpl) constructUpdateQueryAndArgs(user model.User) (string, []interface{}) {
	query := `UPDATE users SET %s WHERE id = $1;`

	setQuery := ""
	args := []interface{}{user.ID.String()}

	if user.FullName != "" {
		setQuery += fmt.Sprintf("full_name = $%d", len(args)+1)
		args = append(args, user.FullName)
	}

	if user.PhoneNumber != "" {
		if len(setQuery) > 0 {
			setQuery += ", "
		}
		setQuery += fmt.Sprintf("phone_number = $%d", len(args)+1)
		args = append(args, user.PhoneNumber)
	}

	if user.PasswordHash != "" {
		if len(setQuery) > 0 {
			setQuery += ", "
		}
		setQuery += fmt.Sprintf("password_hash = $%d", len(args)+1)
		args = append(args, user.PasswordHash)
	}

	return fmt.Sprintf(query, setQuery), args
}
