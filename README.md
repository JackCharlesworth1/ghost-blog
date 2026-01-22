# Ghost Blog - Anonymous Blogging Platform

An anonymous blogging platform where users can submit posts and browse them in a unique snap-scroll interface, similar to Instagram stories. Posts are displayed randomly with country detection from IP addresses.

## Features

- **Anonymous Posting**: Submit blog posts without registration
- **Snap-Scroll Browsing**: Browse posts one at a time with smooth scroll transitions
- **Country Detection**: Automatically detect and display the country of origin
- **Rate Limiting**: IP-based rate limiting (5 posts per hour)
- **Admin Panel**: Password-protected admin interface for post moderation
- **Dark Mode**: Clean dark grey theme throughout

## Tech Stack

### Backend
- Go 1.21+
- Chi Router
- PostgreSQL 14+
- pgx (PostgreSQL driver)

### Frontend
- Vue.js 3 (Composition API)
- Vue Router
- Pinia (State Management)
- Axios
- Vite

## Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- PostgreSQL 14+
- Docker and Docker Compose (optional, for database)

## Quick Start

### 1. Database Setup

#### Option A: Using Docker (Recommended)
```bash
docker-compose up -d
```

#### Option B: Manual PostgreSQL Setup
1. Install PostgreSQL
2. Create a database:
```sql
CREATE DATABASE ghostblog;
```

### 2. Backend Setup

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env with your database credentials
# DATABASE_URL=postgres://postgres:postgres@localhost:5432/ghostblog?sslmode=disable
# PORT=8080
# ADMIN_PASSWORD=your_secure_password

# Install dependencies
go mod download

# Run migrations
psql -h localhost -U postgres -d ghostblog -f migrations/001_create_blogposts.up.sql
psql -h localhost -U postgres -d ghostblog -f migrations/002_create_ratelimits.up.sql

# Run the server
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080`

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will start on `http://localhost:3000`

## Project Structure

```
ghost-blog/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # Application entry point
│   ├── internal/
│   │   ├── config/               # Configuration management
│   │   ├── database/             # Database connection
│   │   ├── handlers/             # HTTP handlers
│   │   ├── middleware/           # Middleware (CORS, rate limiting, auth)
│   │   └── models/               # Data models and repositories
│   ├── migrations/               # Database migrations
│   ├── go.mod
│   └── go.sum
└── frontend/
    ├── src/
    │   ├── assets/
    │   │   └── css/              # Global styles
    │   ├── components/           # Reusable components
    │   ├── router/               # Vue Router configuration
    │   ├── stores/               # Pinia stores
    │   ├── views/                # Page components
    │   ├── App.vue               # Root component
    │   └── main.js               # Application entry point
    ├── index.html
    ├── package.json
    └── vite.config.js
```

## API Endpoints

### Public Endpoints

- `GET /api/posts` - Get random posts
  - Query params: `limit` (default: 10, max: 20), `offset` (default: 0)

- `POST /api/posts` - Create a new post (rate limited)
  - Body: `{ "title": "string", "content": "string" }`

### Admin Endpoints (Basic Auth Required)

- `GET /api/admin/posts` - Get all posts
  - Query params: `limit` (default: 50), `offset` (default: 0)

- `DELETE /api/admin/posts/:id` - Delete a post


## Usage

### Browsing Posts
1. Navigate to the Browse page
2. Scroll with your mouse wheel to navigate between posts
3. Posts are displayed one at a time with snap-scroll behavior

### Submitting Posts
1. Navigate to the Submit page
2. Enter a title (max 255 characters) and content
3. Submit - your country will be automatically detected
4. Rate limit: 5 posts per hour per IP address

### Admin Panel
1. Navigate to `/admin` (not visible in navbar)
2. Login with admin credentials
3. View all posts and delete inappropriate content

## Development

### Backend Development
```bash
cd backend
go run cmd/server/main.go
```

### Frontend Development
```bash
cd frontend
npm run dev
```

### Building for Production

#### Backend
```bash
cd backend
go build -o server cmd/server/main.go
./server
```

#### Frontend
```bash
cd frontend
npm run build
# Serve the dist/ directory with your web server
```

## Database Schema

### blogposts
```sql
id          SERIAL PRIMARY KEY
title       VARCHAR(255) NOT NULL
content     TEXT NOT NULL
country     VARCHAR(100) NOT NULL
created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
```

### ratelimits
```sql
ip_address   VARCHAR(45) PRIMARY KEY
post_count   INTEGER DEFAULT 0
window_start TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
```

## Security Features

- IP-based rate limiting
- Basic authentication for admin panel
- CORS configuration
- SQL injection prevention (parameterized queries)
- Input validation

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
