package e2e

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aclgo/grpc-admin/config"
	"github.com/aclgo/grpc-admin/internal/admin"
	proto "github.com/aclgo/grpc-admin/proto/admin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	l, err := grpc.DialContext(ctx, fmt.Sprintf(":%v", cfg.ApiPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.DialContext: %v", err)
	}

	svcClient := proto.NewAdminServiceClient(l)

	e2e := Newe2eConfig(svcClient)

	// ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	if err := e2e.registerUser(ctx); err != nil {
		log.Fatalf("registerUser: %v", err)
	}

	if err := e2e.searchUser(ctx); err != nil {
		log.Fatalf("searchUser: %v", err)
	}

}

type e2eConfig struct {
	svcClient proto.AdminServiceClient
}

func Newe2eConfig(adminClient proto.AdminServiceClient) *e2eConfig {
	return &e2eConfig{
		svcClient: adminClient,
	}
}

func (e2e *e2eConfig) registerUser(ctx context.Context) error {

	createParams := admin.ParamsCreateAdmin{
		Name:     "aclgo",
		Lastname: "e2e",
		Password: "aclgo_e2e",
		Email:    uuid.NewString() + "@gmail,com",
	}

	_, err := e2e.svcClient.Register(
		ctx, &proto.ParamsCreateAdmin{
			Name:     createParams.Name,
			Lastname: createParams.Lastname,
			Password: createParams.Password,
			Email:    createParams.Email,
		},
	)

	if err != nil {
		return errors.Wrap(err, "registerUser.Register: %v")
	}

	log.Printf("test e2e register new admin or user PASS")

	return nil
}

func (e2e *e2eConfig) searchUser(ctx context.Context) error {

	paramsSearch := admin.ParamsSearchUsers{
		Query: "",
		Role:  "",
		Page:  0,
		Pagination: admin.Pagination{
			OffSet: 0,
			Limit:  0,
		},
	}

	_, err := e2e.svcClient.Search(
		ctx,
		&proto.ParamsSearchRequest{
			Query:  paramsSearch.Query,
			Role:   paramsSearch.Role,
			Page:   int32(paramsSearch.Page),
			Offset: int32(paramsSearch.Pagination.OffSet),
			Limit:  int32(paramsSearch.Pagination.Limit),
		},
	)

	if err != nil {
		return fmt.Errorf("searchUser.Search: %v", err)
	}

	log.Printf("test e2e search admins or users PASS")

	return nil
}
