package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/patrick/test-grpc-docker-gitactions/grpcserver"

	"github.com/spf13/cobra"
)

//Cmd implements the start command
var Cmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var startCmd = &cobra.Command{
	Use: "start",
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("start")
		grpcServer, err := grpcserver.StartGrpc()
		if err != nil {
			return err
		}

		httpServer, err := grpcserver.StartGrpcWeb(grpcServer)
		if err != nil {
			grpcServer.Stop()
			return err
		}

		// wait for ctrl C to exit
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)

		// block until a signal is received
		<-ch

		fmt.Println("stopping the server")
		grpcServer.GracefulStop()
		httpServer.Shutdown(context.Background())

		fmt.Println("server stopped")

		return nil
	},
}

func init() {
	Cmd.AddCommand(startCmd)
}
