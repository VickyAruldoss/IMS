# IMS - Institution Management Service

Go microservice for managing institution member information.

## Tech Stack

- **Language:** Go 1.23
- **Router/Framework:** Gin (`github.com/gin-gonic/gin`)
- **Database:** PostgreSQL 16 (Docker) via `database/sql` + `lib/pq`
- **Migrations:** Liquibase (Docker image)
- **Docs:** Swagger via `swaggo/swag`
- **Testing:** `testify` with `testify/mock`

## Project Structure

```
IMS/
├── main.go                  # Entry point — loads config, opens DB, starts server
├── config/config.go         # Env-var config (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, SERVER_PORT)
├── model/member.go          # Member struct + CreateMemberRequest / UpdateMemberRequest DTOs
├── repository/              # MemberRepository interface + PostgreSQL and in-memory implementations
├── service/                 # MemberService interface + business logic
├── controller/              # Gin HTTP handlers (CRUD for members)
├── router/router.go         # Wires all layers, registers routes, serves Swagger UI
├── mocks/                   # testify/mock implementations of repository and service interfaces
├── docs/                    # Auto-generated Swagger files — do not edit manually
├── db/changelog/            # Liquibase changelogs
│   ├── db.changelog-master.xml
│   └── changes/001-create-members-table.xml
├── docker-compose.yml       # PostgreSQL + Liquibase services
└── Makefile                 # All dev commands
```

## Essential Commands

```bash
make db-start   # start PostgreSQL in Docker
make migrate    # run Liquibase migrations
make run        # build + start server on :8080
make test       # run unit tests with race detector
make build      # compile binary to ./ims
make swagger    # regenerate docs/ from annotations
make lint       # go vet ./...
make db-stop    # stop and remove containers
make clean      # remove compiled binary
```

## API

Base path: `/api/v1`

| Method   | Path            | Description       |
|----------|-----------------|-------------------|
| `POST`   | `/members`      | Create member     |
| `GET`    | `/members`      | List all members  |
| `GET`    | `/members/:id`  | Get member by ID  |
| `PUT`    | `/members/:id`  | Update member     |
| `DELETE` | `/members/:id`  | Delete member     |

Swagger UI: `http://localhost:8080/swagger/index.html`

## Environment Variables

| Variable      | Default       | Description          |
|---------------|---------------|----------------------|
| `DB_HOST`     | `localhost`   | PostgreSQL host      |
| `DB_PORT`     | `5432`        | PostgreSQL port      |
| `DB_USER`     | `ims_user`    | PostgreSQL user      |
| `DB_PASSWORD` | `ims_password`| PostgreSQL password  |
| `DB_NAME`     | `ims_db`      | PostgreSQL database  |
| `SERVER_PORT` | `8080`        | HTTP server port     |

## Architecture

Each layer depends only on the layer below it via an interface:

```
Controller → Service (interface) → Repository (interface) → PostgreSQL
```

- **controller** — HTTP request/response, binding, status codes
- **service** — business logic, UUID generation, timestamps
- **repository** — SQL queries; `NewPostgresRepository` for production, `NewInMemoryMemberRepository` for tests
- **mocks** — `MockMemberService` and `MockMemberRepository` used in unit tests

## Testing

Unit tests use mocks — no database required:

```bash
make test
```

- `service/member_service_test.go` — mocks the repository
- `controller/member_controller_test.go` — mocks the service

## Adding a Migration

Create a new file in `db/changelog/changes/` following the naming pattern:

```
002-<description>.xml
```

Register it in `db/changelog/db.changelog-master.xml`:

```xml
<include file="changes/002-<description>.xml" relativeToChangelogFile="true"/>
```

Then run `make migrate`.

## Swagger

After changing any `@` annotation in `controller/` or `main.go`, regenerate:

```bash
make swagger
```

## Commit Convention

```
Vicky | <short description>
```
