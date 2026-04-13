package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolas-andreoli/trainer-app/internal/envutil"
	"github.com/nicolas-andreoli/trainer-app/internal/garmin"
	"github.com/nicolas-andreoli/trainer-app/internal/tools"
)

func main() {
	envutil.LoadDotEnv(".env")

	s := server.NewMCPServer(
		"garmin-mcp",
		"1.0.0",
	)

	client := garmin.NewClient()
	tools.RegisterGarminTools(s, client)

	stdio := server.NewStdioServer(s)
	if err := stdio.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "garmin-mcp server error: %v\n", err)
		os.Exit(1)
	}
}
