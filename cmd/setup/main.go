package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nicolas-andreoli/trainer-app/internal/envutil"
	"github.com/nicolas-andreoli/trainer-app/internal/garmin"
)

const (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	bold   = "\033[1m"
	reset  = "\033[0m"
)

var scanner = bufio.NewScanner(os.Stdin)

type stepResult struct {
	name   string
	passed bool
}

var results []stepResult

func main() {
	printBanner()

	checkGo()
	checkClaude()
	configureEnv()
	stravaOAuth()
	testGarmin()
	testMCPServers()
	printSummary()
}

func printBanner() {
	fmt.Println(bold + "╔══════════════════════════════════════╗")
	fmt.Println("║       Trainer App Setup Wizard       ║")
	fmt.Println("╚══════════════════════════════════════╝" + reset)
	fmt.Println()
	fmt.Println("This wizard will:")
	fmt.Println("  1. Verify prerequisites (Go, Claude Code)")
	fmt.Println("  2. Configure API credentials (.env)")
	fmt.Println("  3. Authorize Strava OAuth")
	fmt.Println("  4. Test Garmin login")
	fmt.Println("  5. Verify MCP servers")
	fmt.Println()
}

func checkGo() {
	fmt.Println(bold + "Step 1: Checking Go installation..." + reset)
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		fmt.Println(red + "  Go is required. Install from https://go.dev/dl/" + reset)
		results = append(results, stepResult{"Go installed", false})
		os.Exit(1)
	}
	fmt.Printf(green+"  %s"+reset+"\n\n", strings.TrimSpace(string(out)))
	results = append(results, stepResult{"Go installed", true})
}

func checkClaude() {
	fmt.Println(bold + "Step 2: Checking Claude Code..." + reset)
	out, err := exec.Command("claude", "--version").Output()
	if err != nil {
		fmt.Println(yellow + "  Warning: Claude Code not found. It is needed to use the @trainer agent." + reset)
		fmt.Println(yellow + "  Install: https://docs.anthropic.com/en/docs/claude-code" + reset)
		results = append(results, stepResult{"Claude Code", false})
	} else {
		fmt.Printf(green+"  Claude Code %s"+reset+"\n", strings.TrimSpace(string(out)))
		results = append(results, stepResult{"Claude Code", true})
	}
	fmt.Println()
}

func configureEnv() {
	fmt.Println(bold + "Step 3: Configuring .env file..." + reset)

	if _, err := os.Stat(".env"); err == nil {
		fmt.Print("  .env already exists. Reconfigure? (y/n): ")
		if !askYesNo() {
			fmt.Println("  Skipping .env configuration.")
			results = append(results, stepResult{".env configured", true})
			fmt.Println()
			return
		}
	}

	stravaID := prompt("  Strava Client ID (create app at https://www.strava.com/settings/api): ")
	stravaSecret := prompt("  Strava Client Secret: ")
	garminEmail := prompt("  Garmin Email: ")
	garminPassword := prompt("  Garmin Password (stored in .env, which is gitignored): ")

	content := fmt.Sprintf(`# Strava OAuth2 credentials
STRAVA_CLIENT_ID=%s
STRAVA_CLIENT_SECRET=%s

# Garmin Connect credentials
GARMIN_EMAIL=%s
GARMIN_PASSWORD=%s
`, stravaID, stravaSecret, garminEmail, garminPassword)

	if err := os.WriteFile(".env", []byte(content), 0600); err != nil {
		fmt.Printf(red+"  Failed to write .env: %v"+reset+"\n", err)
		results = append(results, stepResult{".env configured", false})
	} else {
		fmt.Println(green + "  .env file written successfully." + reset)
		results = append(results, stepResult{".env configured", true})
	}
	fmt.Println()
}

