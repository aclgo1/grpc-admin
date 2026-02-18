package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/aclgo/grpc-admin/internal/admin"
	"github.com/aclgo/grpc-admin/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	sqlCreate = `INSERT INTO users (user_id, name, last_name, password, email,
		role, verified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING user_id, name, last_name, password, email, role, verified, created_at, updated_at`
	sqlDelete = `delete from users where user_id=$1`
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewpostgresRepo(db *sqlx.DB) *postgresRepo {
	return &postgresRepo{
		db: db,
	}
}

func (a *postgresRepo) Create(ctx context.Context, user *models.ParamsCreateAdmin) (*models.ParamsUser, error) {

	var created models.ParamsUser

	err := a.db.QueryRowxContext(ctx, sqlCreate,
		user.Id,
		user.Name,
		user.Lastname,
		user.Password,
		user.Email,
		user.Role,
		user.Verified,
		user.CreatedAt,
		user.UpdatedAt,
	).StructScan(&created)
	switch {
	case errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled):
		return nil, err
	case err != nil:
		return nil, errors.Wrap(err, "Create.QueryRowContext")
	default:
		return &created, nil
	}
}

func (a *postgresRepo) Search(ctx context.Context, params *admin.ParamsSearchUsers) (*models.DataSearchedUser, error) {
	var (
		args = []any{"%" + params.Query + "%"}
		w    = []string{"name LIKE $1"}
	)

	if params.Role != "" {
		args = append(args, params.Role)
		w = append(w, fmt.Sprintf(` "role" = $%v`, len(args)))
	}

	where := strings.Join(w, " AND ")

	var total int

	sqlTotal := fmt.Sprintf(`SELECT COUNT(*) AS total FROM "users" WHERE %s`, where)
	err := a.db.QueryRowContext(ctx, sqlTotal, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(`SELECT * FROM "users" WHERE %s ORDER BY "user_id" DESC`, where)
	if params.Pagination.Limit > 0 {
		args = append(args, params.Pagination.Limit)
		sql += fmt.Sprintf(` LIMIT $%d`, len(args))
	}

	if params.Pagination.OffSet > 0 {
		args = append(args, params.Pagination.OffSet)
		sql += fmt.Sprintf(` OFFSET $%d`, len(args))
	}

	rows, err := a.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []*models.ParamsUser

	for rows.Next() {
		var item models.ParamsUser

		if err := rows.StructScan(&item); err != nil {
			return nil, errors.Wrap(err, "Search.StructScan")
		}

		item.ClearPass()

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "Search.Err")
	}

	if err := rows.Close(); err != nil {
		return nil, errors.Wrap(err, "Search.Close")
	}

	return &models.DataSearchedUser{Total: total, Users: items}, nil
}

func (a *postgresRepo) Delete(ctx context.Context, params *models.ParamsDeleteUser) error {
	if _, err := a.db.Exec(sqlDelete, params.UserId); err != nil {
		return err
	}

	return nil
}
