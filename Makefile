.PHONY: up down migrate run run-api run-api-dev run-all build build-frontend frontend-dev frontend-install test secrets-keygen secrets-encrypt secrets-decrypt

up:
	docker-compose up -d

down:
	docker-compose down

run: secrets-decrypt
	go run cmd/bot/main.go

run-api: secrets-decrypt
	go run cmd/api/main.go

run-api-dev: secrets-decrypt
	DEV_MODE=true DEV_USER_ID=123456789 go run cmd/api/main.go

# Build everything
build: build-frontend
	go build -o bin/api cmd/api/main.go
	go build -o bin/bot cmd/bot/main.go

build-frontend:
	cd frontend && npm install && npm run build

frontend-dev:
	cd frontend && npm run dev

frontend-install:
	cd frontend && npm install

test:
	go test -v ./...

# === Secrets (age encryption) ===

# Generate new age keypair
secrets-keygen:
	@if [ -f .age-key.txt ]; then \
		echo "Key already exists: .age-key.txt"; \
		exit 1; \
	fi
	age-keygen -o .age-key.txt
	@echo "Public key for .age-recipients:"
	@grep "public key:" .age-key.txt | cut -d: -f2 | tr -d ' '

# Encrypt .secret.env -> .secret.enc.env
secrets-encrypt:
	@if [ ! -f .age-recipients ]; then \
		echo "Create .age-recipients with public keys first"; \
		exit 1; \
	fi
	age -R .age-recipients -o .secret.enc.env .secret.env
	@echo "Encrypted: .secret.env -> .secret.enc.env"

# Decrypt .secret.enc.env -> .secret.env
secrets-decrypt:
	@if [ ! -f .age-key.txt ]; then \
		echo "No key found: .age-key.txt"; \
		exit 1; \
	fi
	age -d -i .age-key.txt -o .secret.env .secret.enc.env
	@echo "Decrypted: .secret.enc.env -> .secret.env"
