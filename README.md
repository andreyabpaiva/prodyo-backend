# Prodyo Backend API

A REST API for managing projects and users with automatic database migrations.

## Features

- **Project Management**: Create, read, update, and delete projects
- **User Management**: Create, read, update, and delete users
- **Project Members**: Many-to-many relationship between projects and users
- **Automatic Migrations**: Database migrations run automatically on startup
- **Productivity Ranges**: Configure productivity thresholds for projects
- **Docker Support**: Easy deployment with Docker Compose

## Project Entity Structure

The Project entity has the following attributes:
- `id`: UUID (Primary Key)
- `name`: string (Required)
- `description`: string
- `members`: User[] (Array of users)
- `color`: string
- `prodRange`: ProductivityRange
  - `ok`: int
  - `alert`: int
  - `critical`: int
- `createdAt`: Time
- `updatedAt`: Time

## User Entity Structure

The User entity has the following attributes:
- `id`: UUID (Primary Key)
- `name`: string (Required)
- `email`: string (Required, Unique)
- `createdAt`: Time
- `updatedAt`: Time

## Quick Start

### Using Docker Compose (Recommended)

1. **Start the services:**
```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- PgAdmin on port 8080
- API server on port 8081

2. **Access the API:**
- API: http://localhost:8081
- Swagger UI: http://localhost:8081/swagger/
- PgAdmin: http://localhost:8080

### Manual Setup

1. **Install dependencies:**
```bash
go mod tidy
```

2. **Set environment variables:**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=prodyo_user
export DB_PASSWORD=prodyo_password
export DB_NAME=prodyo_db
```

3. **Start PostgreSQL database**

4. **Run the API:**
```bash
go run cmd/api/main.go
```

## API Endpoints

### Base URL
```
http://localhost:8081/api/v1
```

### Project Endpoints

#### Get All Projects
- **Method**: `GET`
- **URL**: `/projects`
- **Description**: Retrieve all projects with their members

**Example Response:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "My Project",
    "description": "A sample project",
    "color": "#FF5733",
    "prod_range": {
      "ok": 80,
      "alert": 60,
      "critical": 40
    },
    "members": [
      {
        "id": "456e7890-e89b-12d3-a456-426614174001",
        "name": "John Doe",
        "email": "john@example.com",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

#### Get Project by ID
- **Method**: `GET`
- **URL**: `/projects/{id}`
- **Description**: Retrieve a specific project by its ID

#### Create Project
- **Method**: `POST`
- **URL**: `/projects`
- **Description**: Create a new project

**Request Body:**
```json
{
  "name": "New Project",
  "description": "Project description",
  "color": "#FF5733",
  "prod_range": {
    "ok": 80,
    "alert": 60,
    "critical": 40
  },
  "member_ids": ["456e7890-e89b-12d3-a456-426614174001"]
}
```

#### Update Project
- **Method**: `PUT`
- **URL**: `/projects/{id}`
- **Description**: Update an existing project

#### Delete Project
- **Method**: `DELETE`
- **URL**: `/projects/{id}`
- **Description**: Delete a project

### User Endpoints

#### Get All Users
- **Method**: `GET`
- **URL**: `/users`
- **Description**: Retrieve all users

#### Get User by ID
- **Method**: `GET`
- **URL**: `/users/{id}`
- **Description**: Retrieve a specific user by ID

#### Create User
- **Method**: `POST`
- **URL**: `/users`
- **Description**: Create a new user

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

#### Update User
- **Method**: `PUT`
- **URL**: `/users/{id}`
- **Description**: Update an existing user

#### Delete User
- **Method**: `DELETE`
- **URL**: `/users/{id}`
- **Description**: Delete a user

### Health Check
- **Method**: `GET`
- **URL**: `/health`
- **Description**: Check API health status

## Database Migrations

The API automatically runs database migrations on startup. Migrations are located in the `cmd/migrations/` directory:

1. `001_init_projects_table.up.sql` - Creates the initial projects table
2. `002_create_users_table.up.sql` - Creates the users table
3. `003_update_projects_table.up.sql` - Updates projects table to add members relationship

### Adding New Migrations

To add a new migration:

1. Create a new migration file in `cmd/migrations/` with the format:
   - `{number}_{description}.up.sql`
   - `{number}_{description}.down.sql`

2. The migration will run automatically on the next API startup

## Testing

Use the provided test script to test the API:

```bash
./test_api.sh
```

Make sure you have `jq` installed for JSON formatting:
```bash
# On Ubuntu/Debian
sudo apt-get install jq

# On macOS
brew install jq

# On Windows (with Chocolatey)
choco install jq
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `postgres` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `prodyo_user` | Database username |
| `DB_PASSWORD` | `prodyo_password` | Database password |
| `DB_NAME` | `prodyo_db` | Database name |

## Development

### Project Structure
```
prodyo-backend/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── migrations/
│       ├── 001_init_projects_table.up.sql
│       ├── 001_init_projects_table.down.sql
│       ├── 002_create_users_table.up.sql
│       ├── 002_create_users_table.down.sql
│       ├── 003_update_projects_table.up.sql
│       └── 003_update_projects_table.down.sql
├── cmd/internal/
│   ├── config/
│   ├── handlers/
│   ├── migrations/
│   ├── models/
│   ├── repositories/
│   └── usecases/
├── docs/
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

### Adding New Features

1. Create the model in `cmd/internal/models/`
2. Create the repository in `cmd/internal/repositories/`
3. Create the use case in `cmd/internal/usecases/`
4. Create the handler in `cmd/internal/handlers/`
5. Add routes in `cmd/internal/handlers/routes.go`
6. Create database migration in `cmd/migrations/`
7. Update `cmd/internal/repositories/repositories.go` to include new repository

## License

MIT
