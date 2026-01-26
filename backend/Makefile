# Database configuration
DB_URL=postgres://bayarin_user:bayarin_pass_2024@localhost:5432/bayarin_db?sslmode=disable
DB_POSTGRES_URL=postgres://postgres@localhost:5432/postgres?sslmode=disable

.PHONY: migrate-up migrate-down migrate-create db-reset db-setup db-drop db-test

# Setup database (create user & database)
db-setup:
	@echo "Setting up database..."
	psql -U postgres -f scripts/setup_db.sql
	@echo "✅ Database setup complete!"

# Drop database
db-drop:
	@echo "Dropping database..."
	psql -U postgres -c "DROP DATABASE IF EXISTS bayarin_db;"
	psql -U postgres -c "DROP USER IF EXISTS bayarin_user;"
	@echo "✅ Database dropped!"

# Test database connection
db-test:
	@echo "Testing database connection..."
	psql "$(DB_URL)" -c "SELECT current_database(), current_user;"
	@echo "✅ Connection successful!"

# Run migrations up
migrate-up:
	@echo "Running migrations..."
	migrate -path migrations -database "$(DB_URL)" up
	@echo "✅ Migrations complete!"

# Run migrations down
migrate-down:
	@echo "Rolling back migrations..."
	migrate -path migrations -database "$(DB_URL)" down
	@echo "✅ Rollback complete!"

# Create new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Reset database (drop + setup + migrate)
db-reset:
	@echo "Resetting database..."
	make db-drop
	make db-setup
	make migrate-up
	@echo "✅ Database reset complete!"

# Database seed
db-seed:
	@echo "Seeding database..."
	psql "$(DB_URL)" < migrations/seed.sql
	@echo "✅ Seed complete!"

# Run application
run:
	go run cmd/api/main.go

# Build application
build:
	go build -o bin/bayarin cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf logs/*.log