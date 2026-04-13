package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolas-andreoli/trainer-app/internal/envutil"
	"github.com/nicolas-andreoli/trainer-app/internal/strava"
	"github.com/nicolas-andreoli/trainer-app/internal/tools"
)

func main() {
	envutil.LoadDotEnv(".env")

	if len(os.Args) > 1 && os.Args[1] == "auth" {
		runAuthFlow()
		return
	}

	client, err := strava.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Strava client: %v", err)
	}

	s := server.NewMCPServer("strava-mcp", "1.0.0")
	tools.RegisterStravaTools(s, client)

	stdio := server.NewStdioServer(s)
	if err := stdio.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func runAuthFlow() {
	clientID := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("STRAVA_CLIENT_ID and STRAVA_CLIENT_SECRET environment variables are required")
	}

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			errCh <- fmt.Errorf("missing authorization code in callback")
			return
		}
		fmt.Fprint(w, "<html><body><h1>Authorization successful!</h1><p>You can close this window.</p></body></html>")
		codeCh <- code
	})

	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- fmt.Errorf("callback server error: %w", err)
		}
	}()

	authURL := fmt.Sprintf(
		"https://www.strava.com/oauth/authorize?client_id=%s&redirect_uri=http://localhost:8080/callback&response_type=code&scope=read_all,activity:read_all,profile:read_all",
		clientID,
	)

	fmt.Printf("Opening browser for Strava authorization...\n")
	fmt.Printf("If the browser doesn't open, visit:\n%s\n\n", authURL)
	openBrowser(authURL)

	select {
	case code := <-codeCh:
		fmt.Println("Exchanging authorization code for token...")
		token, err := strava.ExchangeCode(clientID, clientSecret, code)
		if err != nil {
			log.Fatalf("Failed to exchange code: %v", err)
		}
		fmt.Printf("Authentication successful! Token saved (expires at: %d)\n", token.ExpiresAt)

	case err := <-errCh:
		log.Fatalf("Authentication failed: %v", err)
	}

	srv.Shutdown(context.Background())
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return
	}
	cmd.Start()
}
