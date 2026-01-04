# Setup Instructions

## Prerequisites

- Node.js (v18 or higher)
- Go (v1.21 or higher)
- PostgreSQL (v12 or higher)

## Frontend Setup

1. Install dependencies:
```bash
npm install
```

2. Start development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:5173`

## Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. Create `.env` file (copy from `.env.example`):
```bash
cp .env.example .env
```

4. Update `.env` with your PostgreSQL credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=expenses_db
DB_SSLMODE=disable
PORT=8080
```

5. Create PostgreSQL database:
```sql
CREATE DATABASE expenses_db;
```

6. Run the backend:
```bash
go run main.go
```

The backend API will be available at `http://localhost:8080`

## API Endpoints

- `POST /api/expenses` - Create a new expense
- `GET /api/expenses/months` - Get all months with totals
- `GET /api/expenses/month/:month` - Get month details by category (format: YYYY-MM)
- `GET /api/expenses/date/:date` - Get expenses for a specific date (format: YYYY-MM-DD)
- `PUT /api/expenses/:id` - Update an expense
- `DELETE /api/expenses/:id` - Delete an expense

## Building for Production

### Frontend
```bash
npm run build
```

### Backend
```bash
cd backend
go build -o expenses-tracker
```

## PWA Icons

You'll need to add PWA icons:
- `public/pwa-192x192.png` (192x192 pixels)
- `public/pwa-512x512.png` (512x512 pixels)

You can generate these using online PWA icon generators or create them manually.

