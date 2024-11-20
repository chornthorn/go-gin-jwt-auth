# JWT Authentication Service

A Go-based REST API service implementing JWT authentication with refresh tokens using Gin framework and GORM.

## Features

- User registration and authentication
- JWT-based authentication with access and refresh tokens
- RSA-based token signing
- User profile management
- PostgreSQL database integration
- Middleware for protected routes

## Prerequisites

- Go 1.22 or higher
- PostgreSQL

## Environment Setup

1. Clone the repository
```bash
git clone https://github.com/chornthorn/jwt-auth-app.git
cd jwt-auth-app
```

2. Create a `.env` file in the root directory by copying the `.env.example` file.

3. Generate RSA key pairs for JWT signing:
```bash
# Create keys directory
mkdir -p keys

# Generate access token keys
openssl genrsa -out keys/access_private.pem 2048
openssl rsa -in keys/access_private.pem -pubout -out keys/access_public.pem

# Generate refresh token keys
openssl genrsa -out keys/refresh_private.pem 2048
openssl rsa -in keys/refresh_private.pem -pubout -out keys/refresh_public.pem
```

## Running the Application

1. Install dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run main.go
```

## API Endpoints

### Public Routes
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get tokens
- `POST /api/v1/auth/refresh` - Refresh access token (requires refresh token)

### Protected Routes
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile
- `GET /api/v1/token/info` - Get token information

## Example Requests

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```


## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details

