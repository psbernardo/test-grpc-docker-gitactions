package table

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var UserTable = []*gormigrate.Migration{
	{
		ID: "202109015-000",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
		IF OBJECT_ID('dbo.user1') IS NOT NULL
        	DROP TABLE dbo.user1
    

			CREATE TABLE dbo.user1(
				id int IDENTITY(1,1) NOT NULL
				,name                                    varchar(100) DEFAULT('')
				,last_name                                varchar(100) DEFAULT('')
				,date_time   datetime DEFAULT(getdate())
			);
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`
			DROP TABLE  dbo.user; 
			`).Error
		},
	},
}
