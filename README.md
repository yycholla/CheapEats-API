# CheapEats API

A Go-based REST API for fetching and tracking restaurant prices using Google Places API.

## Features

- Fetch restaurants from Google Places API
- Store restaurant and menu data in PostgreSQL with GORM ORM
- Track price history for menu items
- Search restaurants by location, cuisine, and price range
- RESTful API built with Chi router
- Swagger/OpenAPI documentation with interactive UI
- Docker support for easy deployment
- Hot reload development with Air

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Google Places API key

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and add your Google Places API key
3. Run with Docker Compose:

```bash
docker-compose up -d
```

## API Endpoints

### Health Check
- `GET /api/v1/health` - Check API health status

### Restaurants
- `GET /api/v1/restaurants` - Get all restaurants
  - Query params: `city`, `cuisine`, `price_range`
- `GET /api/v1/restaurants/search` - Search nearby restaurants
  - Query params: `lat`, `lng`, `radius` (in meters)
- `GET /api/v1/restaurants/{id}` - Get restaurant details
- `GET /api/v1/restaurants/{id}/menu` - Get restaurant menu items
  - Query params: `category`, `max_price`

### Menu Items
- `GET /api/v1/menu-items/{itemId}` - Get menu item details
- `GET /api/v1/menu-items/{itemId}/price-history` - Get price history for item

## Swagger Documentation

The API includes interactive Swagger documentation. After starting the server, visit:
- Swagger UI: `http://localhost:8080/swagger/index.html`
- OpenAPI JSON: `http://localhost:8080/swagger/doc.json`

To regenerate Swagger docs:
```bash
make swagger
# or
swag init -g cmd/api/main.go --output docs
```

## Development

### Using Makefile

```bash
# Display all available commands
make help

# Install dependencies and tools
make install

# Generate Swagger documentation
make swagger

# Run the application locally
make run

# Build the application
make build

# Start Docker containers
make docker-up

# Stop Docker containers
make docker-down
```

### Local Development

```bash
# Install dependencies
go mod download

# Install Swag for documentation
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
swag init -g cmd/api/main.go --output docs

# Run database
docker-compose up postgres -d

# Run the API with hot reload
air

# Or run without hot reload
go run cmd/api/main.go
```

### Building

```bash
# Build Docker image
docker build -t cheapeats-api .

# Or build locally
go build -o bin/api cmd/api/main.go
```

## Database Schema

The API uses PostgreSQL with the following main tables:
- `restaurants` - Restaurant information
- `menu_items` - Menu items with prices
- `price_history` - Historical price tracking
- `scraped_data` - Raw API response storage

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | Database user | cheapeats |
| DB_PASSWORD | Database password | cheapeats_pass |
| DB_NAME | Database name | cheapeats_db |
| DB_SSLMODE | SSL mode | disable |
| GOOGLE_PLACES_API_KEY | Google Places API key | (required) |

## Notes

- The Google Places API integration currently generates sample menu items with prices based on the restaurant's price level
- For production use, consider implementing actual menu scraping or partnering with restaurant data providers
- Rate limiting is implemented with a 100ms delay between API calls to respect Google's usage limits