.PHONY: help install install-frontend install-backend sync dev backend frontend seed createuser build build-frontend build-backend run clean test lint check fmt

PKG_MGR ?= $(shell command -v pnpm 2>/dev/null || command -v npm)
BACKEND_DIR := backend
FRONTEND_DIR := frontend
BUILD_DIR := $(BACKEND_DIR)/internal/static/dist
NODE_MODULES := $(FRONTEND_DIR)/node_modules

help:
	@echo "lynxlinkage development targets"
	@echo ""
	@echo "  make install         Install backend (go mod) + frontend (npm/pnpm) deps"
	@echo "  make dev             Run backend + frontend dev servers in parallel"
	@echo "  make backend         Run only the Go server (with air if available)"
	@echo "  make frontend        Run only the SvelteKit dev server"
	@echo "  make seed            Load backend/seed/*.yaml into the PostgreSQL database"
	@echo "  make build           Build production binary with the frontend embedded"
	@echo "  make run             Run the embedded production binary"
	@echo "  make test            Run backend tests"
	@echo "  make lint            Run backend vet + frontend lint"
	@echo "  make check           Run svelte-check"
	@echo "  make fmt             Format frontend with prettier"
	@echo "  make clean           Remove build artifacts"

install: install-backend install-frontend

install-backend:
	cd $(BACKEND_DIR) && go mod download

install-frontend: $(NODE_MODULES)

$(NODE_MODULES): $(FRONTEND_DIR)/package.json
	cd $(FRONTEND_DIR) && $(PKG_MGR) install
	@touch $(NODE_MODULES)

sync: install-frontend
	cd $(FRONTEND_DIR) && $(PKG_MGR) exec svelte-kit sync

dev: sync
	@echo "Starting backend on :8080 and frontend on :5173 (Ctrl+C to stop)…"
	@$(MAKE) -j2 backend frontend

backend:
	@if command -v air >/dev/null 2>&1; then \
		cd $(BACKEND_DIR) && air; \
	else \
		echo "(install github.com/air-verse/air for hot reload; running plain go run)"; \
		cd $(BACKEND_DIR) && go run ./cmd/server; \
	fi

frontend: sync
	cd $(FRONTEND_DIR) && $(PKG_MGR) run dev

seed:
	cd $(BACKEND_DIR) && go run ./cmd/seed -dir ./seed

# Create an HR user. Pass arguments via EMAIL=... PASSWORD=... or omit
# PASSWORD to be prompted on stdin.
createuser:
	@if [ -z "$(EMAIL)" ]; then echo "usage: make createuser EMAIL=hr@example.com [PASSWORD=...]"; exit 2; fi
	@cd $(BACKEND_DIR) && go run ./cmd/createuser -email "$(EMAIL)" $(if $(PASSWORD),-password "$(PASSWORD)",)

build: build-frontend build-backend

build-frontend: sync
	cd $(FRONTEND_DIR) && $(PKG_MGR) run build

build-backend: build-frontend
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)
	@cp -R $(FRONTEND_DIR)/build/. $(BUILD_DIR)/
	cd $(BACKEND_DIR) && go build -tags=embed -o ../bin/lynxlinkage ./cmd/server
	@echo "built ./bin/lynxlinkage with frontend embedded"

run: build
	./bin/lynxlinkage

test:
	cd $(BACKEND_DIR) && go test ./...

lint:
	cd $(BACKEND_DIR) && go vet ./...
	cd $(FRONTEND_DIR) && $(PKG_MGR) run lint || true

check:
	cd $(FRONTEND_DIR) && $(PKG_MGR) run check

fmt:
	cd $(FRONTEND_DIR) && $(PKG_MGR) run format

clean:
	rm -rf bin/ $(FRONTEND_DIR)/build $(FRONTEND_DIR)/.svelte-kit $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)
	@touch $(BUILD_DIR)/.gitkeep
