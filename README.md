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
- Go 1.16 or higher
- Docker

## Installation
1. Clone the repository:
    
    git clone https://github.com/yourusername/go-auth-app.git

2. Navigate to the project directory:
    
    cd go-auth-app
    

## Configuration
1. Create a `.env` file in the root directory and add the following environment variables:
    
    DB_HOST=db
    DB_USER=postgres
    DB_PASSWORD=0663058639
    DB_NAME=go_auth_db
    DB_PORT=5432
    JWT_SECRET=your_jwt_secret
    

## Usage
1. Run the application:
    
    docker-compose up --build
    
2. The application will be available at `http://localhost:8080`.

## Running Database Migrations
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
        Both old and new passwords are required
        password must be at least 6 characters long
    - 401 Unauthorized:
        Incorrect old password
        Unauthorized
    - 500 Internal Server Error:
        Failed to hash new password
        Failed to update password


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






