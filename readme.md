# Users REST API

This is a REST API for managing users and their data, built using the Echo framework. The application supports user registration, authentication, and admin management of user accounts.

## Features

- User registration and authentication
- Admin functionality for user management (create, update, delete, and retrieve users)
- Secure user data handling
- Swagger API documentation for easy exploration of endpoints (http://host:port/swagger/*)

docker-compose build:

```
docker-compose up --build
```

Initial Admin User
Upon starting the application, an admin user is automatically initialized with the following credentials:

```
Username: admin
Password: admin
```

DB diagram: https://dbdiagram.io/d/users-rest-api-66fa9f64f9b1444815e7863a

DB migration: https://github.com/golang-migrate/migrate

.env:

```
DB_PASSWORD=
DB_HOST=users-rest-api-db-1 (or localhost if u dont use docker)
DB_PORT=
DB_USERNAME=
DB_DBNAME=
DB_SSLMODE=disable
PASSWORD_HASH_SALT=
TOKEN_SECRET_KEY=
```
