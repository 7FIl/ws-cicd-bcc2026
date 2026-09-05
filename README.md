# ws-cicd-bcc2026

Simple Todo CRUD REST API built with [Go Fiber](https://gofiber.io/), [GORM](https://gorm.io/), and PostgreSQL 18.
The project is structured with dependency injection (handler → usecase → repository).

## Prerequisites

| Tool                                                           | Version | Needed for                                   |
| -------------------------------------------------------------- | ------- | -------------------------------------------- |
| [Go](https://go.dev/dl/)                                       | 1.25+   | running the app locally                      |
| [Docker](https://docs.docker.com/get-docker/) + Docker Compose | v2      | running PostgreSQL and the containerized app |
| [Make](https://www.gnu.org/software/make/)                     | any     | the shortcut commands (optional)             |
| [golangci-lint](https://golangci-lint.run/)                    | v2.13.2 | linting (optional, `make lint-install`)      |

> On Windows, `make` is available through [Chocolatey](https://community.chocolatey.org/packages/make) (`choco install make`) or Git Bash + MSYS2. Every `make` target can also be copied and run directly.

## Project Structure

```
cmd/
  api/                      application entrypoint
  migrate/                  migration entrypoint (-m up | -m down)
internal/
  bootstrap/                env loading, fiber app, route wiring (dependency injection)
  app/todo/
    entity/                 todo model and request payloads
    handler/                http layer
    usecase/                business logic
    repository/             database access
  infra/postgresql/         gorm connection and migrations
pkg/helper/                 env reader and json response helper
```

## Configuration

Copy the example env file and adjust it if needed:

```bash
cp .env.example .env
```

| Variable      | Default                  | Description                                      |
| ------------- | ------------------------ | ------------------------------------------------ |
| `APP_NAME`    | `ws-cicd-bcc2026`        | application name                                 |
| `APP_ENV`     | `development`            | `development` enables SQL query logging          |
| `APP_PORT`    | `3000`                   | port the API listens on                          |
| `DB_HOST`     | `localhost`              | database host (`postgres` inside Docker Compose) |
| `DB_PORT`     | `5432`                   | database port                                    |
| `DB_USER`     | `postgres`               | database user                                    |
| `DB_PASSWORD` | `postgres`               | database password                                |
| `DB_NAME`     | `todo_db`                | database name                                    |
| `DB_SSLMODE`  | `disable`                | database ssl mode                                |
| `IMAGE_NAME`  | `ws-cicd-bcc2026`        | image name used by Compose for the app service   |
| `IMAGE_TAG`   | `latest`                 | image tag used by Compose for the app service    |

## Run Locally

The app still needs PostgreSQL, so start only the database with Compose and run the Go binary on your machine.

```bash
cp .env.example .env

docker compose up -d postgres

go mod download
go run ./cmd/migrate -m up
make run
```

The API is available at http://localhost:3000 (or whatever `APP_PORT` is set to).

> `make migrate-up` / `make migrate-down` run inside the app container. When you run the app on your machine, call the migration binary directly with `go run ./cmd/migrate -m up` (or `-m down`).

## Run with Docker Compose

This runs both PostgreSQL and the API in containers. `DB_HOST` is overridden to `postgres` inside the network, so no change to `.env` is required.

```bash
cp .env.example .env
make compose-up
```

Compose uses the image `ghcr.io/${IMAGE_NAME}:${IMAGE_TAG}` when it can be pulled, and falls back to building the local `Dockerfile` when it cannot.

Run the migration once the containers are up:

```bash
make migrate-up
```

Rebuild after changing the code, and stop everything when you are done:

```bash
make compose-restart
make compose-down
```

## Make Commands

| Command                | Description                      |
| ---------------------- | -------------------------------- |
| `make run`             | run the API locally              |
| `make build`           | build both binaries into `bin/`  |
| `make test`            | run the unit tests with coverage |
| `make lint`            | run golangci-lint                |
| `make lint-install`    | install golangci-lint            |
| `make tidy`            | tidy go modules                  |
| `make compose-up`      | `docker compose up -d`           |
| `make compose-down`    | `docker compose down`            |
| `make compose-restart` | `docker compose up -d --build`   |
| `make compose-logs`    | follow the container logs        |
| `make migrate-up`      | create the `todos` table (in the app container) |
| `make migrate-down`    | drop the `todos` table (in the app container)   |

## API Endpoints

| Method   | Path             | Description                 |
| -------- | ---------------- | --------------------------- |
| `GET`    | `/`              | welcome message             |
| `GET`    | `/hello`         | hello message with app info |
| `GET`    | `/health`        | app and database health     |
| `POST`   | `/api/todos`     | create a todo               |
| `GET`    | `/api/todos`     | list all todos              |
| `GET`    | `/api/todos/:id` | get a todo by id            |
| `PUT`    | `/api/todos/:id` | update a todo               |
| `DELETE` | `/api/todos/:id` | delete a todo               |

### Examples

```bash
curl http://localhost:3000/

curl -X POST http://localhost:3000/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Build a Fiber API"}'

curl http://localhost:3000/api/todos

curl -X PUT http://localhost:3000/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'

curl -X DELETE http://localhost:3000/api/todos/1
```

Response format:

```json
{
  "success": true,
  "message": "todo created",
  "data": {
    "id": 1,
    "title": "Learn Go",
    "description": "Build a Fiber API",
    "completed": false,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

## Testing and Linting

```bash
make test
make lint-install
make lint
```
