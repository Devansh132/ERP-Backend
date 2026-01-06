# Backend Architecture & Flow - Detailed Explanation

## 📐 Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT (Frontend)                        │
│                    Angular / Postman / Swagger                   │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTP Request
                             │ (JSON)
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    GIN WEB SERVER (Port 8080)                   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  CORS Middleware                                         │  │
│  │  - Handles cross-origin requests                        │  │
│  │  - Sets headers for frontend communication              │  │
│  └──────────────────────────────────────────────────────────┘  │
│                             │                                    │
│                             ▼                                    │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Router (Gin Router)                                    │  │
│  │  - Routes: /api/auth, /api/admin, /api/teacher, etc.   │  │
│  └──────────────────────────────────────────────────────────┘  │
│                             │                                    │
│                             ▼                                    │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Middleware Chain                                         │  │
│  │  1. AuthMiddleware - Validates JWT token                 │  │
│  │  2. RoleMiddleware - Checks user role permissions        │  │
│  └──────────────────────────────────────────────────────────┘  │
│                             │                                    │
│                             ▼                                    │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Handlers (Business Logic)                                │  │
│  │  - AuthHandler, StudentHandler, TeacherHandler, etc.      │  │
│  └──────────────────────────────────────────────────────────┘  │
│                             │                                    │
│                             ▼                                    │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Repository Layer (Data Access)                          │  │
│  │  - UserRepository, StudentRepository, etc.               │  │
│  └──────────────────────────────────────────────────────────┘  │
│                             │                                    │
│                             ▼                                    │
└─────────────────────────────┼────────────────────────────────────┘
                              │
                              │ GORM Queries
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    PostgreSQL Database                          │
│  - Users, Students, Teachers, Classes, Sections, Subjects, etc. │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Server Startup Flow

```
START
  │
  ├─► 1. Load Configuration (config.LoadConfig())
  │     │
  │     ├─► Read .env file
  │     ├─► Set defaults if missing
  │     └─► Store in AppConfig struct
  │
  ├─► 2. Connect to Database (database.Connect())
  │     │
  │     ├─► Build DSN (Data Source Name)
  │     ├─► Connect using GORM
  │     └─► Store connection in database.DB
  │
  ├─► 3. Auto Migrate Models (database.DB.AutoMigrate())
  │     │
  │     ├─► Create/Update tables based on models
  │     ├─► User, Student, Teacher, Class, Section, etc.
  │     └─► Handles schema changes automatically
  │
  ├─► 4. Initialize Handlers
  │     │
  │     ├─► authHandler = NewAuthHandler()
  │     ├─► studentHandler = NewStudentHandler()
  │     ├─► teacherHandler = NewTeacherHandler()
  │     ├─► classHandler = NewClassHandler()
  │     ├─► sectionHandler = NewSectionHandler()
  │     └─► subjectHandler = NewSubjectHandler()
  │
  ├─► 5. Setup Router (gin.Default())
  │     │
  │     ├─► Create Gin router instance
  │     ├─► Add CORS middleware
  │     └─► Add health check route
  │
  ├─► 6. Register API Routes
  │     │
  │     ├─► /api/auth (public)
  │     ├─► /api/admin (protected - admin only)
  │     ├─► /api/teacher (protected - teacher only)
  │     └─► /api/student (protected - student only)
  │
  ├─► 7. Setup Swagger Documentation
  │     │
  │     └─► /swagger/*any → Swagger UI
  │
  └─► 8. Start Server (router.Run(":8080"))
        │
        └─► Server listening on port 8080
            └─► READY TO ACCEPT REQUESTS
```

---

## 🔐 Authentication Flow (Login)

```
CLIENT REQUEST: POST /api/auth/login
{
  "email": "admin@school.com",
  "password": "password123"
}
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 1. CORS Middleware                                      │
│    - Adds CORS headers                                  │
│    - Allows request to proceed                          │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 2. Router Matches Route                                 │
│    - Route: /api/auth/login                             │
│    - Handler: authHandler.Login                         │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 3. AuthHandler.Login()                                  │
│    │                                                     │
│    ├─► Parse JSON body → LoginRequest                   │
│    │   - email, password                                │
│    │                                                     │
│    ├─► Validate request (binding)                       │
│    │   - Check required fields                          │
│    │                                                     │
│    ├─► Find user by email                               │
│    │   userRepo.FindByEmail(email)                      │
│    │   └─► Query: SELECT * FROM users WHERE email = ?   │
│    │                                                     │
│    ├─► Check password hash                              │
│    │   utils.CheckPasswordHash(password, hash)          │
│    │   └─► Compare using bcrypt                         │
│    │                                                     │
│    ├─► Check user status                                │
│    │   - Must be "active"                               │
│    │                                                     │
│    ├─► Generate JWT token                               │
│    │   jwt.GenerateToken(userID, email, role)           │
│    │   └─► Create token with claims                     │
│    │       - user_id, email, role, exp, iat            │
│    │                                                     │
│    └─► Return response                                  │
│        {                                                │
│          "token": "eyJhbGci...",                        │
│          "user": { id, email, role }                    │
│        }                                                │
└─────────────────────────────────────────────────────────┘
  │
  ▼
CLIENT RECEIVES: 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "admin@school.com",
    "role": "admin"
  }
}
```

