# Prodyo Backend API

A REST API for managing projects and users with automatic database migrations.

## Features

- **Project Management**: Create, read, update, and delete projects
- **User Management**: Create, read, update, and delete users
- **Project Members**: Many-to-many relationship between projects and users
- **Automatic Migrations**: Database migrations run automatically on startup
- **Productivity Ranges**: Configure productivity thresholds for projects
- **Docker Support**: Easy deployment with Docker Compose

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

## License

MIT
