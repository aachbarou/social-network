# Full-Stack Social Network Platform

A fully functional, real-time social networking platform engineered for high concurrency, instant communication, and scalable group management.

## 🚀 System Architecture

This project is separated into two discrete environments allowing for microservice-style deployments and strict separation of concerns:

- **/frontend:** A reactive Single Page Application (SPA) orchestrating the interface using Vue.js 3, Vite, and Pinia State Management.
- **/backend:** A high-performance RESTful API and WebSocket broadcast hub written entirely in Go (Golang) backed by SQLite.

## 🌟 Core Features

- **Real-Time Global Chat & Direct Messaging:** Instantaneous delivery using persistent WebSocket connections.
- **Dynamic Notifications:** Live updates for new followers, group invites, event RSVPs, and chat messages.
- **Social Graph Hierarchy:** Follow/Unfollow users, track active followers, and build custom timelines.
- **Groups & Events:** Create managed groups, invite members, securely post content, and schedule RSVP-based events.
- **Authentication & Security:** Robust encrypted JWT or Session-based security layers bridging frontend clients and the backend APIs.

## ⚙️ Getting Started

You will need to run both the frontend and backend servers concurrently to operate the application locally.

### 1. Launch the Backend
```bash
cd backend
go mod tidy
go run server.go
```
*Runs on port `:8081`.*

### 2. Launch the Frontend
In a new terminal wrapper:
```bash
cd frontend
npm install
npm run dev
```
*Runs on port `:5173`.*

## 🐳 Docker Deployment

For reproducible containerized environments, the root project also includes scripts and configuration files for a complete Docker engine deployment.

Please see the individual `/frontend/README.md` and `/backend/README.md` files for deeper architecture and API documentation layouts.
