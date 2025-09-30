# Credit-Card-Api

### Context
Credit card REST api is developed with golang gin framework with mongo as database, it has exposed few endpoints, which will allow customer to create credit card transactions.
first you can create user(cardholder), customer can create one account with the same account customer can
make multiple transactions. 

### Prerequisites

Ensure the following tools are installed on your system before proceeding:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://go.dev/doc/install) (version **v1.24.0**)
- [Make](https://www.gnu.org/software/make/) utility

You can verify installation with:

```bash
docker --version
docker-compose --version
go version
make --version
```

### Run Application

#### Option 1: Run With Docker Env using Make
```
make clean build

make start

docker ps
```

#### Option 2:Run With docker-compose
- Build Local Application Binary

    ```go build .```

- Run App and Database container

    ```docker-compose up -d && docker ps```

#### Option 3:Run without docker
- Make sure MongoDB is running either in machine or via docker
- Export db url variable

  ```export MONGO_URI=mongodb://localhost:27017```

- Build Local Application Binary

  ```go build .```

- Start an application

  ```go run main.go```

Go to Below URL, Congratulations !!!!!

**URL:** http://localhost:8080/swagger/index.html

**Insomnia API client collection:** please find `Insomnia.json` at root level of project. 

#### Database

*To start / stop database*

```
make start_mongodb

make stop_mongodb
```


*Execute db queries from DB containers*

```docker-compose exec mongodb mongosh -u admin -p admin123 --authenticationDatabase admin```

CheatSheet for Mongo

```bash
  show dbs    //to list out database
  use credit-card-api   //to select database
  show collections  //to list out collections
  db.collection_name.find()   //to list all documents 
```
---

### Schema Definition

Users Collection

| Field Name      | Type     |
|-----------------|----------|
| `user_id`       | `string` |
| `name`          | `string` |
| `email`         | `string` |
| `mobile_number` | `string` |
| `created_at`    | `int64`  |
| `updated_at`    | `int64`  |

Accounts Collection

| Field Name        | Type     |
|-------------------|----------|
| `account_id`      | `string` |
| `user_id`         | `string` |
| `document_number` | `string` |
| `current_balance` | `float`  |
| `status`          | `string` |
| `created_at`      | `int64`  |
| `updated_at`      | `int64`  |

Transactions Collection

| Field Name       | Type     |
|------------------|----------|
| `transaction_id` | `string` |
| `account_id`     | `string` |
| `operation_type` | `string` |
| `amount`         | `float`  |
| `created_at`     | `int64`  |
| `updated_at`     | `int64`  |

---
