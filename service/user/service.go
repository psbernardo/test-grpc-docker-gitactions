package user

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/patrick/test-grpc-docker-gitactions/lib/dbop"
	"github.com/patrick/test-grpc-docker-gitactions/proto/userpb"
)

//Service ORM
type Service struct {
	tx *gorm.DB
}

func New() *Service {
	return &Service{}
}

func (s *Service) WithTx(ctx *context.Context) *Service {
	s.tx = dbop.MSDB()
	return s
}

func (s *Service) TestCall(in *userpb.User) (*userpb.UserResponse, error) {

	if err := s.tx.Exec(`INSERT INTO [api].[status]([account_id]) VALUES ('dfdf')`).Error; err != nil {
		return nil, err
	}

	fmt.Println("call server")
	return &userpb.UserResponse{Status: "GCP call sucess" + in.Name}, nil
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
