package usecase

import (
	"context"
	"time"

	"github.com/aclgo/grpc-admin/internal/admin"
	"github.com/aclgo/grpc-admin/internal/models"
	"github.com/aclgo/grpc-admin/pkg/logger"
	"github.com/google/uuid"
)

type AdminService struct {
	adminRepo admin.AdminRepo
	logger    logger.Logger
	redisRepo admin.RedisRepo
}

func NewAdminService(adminRepo admin.AdminRepo, redisRepo admin.RedisRepo, logger logger.Logger) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
		redisRepo: redisRepo,
		logger:    logger,
	}
}

func (a *AdminService) Create(ctx context.Context, params *admin.ParamsCreateAdmin) (*models.ParamsUser, error) {

	created, err := a.adminRepo.Create(ctx, &models.ParamsCreateAdmin{
		Id:        uuid.NewString(),
		Name:      params.Name,
		Lastname:  params.Lastname,
		Password:  params.HashPass(),
		Email:     params.Email,
		Role:      params.Role,
		Verified:  params.Verified,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (a *AdminService) SearchUsers(ctx context.Context, params *admin.ParamsSearchUsers) (*models.DataSearchedUser, error) {
	searched, err := a.adminRepo.Search(ctx,
		&admin.ParamsSearchUsers{
			Query: params.Query,
			Role:  params.Role,
			Page:  params.Page,
			Pagination: admin.Pagination{
				OffSet: params.Pagination.OffSet,
				Limit:  params.Pagination.Limit,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return searched, nil
}

func (a *AdminService) Delete(ctx context.Context, params *admin.ParamsDeleteUser) error {
	if err := a.adminRepo.Delete(ctx, &models.ParamsDeleteUser{
		UserId: params.UserId,
	}); err != nil {
		logger.Logger.Errorf(a.logger, "a.adminRepo.Delete:%w", err)
		return err
	}

	a.redisRepo.Publish(ctx, "disconnect_channel", admin.FormatTokenDisconnectChannel(params.UserId))

	pipe := a.redisRepo.Pipeline()
	pipe.Del(ctx, admin.FormatActiveSessionAccess(params.UserId))
	pipe.Del(ctx, admin.FormatActiveSessionRefresh(params.UserId))

	if _, err := pipe.Exec(ctx); err != nil {
		logger.Logger.Errorf(a.logger, "pipe.Exec:%w", err)
		return err
	}

	return nil
}
