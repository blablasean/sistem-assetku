# CHANGELOG - Sistem AsetKu Fixes & Updates

## 🌟 Version 1.2.0 - Hotel Management & Operations Release

### 1. 📊 Export Laporan Bulanan ke Excel (.xlsx / .xls)
- Added **`📊 Export Excel`** feature generating formatted Excel tables complete with header titles, blue column headers, BOM encoding, numeric formatting, and total cost summary row.

### 2. 🖨️ Formatted A4 Landscape PDF Printing
- Added `@page { size: A4 landscape; margin: 8mm 10mm; }` enforcing auto-detection of A4 Landscape paper.
- Decoupled report modal from `.page-container` in DOM layout so top UI tables are hidden while the official report table renders with 100% full-width gridlines and zero truncation.

### 3. 🗑️ Full System Data Deletion Controls
- **Work Order Deletion**: Added `POST /workorders/delete` endpoint and `🗑️ Hapus` button for `hod`, `management`, and `admin`.
- **Asset Deletion**: Added `POST /assets/delete` endpoint and `🗑️` button for `hod`, `management`, and `admin`.
- **Preventive Maintenance Schedule Deletion**: Added `POST /maintenance/delete` endpoint and `🗑️ Hapus` button.
- Enforced GORM explicit query syntax (`c.db.Where("id = ?", id).Delete(...)`) and *Optimistic UI Updates* for instant row removal.

### 4. ✏️ Edit Preventive Maintenance Schedules
- Added `POST /maintenance/edit` endpoint and `✏️ Edit` modal allowing HOD, Management, and Admin to edit PM schedule frequency (Daily/Weekly/Monthly/Yearly), next due date (`next_run`), and inspection checklist items.

### 5. 👤 Custom UI Profile & Floating Toast Notifications
- Replaced all browser `alert()` popups with:
  - **User Profile Modal**: Displays avatar, name, concise role label, session badge, and permission rights.
  - **Floating Toast UI Notifications**: Clean toast notifications in the bottom-right corner for all action feedback.

### 6. 📋 Activity Log Overhaul & Status "Finish"
- **Unified Activity Log**: `GET /activitylogs` updated to return completed Work Orders (`Finish` / `Completed`) and completed Maintenance history (`MaintenanceHistory`) in a structured single payload.
- **ActivityLog UI**: Redesigned `ActivityLog.vue` with 3 summary metrics cards (Finished WO count, Maintenance count, Total Maintenance cost), two separate detailed tables, and real-time multi-column search filtering.
- **Status "Closed" -> "Finish"**: Updated Work Order completion status lifecycle across backend (`CloseWorkOrder`), controllers, routes, `WorkOrder.vue`, and `StatusBadge.vue` to use **`Finish`**.

### 7. 🔐 Session Security Enforcements
- Standardized `sessionStorage` usage for tokens and user details so opening a new browser tab/window enforces fresh re-authentication.

---

## 🌟 Version 1.1.0 - Hotel Operational Release

### 1. 👥 4 User Roles & Access Rights (RBAC Matrix)
- **Management Engineer (Supervisor)**: Full access to Dashboard, Asset Management, Work Order Management (Priority, History, Cancel, Close, Assign Worker), Activity Log.
- **Head of Department (HOD) Engineer**: Profile Asset (Registrasi & Edit Asset), Generate QR Code, Asset Management, Maintenance Management, Preventive Maintenance Checklist.
- **Staff Engineer (Teknisi Lapangan)**: Dashboard, Asset Management, Maintenance Management (Reminder, Report, History, PM Schedule), Task Execution & Progress Update.
- **External User (Staff Hotel / Departemen Lain)**: Work Order / Repair Request (Pengajuan tiket kerusakan dengan Lokasi & Priority), Work Order Status Tracking (Real-time tracking).

### 2. 🔧 Work Order Management Enhancements
- Added **Priority** (*Low, Medium, High, Emergency*) and **Location / Room** fields to Work Orders.
- Work Order lifecycle workflow: *Open ➔ In Progress ➔ Under Review ➔ Completed / Closed*.
- Full ticket controls: Assign Worker, Cancel Work Order, Close Work Order with completion date.

### 3. 📱 QR Code Scanning & Printing
- Printable QR Code badge generation for physical hotel assets.
- QR Code scanner simulation modal: Scan asset ➔ View specifications, last maintenance date, location, and PIC ➔ Direct **"Laporkan Kerusakan Aset Ini"** button.

### 4. 🔄 Asset Mutation & Relocation Tracking
- Official logging of physical asset transfers (e.g. moving TV, AC, or chairs between rooms/areas).
- Preserves location history, previous location, new location, new PIC, and mutation reason.

---

## 🔒 Security Fixes

### 1. Password Security
- ❌ **Before:** Password disimpan plain text di database
- ✅ **After:** Password di-hash menggunakan bcrypt (cost 10)
- **File:** `backend/utils/password.go` (baru)

