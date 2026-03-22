# Mosque API (Core SaaS Engine)

Backend service for Mosque SaaS platform, built with a Modular Monolith architecture.

## Tech Stack
- **Language:** Go (Golang)
- **Framework:** Go Fiber v2
- **Database:** PostgreSQL (Raw SQL via `pgx/v5`)
- **Migrations:** `golang-migrate`
- **Logging:** Logrus

## Prerequisites
- Go 1.21+
- Docker & Docker Compose
- `golang-migrate` CLI installed

## Local Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd mosque-api
   ```

2. **Setup Environment Variables:**
    
    Copy .env.example to .env (Create .env file and set the DB_URL and APP_PORT).

    ```Code snippet
    APP_PORT=8080
    DB_URL=postgres://root:secretpassword@localhost:5435/mosque_saas?sslmode=disable
    ```

3. **Run Database:**

    Ensure your docker-compose is running.

    ```Bash
    docker compose up -d
    ```
4. **Run Migrations:**

    ```Bash
    make migrate-up
    ```
5. **Start the Server:**

    ```Bash
    make run
    ```