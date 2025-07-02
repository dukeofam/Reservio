# Reservio
A Go-based secure reservation backend for kindergarten signups. See API documentation and .env config for setup.

---

## API Reference

### Authentication
- All endpoints under `/api/parent`, `/api/admin`, and `/api/user` require authentication (session-based).
- Admin endpoints require the user to have the `admin` role.
- CSRF protection is enabled for all state-changing requests.

### Error Format
All errors are returned as:
```json
{"error": "Error message here"}
```

### Endpoints

#### Auth
- `POST /api/auth/register` — Register a new user
- `POST /api/auth/login` — Login
- `POST /api/auth/logout` — Logout
- `POST /api/auth/request-reset` — Request password reset (email)
- `POST /api/auth/reset-password` — Reset password with token

#### User
- `GET /api/user/profile` — Get current user info
- `PUT /api/user/profile` — Update user info

#### Parent
- `POST /api/parent/children` — Add child
- `GET /api/parent/children` — List children
- `PUT /api/parent/children/:id` — Edit child
- `DELETE /api/parent/children/:id` — Delete child
- `POST /api/parent/reserve` — Make reservation
- `GET /api/parent/reservations` — List own reservations
- `DELETE /api/parent/reservations/:id` — Cancel reservation

#### Admin
- `POST /api/admin/slots` — Create slot
- `PUT /api/admin/approve/:id` — Approve reservation
- `PUT /api/admin/reject/:id` — Reject reservation
- `GET /api/admin/reservations` — List/filter reservations
- `GET /api/admin/users` — List users
- `DELETE /api/admin/users/:id` — Delete user
- `PUT /api/admin/users/:id/role` — Change user role

#### Slots
- `GET /api/slots` — List all slots

#### Health/Version
- `GET /health` — Health check
- `GET /version` — API version info

---

### Example Request/Response

**Register:**
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "parent@example.com",
  "password": "strongpassword"
}
```

**Response:**
```json
{"message": "User registered", "user": "parent@example.com"}
```

**Error Example:**
```json
{"error": "Invalid email format"}
```

---

For more details, see the code or contact the backend team.