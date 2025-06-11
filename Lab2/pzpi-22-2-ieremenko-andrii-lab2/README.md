# Go HTTP Client with JWT Authentication

This is a basic Go HTTP client implementation that provides user registration and JWT authentication functionality.

## Features

- User registration
- User login with JWT authentication
- Get user information
- JWT token management
- Swagger API documentation

## Project Structure

```
.
├── auth/
│   └── jwt.go         # JWT authentication utilities
├── client/
│   └── client.go      # HTTP client implementation
├── models/
│   └── user.go        # User-related data structures
├── server/
│   └── server.go      # HTTP server implementation
├── docs/
│   └── swagger.go     # Swagger documentation
├── main.go            # Example usage
├── go.mod             # Go module file
└── README.md          # This file
```

## Usage

1. First, make sure you have Go 1.21 or later installed.

2. Install the dependencies:
```bash
go mod tidy
```

3. Install Swagger CLI tool:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

4. Generate Swagger documentation:
```bash
swag init
```

5. Run the server:
```bash
go run main.go
```

6. Access the Swagger UI at `http://localhost:8080/swagger/index.html`

## API Documentation

The API is documented using Swagger/OpenAPI. You can access the Swagger UI to:
- View all available endpoints
- Test the API directly from the browser
- View request/response schemas
- See authentication requirements

## API Endpoints

The client expects the following endpoints to be available on the server:

- `POST /register` - Register a new user
- `POST /login` - Login and get JWT token
- `GET /user` - Get user information (requires JWT token)

## Security Notes

- In a production environment, always use HTTPS
- Store the JWT secret key in environment variables
- Implement proper password hashing on the server side
- Add rate limiting and other security measures as needed 