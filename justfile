# Run goshort
run:
  @echo "Starting goshort..."
  go run ./cmd/server

# Set up local environment
setup:
  @echo "Setting up local environment..."
  cp ./.example.env .env

tidy:
  @echo "Tidying module dependencies..."
  go mod tidy
  @echo "Verifying module dependencies..."
  go mod verify
  @echo 'Formatting .go files...'
  go fmt ./...
