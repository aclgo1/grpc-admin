package service

import (
	"context"
	"fmt"

	"github.com/aclgo/grpc-admin/internal/admin"
	"github.com/aclgo/grpc-admin/internal/models"
	proto "github.com/aclgo/grpc-admin/proto/admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AdminService) Register(ctx context.Context, req *proto.ParamsCreateAdmin) (*proto.ParamsUser, error) {
	// userRegisterSuccess, _ := s.observer.Meter.Float64Counter("user-register-success", metric.WithUnit("0"))
	// userRegisterError, _ := s.observer.Meter.Float64Counter("user-register-error", metric.WithUnit("0"))

	// ctx, span := s.observer.Trace.Start(ctx, "user-register")
	// defer span.End()

	// span.AddEvent("user-register")
	result, err := s.adminUC.Create(ctx, &admin.ParamsCreateAdmin{
		Name:     req.Name,
		Lastname: req.Lastname,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
		Verified: req.Verified,
	})

	if err != nil {
		// userRegisterError.Add(ctx, 1)
		// span.SetStatus(otelCodes.Error, err.Error())
		// span.End()
		return nil, err
	}

	// userRegisterSuccess.Add(ctx, 1)
	// span.SetStatus(otelCodes.Ok, "new user registred")

	return parseModelProto([]*models.ParamsUser{result})[0], nil
}

func (s *AdminService) Search(ctx context.Context, req *proto.ParamsSearchRequest) (*proto.ParamsSearchResponse, error) {

	result, err := s.adminUC.SearchUsers(ctx,
		&admin.ParamsSearchUsers{
			Query: req.Query,
			Role:  req.Role,
			Page:  int(req.Page),
			Pagination: admin.Pagination{
				OffSet: int(req.Offset),
				Limit:  int(req.Limit),
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return &proto.ParamsSearchResponse{
		Total: int64(result.Total),
		Users: parseModelProto(result.Users),
	}, nil
}

func parseModelProto(items []*models.ParamsUser) []*proto.ParamsUser {
	var users []*proto.ParamsUser

	for _, user := range items {
		user.ClearPass()

		users = append(users, &proto.ParamsUser{
			UserId:    user.UserID,
			Name:      user.Name,
			Lastname:  user.Lastname,
			Password:  user.Password,
			Email:     user.Email,
			Role:      user.Role,
			Verified:  user.Verified,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return users
}

func (s *AdminService) DeleteUser(ctx context.Context, req *proto.ParamsDeleteUserRequest) (*proto.ParamsDeleteUserResponse, error) {
	i := admin.ParamsDeleteUser{
		UserId: req.UserId,
	}

	if err := s.adminUC.Delete(ctx, &i); err != nil {
		return nil, err
	}

	out := proto.ParamsDeleteUserResponse{
		Msg: fmt.Sprintf("user id %s deleted", i.UserId),
	}

	return &out, nil
}
