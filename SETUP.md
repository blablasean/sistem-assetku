# Setup Instructions - Sistem AsetKu

## Backend Setup

### Prerequisites
- Go 1.25.4 or higher
- MySQL 5.7 or higher

### Installation

1. **Navigate to backend directory:**
```bash
cd backend
```

2. **Copy environment file:**
```bash
cp .env.example .env
```

3. **Edit .env with your configuration:**
```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=your_password
DB_NAME=db_sistemasetku
JWT_SECRET=your-secure-secret-key
SERVER_PORT=8080
```

4. **Create database:**
```bash
mysql -u root -p
CREATE DATABASE db_sistemasetku;
```

5. **Install dependencies:**
```bash
go mod download
go mod tidy
```

6. **Run backend:**
```bash
go run main.go
```

Backend will start on `http://localhost:8080`

---

## Frontend Setup

### Prerequisites
- Node.js 18+ and npm

### Installation

1. **Navigate to frontend directory:**
```bash
cd frontend
```

2. **Install dependencies:**
```bash
npm install
```

3. **Copy environment file:**
```bash
cp .env.example .env.local
```

4. **Edit .env.local (optional, defaults to localhost:8080):**
```bash
VITE_API_BASE=http://localhost:8080
```

5. **Run development server:**
```bash
npm run dev
```

Frontend will start on `http://localhost:5173`

---

## Deployment

### Docker Setup

**Backend Dockerfile exists at:** `backend/Dockerfile`

**Build and run:**
```bash
cd backend
docker build -t sistem-asetku-backend .
docker run -p 8080:8080 --env-file .env sistem-asetku-backend
```

**Frontend Dockerfile exists at:** `frontend/Dockerfile`

```bash
cd frontend
docker build -t sistem-asetku-frontend .
docker run -p 80:80 sistem-asetku-frontend
```

### Docker Compose

Run both services:
```bash
docker-compose up
```

---

## API Endpoints

### Public Endpoints
- `POST /auth/login` - User login
- `POST /auth/register` - User registration
- `GET /` - Health check

### Protected Endpoints (Require JWT Token)

**Assets:**
- `GET /assets?q=search` - Search assets
- `POST /assets` - Create asset (HOD only)
- `GET /assets/{id}` - Get asset details
- `GET /assets/{id}/history` - Get asset history

**Mutations:**
- `POST /mutations` - Create mutation (HOD only)
- `GET /mutations/{assetID}` - Get mutation history

**Work Orders:**
- `POST /workorders` - Create work order (HOD only)
- `POST /workorders/{id}/assign` - Assign worker (HOD only)
- `PUT /workorders/{id}/status` - Update status (HOD only)
- `GET /workorders/{id}/status` - Get status

**Maintenance:**
- `POST /maintenance/schedule` - Create schedule (HOD only)
- `POST /maintenance/{id}/checklist` - Submit checklist
- `GET /maintenance/{assetID}/history` - Get history

---

## Authentication

### Login Flow
1. POST `/auth/login` with username & password
2. Receive JWT token in response
3. Store token in localStorage
4. Include token in `Authorization: Bearer <token>` header for protected routes

### JWT Token
- Expires in 24 hours
- Contains: `user_id`, `username`, `role`
- Signed with HS256

---

## Troubleshooting

**Connection refused:**
- Check if backend is running on port 8080
- Verify MySQL is running

**CORS errors:**
- Backend includes CORS middleware
- Ensure frontend uses correct API base URL

**Database errors:**
- Run schema.sql manually if auto-migration fails
- Check database credentials in .env

**Login fails:**
- Verify user exists in database
- Check password is hashed (use bcrypt)
- Ensure JWT_SECRET is set
