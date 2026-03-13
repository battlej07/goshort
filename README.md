# goshort
A URL shortener written in Go.

## Development

1. Copy the example environment file:

   ```sh
   cp .example.env .env
   ```

2. Start Postgres for local development:

   ```sh
   docker compose -f compose.dev.yml up -d
   ```

3. Run the app:

   ```sh
   go run ./cmd/server
   ```

The server will create the `short_urls` table automatically if it does not already exist.
