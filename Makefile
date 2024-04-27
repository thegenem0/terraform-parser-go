setup: 
	@echo "Setting up PostgreSQL for local development"
	@docker-compose up -d
	@echo "Setting up API"
	@go mod download
	@go run .

