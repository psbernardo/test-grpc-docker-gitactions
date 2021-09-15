package db

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/patrick/test-grpc-docker-gitactions/lib/dbop"
	"github.com/patrick/test-grpc-docker-gitactions/lib/env"
	"github.com/patrick/test-grpc-docker-gitactions/migration"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var Cmd = &cobra.Command{
	Use: "db",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var migrateCmd = &cobra.Command{
	Use: "migrate",
	RunE: func(cmd *cobra.Command, args []string) error {
		db := dbop.MSDB()

		if cmd.Flags().Changed("verbose") {
			db.LogMode(verbose)
		}

		return migrate(db)
	},
}

type Migration struct {
	ID string
}

func migrate(db *gorm.DB) error {
	m := migration.New(db)

	if err := m.Migrate(); err != nil {
		fmt.Println("Error Message 46", err.Error())
		return err
	}

	mig := Migration{}
	if err := db.Model(&Migration{}).Last(&mig).Error; err != nil {
		fmt.Println("Error Message 52", err.Error())
		return err
	}
	fmt.Printf("migrated to %s\n", mig.ID)

	return nil
}

//creatDbCmd create Clearing database for github validation
var createDbCmd = &cobra.Command{
	Use: "create-db",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbName := "testdb"
		connectionString := strings.Replace(env.DBConnectionString, "testdb", "master", 1)
		db, err := dbop.NewMSDB(connectionString)
		if err != nil {
			return err
		}

		db = db.Exec("DROP DATABASE IF EXISTS ?;", dbName)
		if db.Error != nil {
			fmt.Println(db.Error)
		}

		db = db.Exec("CREATE DATABASE testdb;")
		if db.Error != nil {
			fmt.Println(db.Error)
		}

		defer db.Close()

		return nil
	},
}

func init() {
	Cmd.AddCommand(migrateCmd)
	Cmd.AddCommand(createDbCmd)

}
