# Dotfiles Config Sharing API

Simple Go web API for sharing dotfiles configurations. Designed for Railway deployment.

## Features

- Upload/download dotfiles configurations
- Search public configurations
- Featured configurations (most downloaded)
- Simple in-memory storage (use database for production)
- CORS enabled for web frontend integration

## API Endpoints

- `POST /api/configs/upload` - Upload a new config
- `GET /api/configs/:id` - Download a config by ID
- `GET /api/configs/search?q=query` - Search public configs
- `GET /api/configs/featured` - Get featured configs
- `GET /api/configs/stats` - Get platform statistics

## Environment Variables

- `PORT` - Server port (default: 8080, automatically set by Railway)

## Local Development

```bash
go mod tidy
go run main.go
```

Server will start on http://localhost:8080

## Railway Deployment

1. Connect your GitHub repo to Railway
2. Railway will automatically detect this as a Go app
3. Set any environment variables in Railway dashboard
4. Deploy!

## Testing

```bash
# Upload a config
curl -X POST http://localhost:8080/api/configs/upload \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Config",
    "description": "A test configuration",
    "author": "testuser",
    "tags": ["test", "demo"],
    "config": "{\"brews\":[\"git\",\"curl\"],\"casks\":[\"vscode\"],\"taps\":[],\"stow\":[\"vim\"]}",
    "public": true
  }'

# Search configs
curl http://localhost:8080/api/configs/search?q=test

# Get featured configs
curl http://localhost:8080/api/configs/featured
```