func stravaOAuth() {
	fmt.Println(bold + "Step 4: Strava OAuth authorization..." + reset)
	fmt.Print("  Run Strava authorization now? (y/n): ")
	if !askYesNo() {
		fmt.Println("  Skipping Strava OAuth. Run later with: go run ./cmd/strava-mcp auth")
		results = append(results, stepResult{"Strava OAuth", false})
		fmt.Println()
		return
	}

	cmd := exec.Command("go", "run", "./cmd/strava-mcp", "auth")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf(red+"  Strava OAuth failed: %v"+reset+"\n", err)
		results = append(results, stepResult{"Strava OAuth", false})
	} else {
		fmt.Println(green + "  Strava OAuth completed." + reset)
		results = append(results, stepResult{"Strava OAuth", true})
	}
	fmt.Println()
}

func testGarmin() {
	fmt.Println(bold + "Step 5: Testing Garmin login..." + reset)
	envutil.LoadDotEnv(".env")

	client := garmin.NewClient()
	if err := client.EnsureAuthenticated(); err != nil {
		fmt.Printf(red+"  Garmin login failed: %v"+reset+"\n", err)
		results = append(results, stepResult{"Garmin login", false})
	} else {
		fmt.Println(green + "  Garmin authentication successful." + reset)
		results = append(results, stepResult{"Garmin login", true})
	}
	fmt.Println()
}

func testMCPServers() {
	fmt.Println(bold + "Step 6: Testing MCP servers..." + reset)
	for _, srv := range []string{"strava-mcp", "garmin-mcp"} {
		ok := testMCPServer(srv)
		results = append(results, stepResult{srv + " server", ok})
	}
	fmt.Println()
}

func testMCPServer(name string) bool {
	fmt.Printf("  Testing %s...\n", name)

	cmd := exec.Command("go", "run", "./cmd/"+name)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf(red+"    Failed to create stdin pipe: %v"+reset+"\n", err)
		return false
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf(red+"    Failed to create stdout pipe: %v"+reset+"\n", err)
		return false
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf(red+"    Failed to start %s: %v"+reset+"\n", name, err)
		return false
	}

	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	initMsg := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"setup","version":"1.0.0"}}}` + "\n"
	if _, err := io.WriteString(stdin, initMsg); err != nil {
		fmt.Printf(red+"    Failed to send initialize: %v"+reset+"\n", err)
		return false
	}

	respCh := make(chan []byte, 1)
	go func() {
		reader := bufio.NewReader(stdout)
		line, err := reader.ReadBytes('\n')
		if err == nil {
			respCh <- line
		}
	}()

	select {
	case resp := <-respCh:
		var msg map[string]interface{}
		if err := json.Unmarshal(resp, &msg); err != nil {
			fmt.Printf(red+"    Invalid JSON response: %v"+reset+"\n", err)
			return false
		}
		if _, ok := msg["result"]; ok {
			fmt.Printf(green+"    %s is working."+reset+"\n", name)
			return true
		}
		fmt.Printf(red+"    Unexpected response: %s"+reset+"\n", string(resp))
		return false
	case <-time.After(15 * time.Second):
		fmt.Printf(red+"    %s timed out."+reset+"\n", name)
		return false
	}
}

func printSummary() {
	fmt.Println(bold + "═══════════════════════════════════════" + reset)
	fmt.Println(bold + " Setup Summary" + reset)
	fmt.Println(bold + "═══════════════════════════════════════" + reset)
	for _, r := range results {
		icon := red + "✗" + reset
		if r.passed {
			icon = green + "✓" + reset
		}
		fmt.Printf("  %s  %s\n", icon, r.name)
	}
	fmt.Println()

	allPassed := true
	for _, r := range results {
		if !r.passed {
			allPassed = false
			break
		}
	}

	if allPassed {
		fmt.Println(green + bold + "Setup complete!" + reset)
	} else {
		fmt.Println(yellow + "Setup finished with some issues. Review the steps above." + reset)
	}
	fmt.Println("Open Claude Code in this directory and type: " + bold + "@trainer Analyze my last week of training" + reset)
}

func prompt(label string) string {
	fmt.Print(label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func askYesNo() bool {
	scanner.Scan()
	ans := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return ans == "y" || ans == "yes"
}
