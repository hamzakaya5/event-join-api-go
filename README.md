# JOIN EVENT backend api

This is a simple backend API built with Go, PostgreSQL, and Redis, running in Docker containers. It supports horizontal scaling of the backend service and uses environment variables for configuration.

---

## Features

* REST API endpoints (example: `/register`, `/login`, `/join`)
* PostgreSQL database
* Redis caching
* Load-balanced backend using Nginx
* Easy to scale backend containers (docker compose)

---

## Project Structure

```
.
├──── main.go, go.mod, go.sum            # Go backend source code
├── models/ 
├── routes/                # Endpoint definitions
├── handlers/              # Functions for routers, endpoint handlers
├── nginx/                 # Nginx Dockerfile & config (load balancing/routing)
├── database/              # Initializes redis, postgresql connections
├── docker-compose.yaml
└── README.md
```

---

## Environment Variables

The backend reads configuration from environment variables:

| Variable      | Description              | Example    |
| ------------  | ------------------------ | ---------- |
| DB_HOST       | PostgreSQL host          | postgres   |
| DB_PORT       | PostgreSQL port          | 5432       |
| DB_USER       | PostgreSQL username      | myuser     |
| DB_PASSWORD   | PostgreSQL password      | mypassword |
| DB_NAME       | PostgreSQL database name | mydb       |
| REDIS_HOST    | Redis host               | redis      |
| REDIS_PORT    | Redis port               | 6379       |
| REDIS_PASSWORD| Redis password           | mypassword |

---

## Running with Docker Compose

1. Make sure Docker and Docker Compose are installed.
2. Build and start all services:

```bash
docker compose up --build
```

3. Optional: Run in detached mode:

```bash
docker compose up --build -d
```

4. Scale backend containers:

```bash
docker compose up --scale backend=3 -d
```

---

## Accessing the API

* Backend API: `http://localhost:8080`
* Nginx will load-balance requests to multiple backend containers.
* PostgreSQL and Redis are accessible to backend containers by their service names (`postgres`, `redis`). We cannot reach from outside of the network.



## Curl commands for API's

``
curl --location 'http://localhost:8080/register' \
--header 'email: <email>' \
--header 'password: <PASSWORD>' \
--header 'username: <USERNAME>'
``

``
curl --location 'http://localhost:8080/login' \
--header 'email: <EMAİL>' \
--header 'password: <PASSWORD>'
``


* User id, Autorization tokan and level informations are coming for the login response.
``
curl --location 'http://localhost:8080/join' \
--header 'Authorization: Bearer <TOKEN>' \
--header 'eventNo: <EVENT_ID>' \
--header 'userId: <USER_ID>' \
--header 'level: <LEVEL>'
``

---


## Stopping and Removing Containers

```bash
docker compose down
```

---

## Notes

* All services communicate over a private Docker network.
* You can adjust the number of backend replicas without changing the configuration.
* Make sure your `.env` or Docker environment variables are set correctly.
