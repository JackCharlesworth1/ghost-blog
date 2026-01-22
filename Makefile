.PHONY: help setup db-up db-down db-migrate backend frontend clean

help:
	@echo "Ghost Blog - Development Commands"
	@echo ""
	@echo "  make setup       - Set up the project (copy env files)"
	@echo "  make db-up       - Start PostgreSQL with Docker Compose"
	@echo "  make db-down     - Stop PostgreSQL"
	@echo "  make db-migrate  - Run database migrations"
	@echo "  make backend     - Run backend server"
	@echo "  make frontend    - Run frontend dev server"
	@echo "  make clean       - Clean build artifacts"

setup:
	@echo "Setting up environment files..."
	@cp -n backend/.env.example backend/.env 2>/dev/null || echo "backend/.env already exists"
	@cp -n frontend/.env.example frontend/.env 2>/dev/null || echo "frontend/.env already exists"
	@echo "Setup complete! Edit .env files as needed."

db-up:
	@echo "Starting PostgreSQL..."
	@docker-compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Database is ready!"

db-down:
	@echo "Stopping PostgreSQL..."
	@docker-compose down

db-migrate:
	@echo "Running migrations..."
	@docker exec -i ghostblog-postgres psql -U postgres -d ghostblog < backend/migrations/001_create_blogposts.up.sql
	@docker exec -i ghostblog-postgres psql -U postgres -d ghostblog < backend/migrations/002_create_ratelimits.up.sql
	@echo "Migrations complete!"

backend:
	@echo "Starting backend server..."
	@cd backend && go run cmd/server/main.go

frontend:
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
	@echo "Starting frontend dev server..."
	@cd frontend && npm run dev

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf backend/server
	@rm -rf frontend/dist
	@rm -rf frontend/node_modules
	@echo "Clean complete!"
