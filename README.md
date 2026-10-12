# Employee API

A RESTful API built with Go, Echo, and PostgreSQL for employee registration, authentication, and attendance management with geotagging.

## Features

- **Employee Registration** — Register employees with personal, organizational, and account information.
- **Authentication** — Secure login using employee number and password.
- **Password Security** — Passwords are hashed using bcrypt.
- **JWT Authentication** — Protect endpoints using signed JSON Web Tokens.
- **Login Protection** — Track failed login attempts and temporarily lock accounts after repeated failures.
- **Employee Profile** — Retrieve the authenticated employee's identity.
- **Attendance with Geotagging** — Record check-in and check-out with latitude, longitude, GPS accuracy, and location capture timestamp.
- **Database Functions** — Use PostgreSQL functions for employee registration, login lookup, and attendance recording.
- **Database Constraints** — Enforce data integrity through foreign keys, unique constraints, and validation rules.

## Tech Stack

- Go
- Echo Framework
- PostgreSQL 16
- JWT (`github.com/golang-jwt/jwt/v5`)
- pgx PostgreSQL driver
- bcrypt (`golang.org/x/crypto`)
- Docker Compose
- Postman

## Prerequisites

Install the following before running the application:

- Go version compatible with `go.mod` (`1.27.2`)
- Docker and Docker Compose
- Git
- Postman (optional, for API testing)

## Project Structure

```text
employee-api/
├── cmd/
│   └── api/
│       └── main.go
├── db/
│   ├── schema.sql
│   ├── seed.sql
│   ├── functions.sql
│   └── migrate.sh
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── handler/
│   ├── middleware/
│   │   └── jwt_middleware.go
│   ├── repository/
│   ├── routes/
│   │   └── routes.go
│   └── service/
├── docker-compose.yml
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/mardenimihardi/employee-api.git
cd employee-api
```

### 2. Configure environment variables

Create a local environment file:

```bash
cp .env.example .env
```

Update `.env` with your local configuration:

```dotenv
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=lebah
DB_PASSWORD=your_local_database_password
DB_NAME=familymart
DB_SSLMODE=disable

JWT_SECRET=replace_with_a_secure_random_secret
JWT_EXPIRES_IN=15m
```

Use a strong, randomly generated JWT secret. Do not commit `.env` or real credentials to version control.

### 3. Start PostgreSQL

Start the PostgreSQL container:

```bash
docker compose up -d
```

Verify that the container is running:

```bash
docker compose ps
```

Ensure the database credentials in `.env` match the PostgreSQL configuration in `docker-compose.yml`.

### 4. Initialize the database

Make the migration script executable:

```bash
chmod +x db/migrate.sh
```

Run the database setup:

```bash
./db/migrate.sh
```

The script executes the following files in order:

1. `db/schema.sql` — Creates the database tables, indexes, and triggers.
2. `db/seed.sql` — Inserts initial department and position data.
3. `db/functions.sql` — Registers the PostgreSQL functions used by the application.

**Warning:** `db/schema.sql` drops and recreates existing application tables. Running the migration script against an existing database can delete its data. Use it only with a disposable development database unless you have backed up the data and intentionally want to reset the schema.

### 5. Install dependencies

```bash
go mod download
```

### 6. Run the API

```bash
go run ./cmd/api
```

The API will be available at:

```text
http://localhost:8080
```

Verify the service:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

## API Documentation

Base URL:

```text
http://localhost:8080
```

| Method | Endpoint | Authentication | Description |
|---|---|---|---|
| GET | `/` | No | API welcome message |
| GET | `/health` | No | Health check |
| POST | `/api/v1/auth/register` | No | Register an employee |
| POST | `/api/v1/auth/login` | No | Authenticate an employee |
| GET | `/api/v1/auth/me` | Bearer Token | Retrieve authenticated employee ID |
| POST | `/api/v1/attendance` | Bearer Token | Record attendance with geotagging |

### 1. Register Employee

**Request**

`POST /api/v1/auth/register`

Content-Type: `application/json`

```json
{
  "employee_no": "EMP001",
  "department_code": "IT",
  "position_code": "IT-DEV",
  "full_name": "Test Employee",
  "email": "test.employee@example.com",
  "phone": "081234567890",
  "date_of_birth": "1995-05-20",
  "gender": "male",
  "address": "Jl. Sudirman No. 10",
  "city": "Jakarta Selatan",
  "province": "DKI Jakarta",
  "postal_code": "12190",
  "password": "SecurePass123"
}
```

**Successful response — `201 Created`**

```json
{
  "message": "Employee registered successfully"
}
```

Notes:
- `employee_no`, `full_name`, `email`, and `password` are required.
- Passwords must be between 8 and 72 bytes.
- `gender` accepts `male` or `female`; it may be omitted.
- `department_code` and `position_code` must reference existing records.
- Duplicate employee numbers or email addresses return `409 Conflict`.

### 2. Login

**Request**

`POST /api/v1/auth/login`

