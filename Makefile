.PHONY: build
build:
	go build -o payjp-cli .

.PHONY: build-with-version
build-with-version:
	@VERSION=$${VERSION:-dev} && \
	COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
	DATE=$$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
	LDFLAGS="-X github.com/payjp/payjp-cli/internal/version.Version=$$VERSION -X github.com/payjp/payjp-cli/internal/version.Commit=$$COMMIT -X github.com/payjp/payjp-cli/internal/version.Date=$$DATE" && \
	go build -ldflags "$$LDFLAGS" -o payjp-cli .