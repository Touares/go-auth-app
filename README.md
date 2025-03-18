# Go Auth App

## Overview
A lightweight authentication API built with **Go**, **PostgreSQL**, and **Docker**. It provides secure user authentication, JWT-based authorization, and a clean, structured REST API.

## Features
- User Registration & Login  
- JWT-based Authentication (Access & Refresh Tokens)  
- User Management (Fetch, Soft Delete, Update)  
- Secure Password Hashing  
- SQL-based Database with Migrations Management 
- Full CRUD Operations  
- Dockerized for Easy Deployment  
- Well-Tested with Unit & Integration Tests  

## Tech Stack
- **Golang** (Standard Library, `net/http`)
- **PostgreSQL** (Persistent Data Storage)
- **Golang Migrate** (Schema Migrations)
- **Docker & Docker Compose** (Containerized Development)
- **JWT (JSON Web Tokens)** (Authentication & Authorization)


## Requirements
- Go 1.24.1 or higher
- Docker

## Installation
1. Clone the repository:
    
    git clone https://github.com/yourusername/go-auth-app.git

2. Navigate to the project directory:
    
    cd go-auth-app
    

## Configuration
1. Create a `.env` file in the root directory and add the following environment variables:
    
    # Database Configuration
    DB_HOST="your_db_host"             # e.g., localhost or db (if using Docker)
    DB_USER="your_db_user"             # e.g., postgres
    DB_PASSWORD="your_db_password"     # Replace with your actual database password
    DB_NAME="your_db_name"             # e.g., go_auth_db
    DB_PORT="5432"                     # Default PostgreSQL port

    # JWT Authentication
    JWT_SECRET="your_random_access_token_secret"
    JWT_REFRESH_SECRET="your_random_refresh_token_secret"
    JWT_ACCESS_EXPIRATION=15           # Access token expiration time in minutes
    JWT_REFRESH_EXPIRATION=168         # Refresh token expiration time in hours (7 days)

    # Full Database Connection URL
    DATABASE_URL="postgres://your_db_user:your_db_password@your_db_host:5432/your_db_name?sslmode=disable"

    # Example instructions in README:
    # 1. Copy this file to `.env`
    # 2. Replace `your_db_*` and `your_jwt_*` values with your own
    # 3. Ensure `.env` is NOT committed to version control
    

## Usage
1. Run the application:
    
    docker-compose up --build
    
2. The application will be available at `http://localhost:8080`.

## Running Database Migrations
⚠️ Note: Migrations are automatically applied when running `docker-compose up --build`.  
If you need to manually trigger migrations, use the following commands:
- **Apply Migrations (Up)**

    docker-compose run --rm migrate -path=/migrations -database="postgres://postgres:0663058639@db:5432/go_auth_db?sslmode=disable" up

- **Apply Migrations (Down)**

    docker-compose run --rm migrate -path=/migrations -database="postgres://postgres:0663058639@db:5432/go_auth_db?sslmode=disable" down
    
## API Endpoints
### Register
- **URL:** `/register`
- **Method:** `POST`
- **Body:**
    {
        "name": "your_name",
        "email": "your_email",
        "password": "your_password"
    }
- **Response:**
    201 Created
    {
    "id": 1,
    "name": "John Doe",
    "email": "johndoe@example.com"
    }

- **Possible Errors:**
    - 400 Bad Request: Invalid email format. 
    - Password is too short (must be at least 6 characters).
    - 409 Conflict: Email is already in use.


### Login
- **URL:** `/login`
- **Method:** `POST`
- **Body:**
    {
    "email": "johndoe@example.com",
    "password": "securepassword"
    }
- **Response:**
    200 OK
    {
  "access_token": "eyJhbGciOiJIUzI1NiIsInR...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR..."
    }
- **Possible Errors:**
    - 400 Bad Request: Missing email or password.

    - 401 Unauthorized:
        Incorrect password.
        Email not registered.


###  Fetch All Users
⚠️ Note: In a real-world scenario, this endpoint would likely be restricted to admins.
- **URL:** `/users`
- **Method:** `GET`
- **Headers:**  
    Authorization: Bearer <your_jwt_token>
- **Response:**
    200 OK
    [
  {
    "id": 1,
    "name": "John Doe",
    "email": "johndoe@example.com",
    "is_deleted": false
  }
    ]
- **Possible Errors:**
    - 401 Unauthorized:
        Missing Authorization header.
        Invalid or expired token.


###  Fetch User
- **URL:** `/users/me`
- **Method:** `GET`
- **Headers:**  
    Authorization: Bearer <your_jwt_token>
- **Response:**
   200 OK
  {
    "id": 1,
    "name": "John Doe",
    "email": "johndoe@example.com",
    "is_deleted": false
  }
- **Possible Errors:**
    - 401 Unauthorized:
        Missing Authorization header.
        Invalid or expired token.


###  Update User
- **URL:** `/users/me/update`
- **Method:** `PATCH`
- **Body:**
    {
    "name":"new name"
    }
- **Headers:**  
    Authorization: Bearer <your_jwt_token>
- **Response:**
  200 OK
  {
    "id": 1,
    "name": "John Doe",
    "email": "johndoe@example.com",
    "is_deleted": false
  }
- **Possible Errors:**
    - 400 Bad Request: Name must be at least 3 characters long.
    - 401 Unauthorized:
        Missing Authorization header.
        Invalid or expired token.


###  Reset Password
- **URL:** `/users/me/reset-password`
- **Method:** `POST`
- **Body:**
   {
  "old_password": "securepassword",
  "new_password": "newsecurepassword"
    }
- **Headers:**  
    Authorization: Bearer <your_jwt_token>
- **Response:**
  200 OK
    {
  "message": "Password updated successfully"
    }
- **Possible Errors:**
    - 400 Bad Request:
    - Missing `old_password` or `new_password`
    - New password must be at least 6 characters long  
    - 401 Unauthorized: Incorrect old password  
    - 500 Internal Server Error: Unexpected database or hashing failure  
    
###  Soft Delete User
- **URL:** `/users/me/deactivate`
- **Method:** `DELETE`
- **Headers:**  
    Authorization: Bearer <your_jwt_token>
- **Response:**
    200 OK
    {
  "message": "User account deactivated"
    }
- **Possible Errors:**
    - 401 Unauthorized:
        Missing Authorization header.
        Invalid or expired token.
    - 400 Bad Request: User account is already deactivated.






