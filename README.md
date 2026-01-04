# Expenses Tracker PWA

A Progressive Web App for tracking expenses built with Svelte and Golang.

## Features

- **Input Expenses**: Add expenses with categories (daily, monthly, others), date, notes, and amount in Indonesian Rupiah format
- **Expenses History**: View expenses in a calendar format with monthly totals, charts, and detailed date views with CRUD operations

## Frontend Setup

```bash
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173`

## Backend Setup

```bash
cd backend
go mod download
```

Create a `.env` file in the backend directory (see `backend/.env.example`):

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=expenses_db
DB_SSLMODE=disable
PORT=8080
```

Create the PostgreSQL database:
```sql
CREATE DATABASE expenses_db;
```

Run the backend:
```bash
go run main.go
```

The backend API will be available at `http://localhost:8080`

## API Endpoints

- `POST /api/expenses` - Create a new expense
- `GET /api/expenses/months` - Get all months with totals and dates
- `GET /api/expenses/month/:month` - Get month details by category (format: YYYY-MM)
- `GET /api/expenses/date/:date` - Get expenses for a specific date (format: YYYY-MM-DD)
- `PUT /api/expenses/:id` - Update an expense
- `DELETE /api/expenses/:id` - Delete an expense

## PWA Icons

You need to add PWA icons to the `public` directory:
- `public/pwa-192x192.png` (192x192 pixels)
- `public/pwa-512x512.png` (512x512 pixels)

You can generate these using online PWA icon generators.

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

