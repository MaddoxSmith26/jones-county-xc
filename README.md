# Jones County XC

A web application with a React frontend and Go backend.

## Project Structure

```
jones-county-xc/
├── frontend/    # React app (Vite + Tailwind CSS)
├── backend/     # Go HTTP server
├── docs/        # Documentation
└── README.md
```

## Prerequisites

- Node.js 18+
- Go 1.21+

## Getting Started

### Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend runs at http://localhost:5173

### Backend

```bash
cd backend
go run main.go
```

The backend runs at http://localhost:8080

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check |
| `/api/hello` | GET | Hello message |
