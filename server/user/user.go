package user

import (
	"context"
	"fmt"

	"github.com/patrick/test-grpc-docker-gitactions/lib/env"
	"github.com/patrick/test-grpc-docker-gitactions/proto/userpb"
	svc "github.com/patrick/test-grpc-docker-gitactions/service/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct{}

func (*UserServer) UserTestCall(ctx context.Context, in *userpb.User) (*userpb.UserResponse, error) {

	fmt.Println(env.DBConnectionString)
	res, err := svc.New().WithTx(&ctx).TestCall(in)
	fmt.Println(err)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return res, nil
}
