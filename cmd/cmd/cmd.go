package cmd

import (
	"github.com/patrick/test-grpc-docker-gitactions/cmd/db"
	"github.com/patrick/test-grpc-docker-gitactions/cmd/server"

	"github.com/spf13/cobra"
)

//Cmd implements the sas-clearing command
var Cmd = &cobra.Command{
	Use:          "grpc",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(server.Cmd)
	Cmd.AddCommand(db.Cmd)
}
