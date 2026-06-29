# IMS — Institution Management Service

A Go microservice for managing institution member information (the institution admin's
member records). Built with a clean, layered architecture and backed by PostgreSQL.

## Tech Stack

- **Language:** Go 1.23
- **Web framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL 16 (via `database/sql` + `lib/pq`)
- **Migrations:** [Liquibase](https://www.liquibase.org/) (Docker image)
- **API docs:** [Swagger](https://github.com/swaggo/swag)
- **Testing:** [testify](https://github.com/stretchr/testify) with `testify/mock`

## Architecture

Each layer depends only on the layer below it through an interface:

```
Controller → Service (interface) → Repository (interface) → PostgreSQL
```

- **controller** — HTTP request/response handling, binding, status codes
- **service** — business logic, UUID generation, timestamps
- **repository** — SQL queries (`NewPostgresRepository` for production,
  `NewInMemoryMemberRepository` for tests)
- **mocks** — `MockMemberService` and `MockMemberRepository` used in unit tests

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose

### Run

```bash
make db-start   # start PostgreSQL in Docker
make migrate    # run Liquibase migrations
make run        # build + start the server
```

The server listens on `http://localhost:8080` and Swagger UI is available at
`http://localhost:8080/swagger/index.html`.

### Database UI (pgAdmin4)

Optionally start pgAdmin4 to browse the database:

```bash
make pgadmin   # http://localhost:5050  (login: admin@ims.com / admin)
```

When adding a server connection inside pgAdmin, use host `postgres` (the Compose
service name), port `5432`, user `ims_user`, password `ims_password`, database `ims_db`.
Run `make pgadmin-stop` to tear it down.

## API

Base path: `/api/v1`

| Method   | Path           | Description      |
|----------|----------------|------------------|
| `POST`   | `/members`     | Create member    |
| `GET`    | `/members`     | List all members |
| `GET`    | `/members/:id` | Get member by ID |
| `PUT`    | `/members/:id` | Update member    |
| `DELETE` | `/members/:id` | Delete member    |

### Member model

```json
{
  "id": "uuid",
  "name": "string",
  "email": "string",
  "role": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

`name`, `email`, and `role` are required on create; `email` must be a valid address.

## Configuration

Configured via environment variables (defaults shown). Copy `.env.example` to
`.env` to override them — the `Makefile` and Docker Compose load `.env` automatically:

| Variable           | Default         | Description           |
|--------------------|-----------------|-----------------------|
| `DB_HOST`          | `localhost`     | PostgreSQL host       |
| `DB_PORT`          | `5432`          | PostgreSQL port       |
| `DB_USER`          | `ims_user`      | PostgreSQL user       |
| `DB_PASSWORD`      | `ims_password`  | PostgreSQL password   |
| `DB_NAME`          | `ims_db`        | PostgreSQL database   |
| `SERVER_PORT`      | `8080`          | HTTP server port      |
| `PGADMIN_EMAIL`    | `admin@ims.com` | pgAdmin4 login email  |
| `PGADMIN_PASSWORD` | `admin`         | pgAdmin4 login password |
| `PGADMIN_PORT`     | `5050`          | pgAdmin4 host port    |

## Make Commands

```bash
make db-start      # start PostgreSQL in Docker
make migrate       # run Liquibase migrations
make run           # build + start server
make build         # compile binary to ./ims
make test          # run unit tests with the race detector
make swagger       # regenerate docs/ from annotations
make lint          # go vet ./...
make pgadmin       # start pgAdmin4 (http://localhost:5050)
make pgadmin-stop  # stop and remove the pgAdmin4 container
make db-stop       # stop and remove containers
make clean         # remove compiled binary
```

## Testing

Unit tests use mocks — no database required:

```bash
make test
```
