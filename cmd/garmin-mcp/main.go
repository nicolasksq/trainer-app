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
	if os.Getenv("GARMIN_EMAIL") == "" || os.Getenv("GARMIN_PASSWORD") == "" {
		fmt.Fprintf(os.Stderr, "Warning: GARMIN_EMAIL or GARMIN_PASSWORD not set. Run 'go run ./cmd/setup' to configure credentials.\n")
	}
	tools.RegisterGarminTools(s, client)

	stdio := server.NewStdioServer(s)
	if err := stdio.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "garmin-mcp server error: %v\n", err)
		os.Exit(1)
	}
}