---

## 🛡️ Protected Route Flow (Example: GET /api/admin/students)

```
CLIENT REQUEST: GET /api/admin/students
Headers:
  Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 1. CORS Middleware                                      │
│    - Adds CORS headers                                  │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 2. Router Matches Route                                 │
│    - Route: /api/admin/students                         │
│    - Group: admin (has middleware chain)                │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 3. AuthMiddleware (First in chain)                      │
│    │                                                     │
│    ├─► Extract Authorization header                     │
│    │   - Try: Authorization, authorization, AUTHORIZATION│
│    │                                                     │
│    ├─► Parse token                                      │
│    │   - Remove "Bearer " prefix                        │
│    │   - Trim whitespace                                │
│    │                                                     │
│    ├─► Validate token                                   │
│    │   jwt.ValidateToken(tokenString)                   │
│    │   ├─► Parse JWT                                    │
│    │   ├─► Check signature                              │
│    │   ├─► Check expiration                             │
│    │   └─► Extract claims                               │
│    │                                                     │
│    ├─► Set context variables                            │
│    │   c.Set("user_id", claims.UserID)                  │
│    │   c.Set("user_email", claims.Email)                │
│    │   c.Set("user_role", claims.Role)                  │
│    │                                                     │
│    └─► c.Next() → Continue to next middleware           │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 4. RoleMiddleware (Second in chain)                      │
│    │                                                     │
│    ├─► Get user role from context                       │
│    │   c.Get("user_role")                                │
│    │                                                     │
│    ├─► Check if role matches                            │
│    │   - Allowed roles: ["admin"]                       │
│    │   - User role: "admin"                             │
│    │   - Match? → Continue                              │
│    │   - No match? → 403 Forbidden                     │
│    │                                                     │
│    └─► c.Next() → Continue to handler                   │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 5. StudentHandler.GetStudents()                          │
│    │                                                     │
│    ├─► Get query parameters                            │
│    │   - class_id, section_id, page, limit              │
│    │                                                     │
│    ├─► Build query                                      │
│    │   query := DB.Preload("User")                      │
│    │                   .Preload("Class")               │
│    │                   .Preload("Section")             │
│    │                                                     │
│    ├─► Apply filters                                   │
│    │   if classID != "" {                              │
│    │     query = query.Where("class_id = ?", classID)  │
│    │   }                                                │
│    │                                                     │
│    ├─► Execute query                                   │
│    │   query.Find(&students)                           │
│    │   └─► SQL: SELECT * FROM students                 │
│    │              LEFT JOIN users ...                   │
│    │              LEFT JOIN classes ...                 │
│    │              WHERE class_id = ?                    │
│    │                                                     │
│    └─► Return JSON response                            │
│        c.JSON(200, students)                            │
└─────────────────────────────────────────────────────────┘
  │
  ▼
CLIENT RECEIVES: 200 OK
[
  {
    "id": 1,
    "admission_number": "STU001",
    "first_name": "John",
    "last_name": "Doe",
    "user": { "id": 1, "email": "john@school.com" },
    "class": { "id": 1, "name": "1st" },
    "section": { "id": 1, "name": "A" }
  },
  ...
]
```

---

## ➕ Create Operation Flow (Example: POST /api/admin/students)

