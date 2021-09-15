package main

import (
	"os"

	grCMD "github.com/patrick/test-grpc-docker-gitactions/cmd/cmd"
)

func main() {
	if err := grCMD.Cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
