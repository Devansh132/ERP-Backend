# School ERP System - Backend API

Go backend API for the School ERP System with comprehensive Swagger documentation.

## 🚀 Quick Start

### 1. Prerequisites
- Go 1.21 or higher
- PostgreSQL 12+ (running and accessible)
- Git

### 2. Installation

```bash
# Clone repository
git clone <repository-url>
cd Backend

# Install dependencies
go mod download
```

### 3. Environment Setup

Create `.env` file in the Backend directory:

```env
# Server Configuration
SERVER_PORT=8080
ENVIRONMENT=development

# Database Configuration
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=school_erp

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production-min-32-characters
JWT_EXPIRY=24
```

### 4. Database Setup

**Option A: Automatic (Recommended)**
- The server will automatically create and migrate tables on startup
- Just ensure the database exists:
  ```sql
  CREATE DATABASE school_erp;
  ```

**Option B: Using pgAdmin**
1. Open pgAdmin
2. Connect to PostgreSQL server
3. Create database: `school_erp`
4. The server will auto-migrate tables on startup

### 5. Run the Server

```bash
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

### 6. Access Swagger Documentation

Open your browser:
```
http://localhost:8080/swagger/index.html
```

## 📁 Project Structure

```
Backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   │   ├── auth.go
│   │   └── student.go
│   ├── middleware/          # Middleware functions
│   │   └── auth.go
│   ├── models/              # Database models
│   │   ├── user.go
│   │   ├── student.go
│   │   └── ...
│   ├── repository/          # Data access layer
│   │   ├── user_repo.go
│   │   └── student_repo.go
│   └── services/            # Business logic
├── pkg/
│   ├── database/            # Database connection
│   ├── jwt/                 # JWT utilities
│   └── utils/               # Helper functions
├── config/                  # Configuration
├── migrations/              # SQL migration files
├── docs/                    # Swagger documentation
├── go.mod
└── README.md
```

## 🔌 API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - Register new user (Admin only)

### Admin Endpoints
- `GET /api/admin/students` - Get all students
- `GET /api/admin/students/:id` - Get student by ID
- `POST /api/admin/students` - Create student
- `PUT /api/admin/students/:id` - Update student
- `DELETE /api/admin/students/:id` - Delete student

*More endpoints to be implemented...*

## 🔐 Authentication

All protected endpoints require JWT authentication:

1. **Login to get token:**
   ```bash
   POST /api/auth/login
   {
     "email": "admin@school.com",
     "password": "password"
   }
   ```

2. **Use token in requests:**
   - In Swagger: Click "Authorize" → Enter `Bearer <token>`
   - In API calls: Add header `Authorization: Bearer <token>`

## 🗄️ Database Models

The system includes the following models:
- **User** - Authentication and user management
- **Student** - Student information
- **Teacher** - Teacher information
- **Class** - Class management
- **Section** - Section management
- **Attendance** - Attendance records
- **Subject** - Subject management
- **Exam** - Exam information
- **Mark** - Marks/grades
- **Assignment** - Assignment details
- **AssignmentSubmission** - Student submissions
- **Timetable** - Class schedules
- **Notice** - Notice board
- **CalendarEvent** - Calendar events
- **LeaveRequest** - Leave applications

## 🛠️ Development

### Generate Swagger Documentation

After adding new endpoints with Swagger annotations:

```bash
swag init -g cmd/server/main.go -o docs
```

### Build

```bash
go build -o bin/server.exe cmd/server/main.go
```

### Run Tests

```bash
go test ./...
```

## 📝 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | `8080` |
| `DB_DRIVER` | Database driver (postgres/mysql) | `postgres` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | (required) |
| `DB_NAME` | Database name | `school_erp` |
| `JWT_SECRET` | JWT secret key | (required) |
| `JWT_EXPIRY` | JWT expiry in hours | `24` |

## 🔄 Development Workflow

### 1. Setup (First Time)
```bash
# 1. Install dependencies
go mod download

# 2. Create .env file with database credentials

# 3. Create database in PostgreSQL
CREATE DATABASE school_erp;

# 4. Run server (auto-migrates tables)
go run cmd/server/main.go
```

### 2. Adding New Endpoints

1. **Create handler** in `internal/handlers/`
2. **Add Swagger annotations:**
   ```go
   // @Summary Endpoint summary
   // @Description Endpoint description
   // @Tags TagName
   // @Router /path [method]
   // @Security BearerAuth
   ```
3. **Register route** in `cmd/server/main.go`
4. **Regenerate Swagger:**
   ```bash
   swag init -g cmd/server/main.go -o docs
   ```
5. **Restart server**

### 3. Database Changes

- Models are in `internal/models/`
- GORM AutoMigrate runs on server startup
- For manual migrations, use SQL files in `migrations/`

## 🐛 Troubleshooting

### Database Connection Issues

**Error: "password authentication failed"**
- Verify password in `.env` matches PostgreSQL password
- Reset password if needed (see pgAdmin or use trust method)

**Error: "database does not exist"**
- Create database: `CREATE DATABASE school_erp;`

**Error: "connection refused"**
- Check PostgreSQL service is running
- Verify port 5432 is correct

### Swagger Issues

**"No operations defined in spec!"**
- Regenerate Swagger: `swag init -g cmd/server/main.go -o docs`
- Restart server
- Clear browser cache

**Double `/api` in URLs**
- Ensure `@Router` annotations don't include `/api` (basePath already has it)
- Example: `@Router /auth/login [post]` not `@Router /api/auth/login [post]`

### Server Issues

**Port already in use**
- Change `SERVER_PORT` in `.env`
- Or stop the process using port 8080

**Module not found**
- Run: `go mod download`
- Run: `go mod tidy`

## 📚 API Documentation

- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- All endpoints are documented with request/response schemas
- Use "Authorize" button to test protected endpoints

## 🔒 Security Notes

- Change `JWT_SECRET` in production (min 32 characters)
- Use HTTPS in production
- Implement rate limiting
- Regular security updates
- Database connection pooling enabled

## 📄 License

MIT License

## 🤝 Contributing

This is a private project. For questions or issues, contact the development team.

---

## 🎯 Next Steps

1. ✅ Database setup complete
2. ✅ Authentication working
3. ✅ Student CRUD endpoints
4. 🔄 Continue building other modules (Teachers, Classes, Attendance, etc.)
5. 🔄 Add more Swagger documentation
6. 🔄 Implement business logic in services
7. 🔄 Add validation and error handling
8. 🔄 Write unit tests
