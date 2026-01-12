.PHONY: up down migrate run test secrets-keygen secrets-encrypt secrets-decrypt

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	go run cmd/migrate/main.go up

run: secrets-decrypt
	go run cmd/bot/main.go

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
