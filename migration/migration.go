package migration

import (
	"fmt"

	"github.com/jinzhu/gorm"
	tbl "github.com/patrick/test-grpc-docker-gitactions/migration/user/table"
	"gopkg.in/gormigrate.v1"
)

var Migrations = []*gormigrate.Migration{}

func New(db *gorm.DB) *gormigrate.Gormigrate {
	Migrations = append(Migrations, tbl.UserTable...)

	return gormigrate.New(db, &gormigrate.Options{
		UseTransaction: true,
	}, Migrations)
}

type Migration struct {
	ID string
}

//Migrate ...
func Migrate(db *gorm.DB) error {
	m := New(db)
	if err := m.Migrate(); err != nil {
		return err
	}

	mig := Migration{}
	if err := db.Model(&Migration{}).Last(&mig).Error; err != nil {
		return err
	}
	fmt.Printf("migrated to %s\n", mig.ID)

	return nil
}
