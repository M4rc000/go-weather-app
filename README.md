# 🌦️ Weather App — Fullstack Developer Test (Go + PostgreSQL)

This project is a secure fullstack weather tracking app built using Golang, PostgreSQL, JWT, and MeteoSource API.

## 🔧 Features
- User registration & login with JWT
- Secure endpoints with middleware
- Store and manage location data
- Cron scheduler to fetch and save weather data every 1 minute
- CRUD weather data
- Daily and Hourly forecast using MeteoSource API
- Caching mechanism to optimize API limit
- Docker & CI/CD ready (Dockerfile, docker-compose, Jenkinsfile)

## 📦 Tech Stack
- Go (Gin, GORM)
- PostgreSQL
- JWT Auth
- Docker, Jenkins
- MeteoSource API

## 🚀 Getting Started

1. Rename `.env.example` to `.env` and fill in your secrets
2. Start with Docker:

```bash
docker-compose up --build