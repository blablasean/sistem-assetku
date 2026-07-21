# Security Implementation - Sistem AsetKu

## Changes Made

### 1. Password Security ✅
**Problem:** Passwords stored as plain text
**Solution:** Implemented bcrypt hashing
- New file: `backend/utils/password.go`
- Uses bcrypt cost factor 10
- `HashPassword()` - hash on registration
- `VerifyPassword()` - verify on login

### 2. Authentication System ✅
**Problem:** Weak token using string concat
**Solution:** Implemented JWT (JSON Web Tokens)
- New file: `backend/utils/jwt.go`
- HS256 signing method
- 24-hour expiration
- Contains: user_id, username, role

### 3. Authorization Middleware ✅
**Problem:** Role checking scattered in code
**Solution:** Centralized middleware
- New file: `backend/middlewares/auth.go`
- Validates JWT from Authorization header
- Extracts claims to request context
- Returns proper error responses

### 4. CORS Support ✅
**Problem:** Frontend (port 5173) blocked by backend (port 8080)
**Solution:** Implemented CORS middleware
- New file: `backend/middlewares/cors.go`
- Allows cross-origin requests
- Handles preflight OPTIONS requests

### 5. Routes Organization ✅
**Problem:** All routes hardcoded in main.go
**Solution:** Organized routes
- New file: `backend/routes/routes.go`
- Centralized route registration
- Consistent error handling with response helper
- Proper HTTP method routing

### 6. Response Standardization ✅
**Problem:** Inconsistent response formats
**Solution:** Unified response structure
- New file: `backend/utils/response.go`
- `SendSuccess()` - successful responses
- `SendError()` - error responses
- Format: `{ status, message, data, error }`

### 7. Frontend Authentication ✅
**Problem:** No token handling in API calls
**Solution:** Added request/response interceptors
- Updated: `frontend/src/api.js`
- Automatically adds token to requests
- Handles 401 unauthorized (token expiration)
- Redirects to login on auth failure

### 8. Frontend Error Handling ✅
**Problem:** Generic error messages
**Solution:** Better error handling
- Updated: `frontend/src/views/Login.vue`
- Specific error messages
- Network error detection
- Timeout handling
- Loading state feedback

---

## Dependencies Added

### Backend (go.mod)
```
github.com/golang-jwt/jwt/v5 v5.2.0
golang.org/x/crypto v0.31.0
```

### Frontend (package.json)
No new dependencies - uses existing axios + vue-router

---

## Configuration Files

### Backend
- `.env.example` - Environment template
- New variables:
  - `JWT_SECRET` - JWT signing key
  - `SERVER_PORT` - Server port (default 8080)
  - `CORS_ALLOWED_ORIGINS` - (optional)

### Frontend
- `.env.example` - Environment template
- `.env.development` - Development config
- `VITE_API_BASE` - API base URL

---

## Database Auto-Migration

Updated `backend/config/database.go`:
- Auto-migrates all models on startup
- Creates tables if they don't exist
- Maintains schema consistency
- No need to run schema.sql manually

Models auto-migrated:
- User
- Asset
- SparePart
- Mutation
- WorkOrder
- PreventiveMaintenance
- MaintenanceHistory
- ActivityLog

---

## Testing the Security

### 1. Register New User
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123",
    "name": "Test User",
    "role": "engineer"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }'
```
Response:
```json
{
  "status": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 3. Use Token in Protected Route
```bash
curl -X GET http://localhost:8080/assets?q=test \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

---

## Security Best Practices Implemented

✅ Password Hashing with bcrypt (cost 10)
✅ JWT with HS256 signature
✅ Token expiration (24 hours)
✅ Authorization header validation
✅ Role-based access control (RBAC)
✅ CORS headers properly set
✅ Consistent error responses (no info leakage)
✅ Middleware-based auth (DRY principle)

---

## Future Improvements

🔄 Refresh token implementation
🔄 Rate limiting
🔄 Request logging & audit trail
🔄 HTTPS/TLS in production
🔄 Environment-based JWT secret rotation
🔄 Password strength validation
🔄 Account lockout on failed attempts
🔄 Email verification for registration
