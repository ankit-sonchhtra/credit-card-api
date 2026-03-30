# Credit-Card-Api

### Context

Credit card REST api is developed with golang gin framework with postgres as database, it has exposed few endpoints,
which will allow customer to create credit card transactions.
customer can create one account with the same account customer can
make transactions.

### Prerequisites

Ensure the following tools are installed on your system before proceeding:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://go.dev/doc/install) (version **v1.25.0**)
- [Make](https://www.gnu.org/software/make/) utility

You can verify installation with:

```bash
docker --version
docker-compose --version
go version
make --version
```

### Run Application

#### Option 1: Run With Docker Env using Make (Recommended)

```
make clean build

make start

docker ps
```

#### Option 2:Run With docker-compose

- Build Local Application Binary

  ```CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o credit-card-api .```

- Run App and Database container

  ```docker-compose up -d && docker ps```

#### Option 3:Run without docker

- Make sure Postgres is running either in machine or via docker
- Export db url variable

  ```export DB_URL=postgresql://user:password@localhost:5432/credit_card_api?sslmode=disable```

- Start an application

  ```go run main.go```

#### Swagger

**URL:** http://localhost:8080/swagger/index.html

**Insomnia API client collection:** please find `Insomnia.json` at root level of project.

#### Database

*To start / stop database*

```
make start_postgres

make stop_postgres
```

---

### Schema Definition

accounts

| Field Name        | Type        |
|-------------------|-------------|
| `account_id`      | `BIGINT`    |
| `document_number` | `VARCHAR`   |
| `created_at`      | `TIMESTAMP` |

transactions

| Field Name       | Type        |
|------------------|-------------|
| `transaction_id` | `BIGINT`    |
| `account_id`     | `BIGINT`    |
| `operation_type` | `BIGINT`    |
| `amount`         | `NUMERIC`   |
| `created_at`     | `TIMESTAMP` |

operation_types

| Field Name          | Type      |
|---------------------|-----------|
| `operation_type_id` | `INT`     |
| `document_number`   | `VARCHAR` |

---
