# Todo 

This is a simple, yet powerful Todo built with React for the frontend, Golang (Fiber framework) for the backend, and SQLC for type-safe database operations.

## Features

- Create new todos
- List all todos
- Mark todos as completed
- Delete todos
- Responsive design

## Tech Stack

- Frontend: React with TypeScript
- Backend: Golang with Fiber framework
- Database: MySQL
- ORM: SQLC for type-safe SQL
- API: RESTful API
- Styling: CSS

## Prerequisites

- Node.js and npm
- Go 1.22 or later
- MySQL

## Setup and Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/todo-app.git
   cd todo-app
   ```

2. Set up the backend:
   ```
   cd backend
   go mod download
   ```

3. Set up the database:
   - Create a MySQL database
   - Update the `.env` file with your database credentials

4. Generate SQLC code:
   ```
   sqlc generate
   ```

5. Run the backend:
   ```
   go run main.go
   ```

6. Set up the frontend:
   ```
   cd ../frontend
   npm install
   ```

7. Run the frontend:
   ```
   npm run dev
   ```

8. Open your browser and navigate to `http://localhost:5173` (or the port specified by Vite)

## API Endpoints

- GET /api/todos: Get all todos
- POST /api/todos: Create a new todo
- PATCH /api/todos/:id: Toggle todo completion status
- DELETE /api/todos/:id: Delete a todo