```json
{
  "employee_no": "EMP001",
  "password": "SecurePass123"
}
```

**Successful response — `200 OK`**

```json
{
  "access_token": "<JWT_ACCESS_TOKEN>",
  "token_type": "Bearer",
  "expires_in": 900,
  "employee": {
    "id": "<EMPLOYEE_UUID>",
    "employee_no": "EMP001",
    "full_name": "Test Employee",
    "email": "test.employee@example.com",
    "role": "employee"
  }
}
```

The `expires_in` value is expressed in seconds. Its value depends on `JWT_EXPIRES_IN`.

The API tracks failed login attempts and temporarily locks an account after five failed attempts. A locked account can retry after the 15-minute lock period expires, subject to the application's account status and login rules.

| Status | Description |
|---|---|
| `200 OK` | Login successful |
| `401 Unauthorized` | Invalid employee number or password |
| `403 Forbidden` | Account is not active |
| `429 Too Many Requests` | Account is temporarily locked |

### 3. Get Authenticated Employee

**Request**

`GET /api/v1/auth/me`

Header:

```text
Authorization: Bearer <JWT_ACCESS_TOKEN>
```

**Successful response — `200 OK`**

```json
{
  "message": "Authenticated successfully",
  "employee_id": "<EMPLOYEE_UUID>"
}
```

The endpoint requires a valid JWT. The middleware also checks that the employee exists, has not been soft-deleted, and has an active account.

### 4. Record Attendance

**Request**

`POST /api/v1/attendance`

Header:

```text
Authorization: Bearer <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

Example request for check-in:

```json
{
  "attendance_type": "check_in",
  "latitude": -6.200000,
  "longitude": 106.816666,
  "accuracy_m": 10.5,
  "location_captured_at": "2026-10-12T18:18:00+07:00"
}
```

Example request for check-out:

```json
{
  "attendance_type": "check_out",
  "latitude": -6.200000,
  "longitude": 106.816666,
  "accuracy_m": 10.5,
  "location_captured_at": "2026-10-12T18:20:00+07:00"
}
```

Replace the example timestamps with the actual location capture time. The timestamp must be within the last 10 minutes and no more than 1 minute in the future.

**Successful response — `201 Created`**

The response contains a success message and the recorded attendance data:

```json
{
  "message": "Attendance recorded successfully",
  "data": {
    "id": "<ATTENDANCE_UUID>",
    "employee_id": "<EMPLOYEE_UUID>",
    "attendance_type": "check_in",
    "recorded_at": "<SERVER_TIMESTAMP>",
    "latitude": -6.2,
    "longitude": 106.816666,
    "accuracy_m": 10.5,
    "location_captured_at": "<LOCATION_CAPTURE_TIMESTAMP>"
  }
}
```

Notes:
- The employee identity is derived from the authenticated JWT, not from the request body.
- `attendance_type` accepts `check_in` or `check_out`.
- Latitude must be between `-90` and `90`.
- Longitude must be between `-180` and `180`.
- `accuracy_m` is optional and must be non-negative when provided.
- `recorded_at` is generated by the database; `location_captured_at` represents when the location was captured.
- The API validates the location timestamp, but does not independently verify GPS authenticity or physical presence.
- The current implementation does not enforce check-in/check-out ordering or prevent duplicate attendance records.

## Database Design

The application uses the following primary tables:

| Table | Purpose |
|---|---|
| `departments` | Stores department reference data |
| `positions` | Stores job position reference data |
| `employees` | Stores employee profiles, credentials, roles, and account status |
| `auth_sessions` | Stores session-related data for potential session management |
| `attendance_records` | Stores attendance records and geolocation information |

Database functions:

- `fn_register_employee` — Registers an employee and returns the new employee ID.
- `fn_login_employee` — Retrieves employee authentication data by employee number.
- `fn_record_attendance` — Records attendance and returns the inserted record.

The database schema also defines constraints, indexes, soft deletion for employees, and an update timestamp trigger.

## Testing with Postman

Import the **Employee API** collection into Postman.

[Open the Employee API Postman collection](https://go.postman.co/collection/364686-5c74860e-e746-4045-9d27-6319c6ae6ae9?source=collection_link)

Recommended testing sequence:

1. **Health** — Verify the API is running.
2. **Register** — Create an employee.
3. **Login** — Authenticate and automatically save `access_token` to the collection variable.
4. **Me** — Verify the authenticated identity.
5. **Attendance** — Submit a check-in record with a fresh location timestamp.
6. **Check-out** — Submit a check-out record with a fresh location timestamp.

The protected requests use the collection variable `{{access_token}}` as their Bearer Token.

## Security Considerations

- Store database credentials and JWT secrets in environment variables.
- Never commit `.env` or real credentials to the repository.
- Use HTTPS and an appropriate PostgreSQL SSL configuration in production.
- Use a strong JWT secret and configure token expiration appropriately.
- Avoid returning database internals or raw database errors to API clients in production.
- Apply database backups and safe, versioned migrations before deploying schema changes.