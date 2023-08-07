package repository

import (
	"context"
	"database/sql"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/google/uuid"
)

type LoginLogRepositoryImplOptions struct {
	DB *sql.DB
}

type LoginLogRepositoryImpl struct {
	opts *LoginLogRepositoryImplOptions
}

func NewLoginLogRepositoryImpl(opts LoginLogRepositoryImplOptions) *LoginLogRepositoryImpl {
	return &LoginLogRepositoryImpl{
		opts: &opts,
	}
}

func (r *LoginLogRepositoryImpl) Save(ctx context.Context, log model.LoginLog) (uuid.UUID, *common.CustomError) {
	query := `INSERT INTO login_logs (id, user_id, login_at) VALUES ($1, $2, $3);`

	log.ID = uuid.New()

	if _, err := r.opts.DB.ExecContext(ctx, query, log.ID.String(), log.UserID.String(), log.LoginAt); err != nil {
		return uuid.Nil, common.NewCustomError(common.ErrUnexpectedError, err.Error())
	}
	return log.ID, nil
}