### 2. Authentication & Authorization
- ❌ **Before:** Token menggunakan string concat sederhana
- ✅ **After:** JWT (JSON Web Tokens) dengan HS256 signing
- **Files:** 
  - `backend/utils/jwt.go` (baru)
  - `backend/controllers/auth_controller.go` (diperbarui)

### 3. Token Validation
- ❌ **Before:** Tidak ada validasi token di endpoint
- ✅ **After:** Middleware auth dengan JWT validation
- **File:** `backend/middlewares/auth.go` (baru)

### 4. CORS Headers
- ❌ **Before:** Frontend (5173) blocked oleh backend (8080)
- ✅ **After:** CORS middleware mengizinkan cross-origin requests
- **File:** `backend/middlewares/cors.go` (baru)

---

## 🏗️ Architecture Improvements

### 5. Routes Organization
- ❌ **Before:** Semua routes hardcoded di main.go
- ✅ **After:** Routes terorganisir dalam separate file
- **File:** `backend/routes/routes.go` (baru)

### 6. Response Standardization
- ❌ **Before:** Inconsistent responses (mix text & JSON)
- ✅ **After:** Unified response structure
- **File:** `backend/utils/response.go` (baru)
- **Format:** 
```json
{
  "status": true,
  "message": "Operation successful",
  "data": { ... },
  "error": null
}
```

### 7. Database Auto-Migration
- ❌ **Before:** Manual schema setup required
- ✅ **After:** Auto-migrate models on startup
- **File:** `backend/config/database.go` (diperbarui)

---

## 🎨 Frontend Improvements

### 8. API Interceptors
- ❌ **Before:** No automatic token handling
- ✅ **After:** Request interceptor adds token to Authorization header
- **File:** `frontend/src/api.js` (diperbarui)

### 9. Response Interceptor
- ❌ **Before:** No session management on 401
- ✅ **After:** Auto logout on token expiration
- **File:** `frontend/src/api.js` (diperbarui)

### 10. Better Error Handling
- ❌ **Before:** Generic error messages
- ✅ **After:** 
  - Specific error messages
  - Network error detection
  - Timeout handling
  - Loading state feedback
- **File:** `frontend/src/views/Login.vue` (diperbarui)

---

## 📦 Dependencies

### New Go Packages
```
github.com/golang-jwt/jwt/v5 v5.2.0
golang.org/x/crypto v0.31.0
```

### Updated Frontend
```
type: "module" (changed from "commonjs")
```

---

## 📄 New Files Created

### Backend
- `backend/utils/password.go` - Bcrypt password utilities
- `backend/utils/jwt.go` - JWT token generation & validation
- `backend/utils/response.go` - Standardized response helpers
- `backend/middlewares/cors.go` - CORS middleware
- `backend/middlewares/auth.go` - JWT auth middleware
- `backend/routes/routes.go` - Centralized route registration
- `backend/.env.example` - Environment template

### Frontend
- `frontend/.env.example` - Environment template
- `frontend/.env.development` - Development config

### Documentation
- `SETUP.md` - Setup & deployment guide
- `SECURITY.md` - Security implementation details
- `CHANGELOG.md` - This file

---

## 🔄 Modified Files

### Backend
- `backend/main.go` - Simplified, uses routes & middleware
- `backend/go.mod` - Added JWT & crypto dependencies
- `backend/config/database.go` - Added auto-migration
- `backend/controllers/auth_controller.go` - Now uses bcrypt & JWT

### Frontend
- `frontend/src/api.js` - Added interceptors
- `frontend/src/views/Login.vue` - Better error handling
- `frontend/package.json` - Minor cleanup

---

## 🚀 Quick Start

### Backend
```bash
cd backend
cp .env.example .env
# Edit .env with your database credentials
go mod tidy
go run main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

---

## ✅ Testing Checklist

- [ ] Backend starts without errors
- [ ] Frontend connects to backend
- [ ] User registration works
- [ ] User login returns JWT token
- [ ] Protected routes require token
- [ ] Token expires after 24 hours
- [ ] CORS requests work
- [ ] Error messages are helpful
- [ ] Password is hashed in database
- [ ] Logout clears localStorage

---

## 📝 Notes

1. **JWT Secret**: Ganti `JWT_SECRET` di .env untuk production
2. **CORS Origins**: Customize allowed origins jika perlu
3. **Token Expiration**: Saat ini 24 jam, bisa disesuaikan
4. **Database**: Auto-migration berjalan on startup

---

## 🔮 Future Enhancements

- [ ] Refresh token implementation
- [ ] Rate limiting
- [ ] Request logging
- [ ] Email verification
- [ ] Password reset flow
- [ ] Account lockout on failed login
- [ ] Role-based dashboard customization
- [ ] API documentation (Swagger/OpenAPI)
