# Run goshort
run:
  @echo "Starting goshort"
  go run ./cmd/server

setup:
  @echo "Setting up local environment"
  cp ./.example.env .env
