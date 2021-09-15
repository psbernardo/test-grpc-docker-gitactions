package user

import (
	"github.com/patrick/test-grpc-docker-gitactions/proto/userpb"
)

//Administrator is model for database column mapping
type User struct {
	Id       uint32 `gorm:"column:id;primary_key"`
	Name     string `gorm:"column:name"`
	LastName string
}

//TableName returns administrator table name
func (u User) TableName() string {
	return "dbo.user1"
}

//NewAdministrator creates new administrator model
func NewUser() *User {
	return &User{}
}

//ToProto convert administrator model to proto
func (a *User) ToProto() *userpb.User {

	return nil
}