```
CLIENT REQUEST: POST /api/admin/students
Headers:
  Authorization: Bearer eyJhbGci...
Body:
{
  "user_id": 2,
  "admission_number": "STU002",
  "first_name": "Jane",
  "last_name": "Smith",
  "date_of_birth": "2010-05-15",
  "gender": "female",
  "class_id": 1,
  "section_id": 1
}
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 1-2. CORS & Router (same as above)                     │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 3-4. AuthMiddleware & RoleMiddleware (same as above)   │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 5. StudentHandler.CreateStudent()                       │
│    │                                                     │
│    ├─► Parse JSON body                                  │
│    │   c.ShouldBindJSON(&req)                           │
│    │   └─► Validates required fields                     │
│    │                                                     │
│    ├─► Parse date                                       │
│    │   time.Parse("2006-01-02", req.DateOfBirth)       │
│    │   └─► Convert string to time.Time                  │
│    │                                                     │
│    ├─► Create Student model                             │
│    │   student := &models.Student{                      │
│    │     UserID: req.UserID,                            │
│    │     AdmissionNumber: req.AdmissionNumber,           │
│    │     FirstName: req.FirstName,                      │
│    │     ...                                            │
│    │     Status: "active"                                │
│    │   }                                                │
│    │                                                     │
│    ├─► Save to database                                 │
│    │   studentRepo.Create(student)                     │
│    │   └─► DB.Create(student)                           │
│    │       └─► SQL: INSERT INTO students (...)            │
│    │           VALUES (...)                              │
│    │                                                     │
│    └─► Return created student                            │
│        c.JSON(201, student)                            │
└─────────────────────────────────────────────────────────┘
  │
  ▼
CLIENT RECEIVES: 201 Created
{
  "id": 2,
  "admission_number": "STU002",
  "first_name": "Jane",
  "last_name": "Smith",
  "date_of_birth": "2010-05-15T00:00:00Z",
  "status": "active",
  "created_at": "2025-12-25T10:00:00Z"
}
```

---

## 🔄 Update Operation Flow (Example: PUT /api/admin/students/:id)

```
CLIENT REQUEST: PUT /api/admin/students/2
Headers:
  Authorization: Bearer eyJhbGci...
Body:
{
  "first_name": "Jane Updated",
  "phone": "1234567890"
}
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 1-4. Middleware chain (same as above)                  │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 5. StudentHandler.UpdateStudent()                       │
│    │                                                     │
│    ├─► Parse ID from URL                               │
│    │   strconv.ParseUint(c.Param("id"), 10, 32)        │
│    │   └─► id = 2                                       │
│    │                                                     │
│    ├─► Parse JSON body                                 │
│    │   c.ShouldBindJSON(&req)                           │
│    │                                                     │
│    ├─► Find existing student                            │
│    │   studentRepo.FindByID(2)                          │
│    │   └─► SELECT * FROM students WHERE id = 2         │
│    │   └─► If not found → 404                           │
│    │                                                     │
│    ├─► Update fields (only provided ones)              │
│    │   if req.FirstName != "" {                         │
│    │     student.FirstName = req.FirstName             │
│    │   }                                                │
│    │   if req.Phone != "" {                             │
│    │     student.Phone = req.Phone                     │
│    │   }                                                │
│    │                                                     │
│    ├─► Save changes                                    │
│    │   studentRepo.Update(student)                     │
│    │   └─► DB.Save(student)                             │
│    │       └─► SQL: UPDATE students SET ...            │
│    │                 WHERE id = 2                       │
│    │                                                     │
│    └─► Return updated student                           │
│        c.JSON(200, student)                             │
└─────────────────────────────────────────────────────────┘
  │
  ▼
CLIENT RECEIVES: 200 OK
{
  "id": 2,
  "first_name": "Jane Updated",
  "phone": "1234567890",
  "updated_at": "2025-12-25T10:05:00Z"
}
```

---

## 🗑️ Delete Operation Flow (Example: DELETE /api/admin/students/:id)

```
CLIENT REQUEST: DELETE /api/admin/students/2
Headers:
  Authorization: Bearer eyJhbGci...
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 1-4. Middleware chain (same as above)                  │
└─────────────────────────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 5. StudentHandler.DeleteStudent()                       │
│    │                                                     │
│    ├─► Parse ID from URL                               │
│    │   id = 2                                           │
│    │                                                     │
│    ├─► Delete from database                             │
│    │   studentRepo.Delete(2)                           │
│    │   └─► DB.Delete(&models.Student{}, 2)             │
│    │       └─► SQL: UPDATE students                      │
│    │                 SET deleted_at = NOW()              │
│    │                 WHERE id = 2                        │
│    │       (Soft delete - GORM default)                │
│    │                                                     │
│    └─► Return success message                           │
│        c.JSON(200, {"message": "Student deleted..."})  │
└─────────────────────────────────────────────────────────┘
  │
  ▼
CLIENT RECEIVES: 200 OK
{
  "message": "Student deleted successfully"
}
```

---

## 📦 Component Layers Explained

### 1. **Config Layer** (`config/config.go`)
- **Purpose**: Load environment variables and configuration
- **Flow**: 
  - Reads `.env` file
  - Sets defaults if missing
  - Stores in `AppConfig` struct
- **Used by**: Database connection, JWT secret, server port

### 2. **Database Layer** (`pkg/database/database.go`)
- **Purpose**: Database connection and management
- **Flow**:
  - Builds DSN from config
  - Connects using GORM
  - Exposes `DB` variable globally
