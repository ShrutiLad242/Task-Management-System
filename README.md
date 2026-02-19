# Task Management System (Golang)

## Overview

This project is a scalable RESTful Task Management Service built using Go. It demonstrates clean architecture, JWT authentication, PostgreSQL persistence, background workers using goroutines and channels, and proper separation of concerns.

The system allows users to create, view, and delete tasks. Tasks are automatically processed and marked as completed by a background worker after a configurable delay.

---

## Features

### Core Features

- Create task
- Get all tasks
- Get task by ID
- Delete task
- PostgreSQL persistence
- Clean architecture (Controller → Service → Repository → DB)
- Background worker with goroutines and channels
- Automatic task completion after configurable delay
- JWT authentication
- Authorization (users can access only their own tasks)
- Environment-based configuration
- Graceful shutdown

---

## Task Lifecycle

Each task follows this lifecycle:

