# Reservio Backend

A Go Fiber backend for managing reservations, users, children, and admin operations.

## 🚀 Getting Started

### Requirements
- Go 1.20+
- PostgreSQL

### Setup
1. **Clone the repo:**
   ```sh
   git clone <repo-url>
   cd reservio
   ```
2. **Create a `.env` file:**
   See the example below. At minimum, set your `DATABASE_URL`.
3. **Install dependencies:**
   ```sh
   go mod tidy
   ```
4. **Setup the database:**
   ```sh
   ./setup_db.sh
   ```
5. **Run the server:**
   ```sh
   go run cmd/main.go
   ```
6. **Run tests:**
   ```sh
   ./run_tests.sh
   ```

### Example `.env` file
```
DATABASE_URL=postgres://reservio:reservio@localhost:5432/reservio?sslmode=disable
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

- **DB URI format:**
  `postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=disable`

## 🛠️ API Endpoints (Summary)

### Auth
- `POST /api/auth/register` — Register new user
- `POST /api/auth/login` — Login
- `POST /api/auth/logout` — Logout
- `POST /api/auth/request-reset` — Request password reset
- `POST /api/auth/reset-password` — Reset password

### User
- `GET /api/user/profile` — Get profile
- `PUT /api/user/profile` — Update profile

### Parent
- `POST /api/parent/children` — Add child
- `GET /api/parent/children` — List children
- `PUT /api/parent/children/:id` — Edit child
- `DELETE /api/parent/children/:id` — Delete child
- `POST /api/parent/reserve` — Make reservation
- `GET /api/parent/reservations` — List my reservations
- `DELETE /api/parent/reservations/:id` — Cancel reservation

### Admin (all require admin role)
- `POST /api/admin/slots` — Create slot
- `PUT /api/admin/approve/:id` — Approve reservation
- `PUT /api/admin/reject/:id` — Reject reservation
- `GET /api/admin/reservations` — List reservations (filter by status)
- `GET /api/admin/users` — List users
- `DELETE /api/admin/users/:id` — Delete user
- `PUT /api/admin/users/:id/role` — Update user role

### Public
- `GET /api/slots` — List slots
- `GET /health` — Health check
- `GET /version` — Version info

## 🧪 Testing
- Integration and unit tests are in `controllers/tests/` and `utils/`.
- Use `./run_tests.sh` to run all tests with a dedicated test DB.

## 📝 Notes
- All endpoints return JSON.
- CSRF protection and session security are enabled by default.
- For production, ensure all secrets are set via environment variables and HTTPS is enforced.

## 📖 API Documentation & OpenAPI (Swagger)

- The OpenAPI spec is in `docs/swagger.yaml` and covers all endpoints.
- View it locally:
  - Web: https://editor.swagger.io/ (import `docs/swagger.yaml`)
  - Docker: `docker run -p 8081:8080 -v $(pwd)/docs/swagger.yaml:/usr/share/nginx/html/swagger.yaml swaggerapi/swagger-ui`
- See `docs/README.md` for more details.

### Automated Swagger Generation

- This project is ready for [swaggo/swag](https://github.com/swaggo/swag) for Go doc-based OpenAPI generation.
- To enable, install swag:
  ```sh
  go install github.com/swaggo/swag/cmd/swag@latest
  swag init -g cmd/main.go -o docs/generated
  ```
- This will generate OpenAPI docs from Go comments in `docs/generated/swagger.yaml`.
- You can then view or merge this with the hand-written spec.

---

For more details, see the code and comments, or open an issue!