- **Features**: Supports PostgreSQL and MySQL

### 3. **Model Layer** (`internal/models/`)
- **Purpose**: Define data structures
- **Features**:
  - GORM tags for database mapping
  - JSON tags for API responses
  - Relationships (foreign keys)
  - Soft deletes

### 4. **Repository Layer** (`internal/repository/`)
- **Purpose**: Data access abstraction
- **Methods**: Create, Read, Update, Delete, FindBy*
- **Benefits**: 
  - Separates business logic from database
  - Easy to test
  - Reusable queries

### 5. **Handler Layer** (`internal/handlers/`)
- **Purpose**: HTTP request handling
- **Responsibilities**:
  - Parse request body/params
  - Validate input
  - Call repository methods
  - Return JSON responses
  - Handle errors

### 6. **Middleware Layer** (`internal/middleware/`)
- **AuthMiddleware**: Validates JWT tokens
- **RoleMiddleware**: Checks user permissions
- **CORS Middleware**: Handles cross-origin requests

### 7. **JWT Layer** (`pkg/jwt/jwt.go`)
- **Purpose**: Token generation and validation
- **GenerateToken**: Creates JWT with claims
- **ValidateToken**: Verifies and extracts claims

### 8. **Utils Layer** (`pkg/utils/`)
- **Password hashing**: bcrypt for secure password storage
- **Password checking**: Compare hashed passwords

---

## 🔄 Complete Request Lifecycle

```
1. CLIENT SENDS REQUEST
   └─► HTTP Method + URL + Headers + Body

2. SERVER RECEIVES
   └─► Gin router matches route

3. MIDDLEWARE CHAIN
   ├─► CORS Middleware
   ├─► AuthMiddleware (if protected)
   └─► RoleMiddleware (if role-based)

4. HANDLER EXECUTION
   ├─► Parse request
   ├─► Validate input
   ├─► Call repository
   └─► Return response

5. REPOSITORY LAYER
   └─► Execute GORM queries

6. DATABASE
   └─► Execute SQL queries

7. RESPONSE FLOW (reverse)
   Database → Repository → Handler → Middleware → Client

8. CLIENT RECEIVES
   └─► JSON response with status code
```

---

## 🎯 Key Design Patterns

### 1. **Layered Architecture**
- Clear separation of concerns
- Each layer has specific responsibility
- Easy to maintain and test

### 2. **Repository Pattern**
- Abstracts database operations
- Handlers don't know about SQL
- Easy to swap database implementations

### 3. **Middleware Pattern**
- Reusable authentication/authorization
- Applied to route groups
- Executes before handlers

### 4. **Dependency Injection**
- Handlers receive repositories
- Easy to mock for testing
- Loose coupling

---

## 🔒 Security Flow

```
1. User logs in → Receives JWT token
2. Token stored in client (localStorage/cookie)
3. Every protected request includes:
   Authorization: Bearer <token>
4. Server validates token:
   ├─► Signature valid?
   ├─► Not expired?
   └─► Extract user info
5. Check role permissions
6. Allow/deny request
```

---

## 📊 Database Relationships

```
Users (1) ──< (1) Students
Users (1) ──< (1) Teachers
Classes (1) ──< (N) Sections
Classes (1) ──< (N) Students
Sections (1) ──< (N) Students
Subjects (1) ──< (N) Exams
Students (1) ──< (N) Marks
```

---

## 🚨 Error Handling Flow

```
Error Occurs
  │
  ├─► Database Error?
  │   └─► Return 500 Internal Server Error
  │
  ├─► Validation Error?
  │   └─► Return 400 Bad Request
  │
  ├─► Not Found?
  │   └─► Return 404 Not Found
  │
  ├─► Unauthorized?
  │   └─► Return 401 Unauthorized
  │
  └─► Forbidden?
      └─► Return 403 Forbidden
```

---

## 📝 Summary

**Request Flow:**
1. Client → HTTP Request
2. CORS Middleware → Add headers
3. Router → Match route
4. AuthMiddleware → Validate token
5. RoleMiddleware → Check permissions
6. Handler → Process request
7. Repository → Database query
8. Database → Execute SQL
9. Response → JSON back to client

**Key Points:**
- ✅ Clean separation of concerns
- ✅ Reusable middleware
- ✅ Type-safe with Go
- ✅ Automatic database migrations
- ✅ JWT-based authentication
- ✅ Role-based access control
- ✅ Swagger documentation

This architecture ensures:
- **Security**: JWT tokens + role-based access
- **Scalability**: Layered design allows easy scaling
- **Maintainability**: Clear structure, easy to modify
- **Testability**: Each layer can be tested independently



