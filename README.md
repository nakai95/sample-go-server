# sample-go-server

This is a sample web application.

## Project Structure

```
sample-go-server
├── api                      # OpenAPI
├── build                    # Docker configuration
├── cmd
│   └── sample-go-server
│       └── main.go          # Entry point of the application
├── internal
│   ├── adapter              # Interface adapters
│   ├── domain               # Domain model
│   ├── infrastructure       # Handlers and router
│   └── usecase              # Business use cases
├── mock                     # Mocks used for testing
├── go.mod                   # Module definition and dependencies
├── go.sum                   # Checksums for module dependencies
└── compose.yaml             # Docker Compose configuration
```

## Setup Instructions

1. **Clone the repository:**

   ```
   git clone https://github.com/nakai95/sample-go-server.git
   cd sample-go-server
   ```

2. **Install dependencies:**

   ```
   go mod tidy
   ```

3. **Run the application:**
   ```
   go run cmd/sample-go-server/main.go
   ```

## Usage

- **Note**: Only the following credentials are accepted for authentication:
  - **Username**: `demo@example.com`
  - **Password**: `#demo`

### Endpoints

- `POST /auth/token`: Returns a JWT token for a given username and password.

  - **Request Body**:
    - `Content-Type`: `application/x-www-form-urlencoded`
    - `username` (string, required): The username for authentication.
    - `password` (string, required): The password for authentication.
  - **Example Request**:
    ```bash
      curl -X 'POST' \
        'http://localhost:8080/auth/token' \
        -H 'accept: application/json' \
        -H 'Content-Type: application/x-www-form-urlencoded' \
        -d 'username=demo%40example.com&password=%23demo'
    ```
  - **Response**: A JSON object containing the JWT token.
  - **Example Response**:
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

- `GET /events`: List all events.
  - **Note**: This endpoint returns demo data only.
  - **Example Request**:
    ```bash
      curl -X 'GET' \
       'http://localhost:8080/events' \
        -H 'accept: application/json' \
        -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
    ```
  - **Response**: A JSON array of events.
  - **Example**:
    ```json
    [
      {
        "id": "1",
        "name": "Event 1",
        "description": "Description of Event 1",
        "imageUrl": "https://example.com/event1.jpg"
      },
      {
        "id": "2",
        "name": "Event 2",
        "description": "Description of Event 2",
        "imageUrl": "https://example.com/event2.jpg"
      }
    ]
    ```
