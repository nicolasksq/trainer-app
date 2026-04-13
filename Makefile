.PHONY: setup build test-strava test-garmin auth-strava clean

setup:
	go run ./cmd/setup

build:
	go build -o bin/strava-mcp ./cmd/strava-mcp
	go build -o bin/garmin-mcp ./cmd/garmin-mcp

test-strava:
	@printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}\n{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}\n' | go run ./cmd/strava-mcp 2>/dev/null | python3 -c "import sys,json; [print(f\"  {t['name']:30s} {t['description'][:60]}\") for t in json.loads(sys.stdin.readlines()[-1])['result']['tools']]" || echo "Failed - run 'make setup' first"

test-garmin:
	@printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}\n{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}\n' | go run ./cmd/garmin-mcp 2>/dev/null | python3 -c "import sys,json; [print(f\"  {t['name']:30s} {t['description'][:60]}\") for t in json.loads(sys.stdin.readlines()[-1])['result']['tools']]" || echo "Failed - run 'make setup' first"

auth-strava:
	go run ./cmd/strava-mcp auth

clean:
	rm -rf bin/
