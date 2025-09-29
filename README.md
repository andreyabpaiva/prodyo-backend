# Project API Endpoints

This document describes the available API endpoints for the Project entity.

## Base URL
```
http://localhost:8081/api/v1
```

## Endpoints

### 1. Get All Projects
- **Method**: `GET`
- **URL**: `/projects`
- **Description**: Retrieve all projects
- **Response**: Array of project objects

**Example Response:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "My Project",
    "email": "project@example.com"
  }
]
```

### 2. Get Project by ID
- **Method**: `GET`
- **URL**: `/projects/{id}`
- **Description**: Retrieve a specific project by its ID
- **Parameters**: 
  - `id` (path): UUID of the project
- **Response**: Project object

**Example Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "My Project",
  "email": "project@example.com"
}
```

### 3. Create Project
- **Method**: `POST`
- **URL**: `/projects`
- **Description**: Create a new project
- **Request Body**:
```json
{
  "name": "My New Project",
  "email": "newproject@example.com"
}
```
- **Response**: Created project object with generated ID

**Example Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "My New Project",
  "email": "newproject@example.com"
}
```

### 4. Update Project
- **Method**: `PUT`
- **URL**: `/projects/{id}`
- **Description**: Update an existing project
- **Parameters**: 
  - `id` (path): UUID of the project to update
- **Request Body**:
```json
{
  "name": "Updated Project Name",
  "email": "updated@example.com"
}
```
- **Response**: Updated project object

### 5. Delete Project
- **Method**: `DELETE`
- **URL**: `/projects/{id}`
- **Description**: Delete a project
- **Parameters**: 
  - `id` (path): UUID of the project to delete
- **Response**: No content (204 status)

### 6. Health Check
- **Method**: `GET`
- **URL**: `/health`
- **Description**: Check if the API is running
- **Response**:
```json
{
  "status": "healthy"
}
```

## Error Responses

All endpoints may return the following error responses:

- **400 Bad Request**: Invalid request data or malformed UUID
- **404 Not Found**: Project not found
- **500 Internal Server Error**: Server error

## Swagger Documentation

The API includes interactive Swagger documentation that is automatically generated and served when the server starts.

### Accessing Swagger UI

Once the server is running, you can access the Swagger UI at:
```
http://localhost:8081/swagger/index.html
```

The Swagger UI provides:
- Interactive API documentation
- Try-it-out functionality for all endpoints
- Request/response examples
- Schema definitions
- Authentication testing (if implemented)

### Regenerating Documentation

If you make changes to the API annotations, regenerate the Swagger documentation:

```bash
swag init -g cmd/api/main.go -o docs
```

## Docker Setup

This project includes Docker configuration for PostgreSQL database and pgAdmin for database management.

### Prerequisites

- Docker and Docker Compose installed on your system

### Starting the Database Services

1. **Start PostgreSQL and pgAdmin:**
   ```bash
   docker-compose up -d
   ```

2. **Check if services are running:**
   ```bash
   docker-compose ps
   ```

### Database Access

- **PostgreSQL Database:**
  - Host: `localhost`
  - Port: `5432`
  - Database: `prodyo_db`
  - Username: `prodyo_user`
  - Password: `prodyo_password`

- **pgAdmin Web Interface:**
  - URL: `http://localhost:8080`
  - Email: `admin@prodyo.com`
  - Password: `admin123`

### Connecting to Database in pgAdmin

1. Open pgAdmin at `http://localhost:8080`
2. Login with the credentials above
3. Right-click "Servers" → "Create" → "Server"
4. In the "General" tab, enter a name (e.g., "Prodyo DB")
5. In the "Connection" tab, enter:
   - Host name/address: `postgres` (use the Docker service name)
   - Port: `5432`
   - Username: `prodyo_user`
   - Password: `prodyo_password`
6. Click "Save"

### Environment Configuration

The project uses environment variables for database configuration. You can modify the `.env` file to change database settings:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=prodyo_user
DB_PASSWORD=prodyo_password
DB_NAME=prodyo_db
```

### Stopping Services

To stop the database services:
```bash
docker-compose down
```

To stop and remove all data (volumes):
```bash
docker-compose down -v
```

## Running the API

To start the API server:

```bash
go run cmd/api/main.go
```

The server will start on port 8081 by default and will display all available endpoints including the Swagger UI.

**Note:** Make sure the PostgreSQL database is running (via Docker Compose) before starting the API server.
