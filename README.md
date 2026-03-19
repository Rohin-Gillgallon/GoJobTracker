# Go Job Tracker API

A RESTful API built in Go for tracking job applications. Features JWT authentication, full CRUD operations, filtering, and pagination.

🚀 **Live API**: https://gojobtracker-production.up.railway.app

## Tech Stack

- **Language**: Go 1.24
- **Router**: Chi
- **Database**: PostgreSQL + sqlx
- **Auth**: JWT (access + refresh tokens) + bcrypt
- **Migrations**: golang-migrate
- **Containerisation**: Docker + Docker Compose
- **CI/CD**: GitHub Actions + Railway

## Endpoints

### Auth
| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/auth/register` | Register a new user | No |
| POST | `/auth/login` | Login and receive tokens | No |

### Jobs
| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/jobs` | Create a job application | Yes |
| GET | `/jobs` | Get all job applications | Yes |
| GET | `/jobs/{id}` | Get a job application by ID | Yes |
| PUT | `/jobs/{id}` | Update a job application | Yes |
| DELETE | `/jobs/{id}` | Delete a job application | Yes |

### Query Parameters
| Parameter | Type | Description |
|---|---|---|
| `status` | string | Filter by status: `applied`, `interview`, `offer`, `rejected` |
| `page` | int | Page number (default: 1) |
| `limit` | int | Results per page (default: 10, max: 100) |

## Getting Started

### Prerequisites
- Go 1.24+
- Docker Desktop
- golang-migrate

### Local Development

1. Clone the repository
```bash
git clone https://github.com/Rohin-Gillgallon/GoJobTracker.git
cd GoJobTracker
```

2. Copy the example env file
```bash
cp .env.example .env
```

3. Start the database
```bash
docker-compose up -d postgres
```

4. Run migrations
```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/jobtracker?sslmode=disable" up
```

5. Run the server
```bash
go run ./cmd/server
```

The API will be available at `http://localhost:8080`.

### Running Tests

Start the test database first:
```bash
docker-compose -f docker-compose.test.yml up -d
```

Then run the tests:
```bash
go test ./... -v -race
```

### Docker

Run the full stack with Docker Compose:
```bash
docker-compose up --build
```

## CI/CD

- **CI** (`ci.yml`) — runs on every push: lints with `golangci-lint`, runs tests with race detection against a real Postgres instance, builds the binary
- **Deploy** (`deploy.yml`) — runs on merge to `main`: builds and pushes Docker image to GitHub Container Registry. Railway auto-deploys from the connected GitHub repo.

## Example Usage

### Register
```bash
curl -X POST https://gojobtracker-production.up.railway.app/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"yourpassword"}'
```

### Login
```bash
curl -X POST https://gojobtracker-production.up.railway.app/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"yourpassword"}'
```

### Create a Job Application
```bash
curl -X POST https://gojobtracker-production.up.railway.app/jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"company":"Acme Corp","role":"Software Engineer","status":"applied"}'
```

### Get All Jobs (with filtering)
```bash
curl https://gojobtracker-production.up.railway.app/jobs?status=applied&page=1&limit=10 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Project Structure
```
├── cmd/server          # Entry point
├── internal/
│   ├── auth            # JWT generation & middleware
│   ├── config          # Environment variable loading
│   ├── database        # Database connection
│   ├── handlers        # HTTP handlers
│   ├── models          # Structs and request/response types
│   └── repository      # Database queries
├── migrations          # SQL migration files
├── .github/workflows   # CI/CD pipelines
├── Dockerfile
└── docker-compose.yml
```