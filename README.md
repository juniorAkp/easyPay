# easyPay

easyPay is a Go-based payment processing service that integrates AI-powered message parsing with mobile money payments. It is designed to facilitate transactions through natural language commands via WhatsApp.

## Features

- **AI-Powered Intent Recognition**: Uses Google's Gemini AI (`gemini-2.0-flash-exp`) to extract payment details and user intent from natural language messages.
- **Mobile Money Integration**: Built-in support for MTN Mobile Money (MoMo) API for processing collections and disbursements.
- **WhatsApp Integration**: Designed to handle interactions via WhatsApp.
- **Asynchronous Processing**: Utilizes RabbitMQ for reliable, asynchronous task handling and message consumption.
- **Robust Architecture**: Built with the Gin web framework, PostgreSQL for persistence, and Docker for easy deployment.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin
- **Database**: PostgreSQL
- **Message Queue**: RabbitMQ
- **AI/LLM**: Google Gemini
- **Infrastructure**: Docker, Docker Compose

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- MTN MoMo API Credentials
- Google Gemini API Key

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/juniorAkp/easyPay.git
    cd easyPay
    ```

2.  **Environment Setup:**
    Create a `.env` file in the root directory with the necessary configuration (Database URL, RabbitMQ URL, API Keys, etc.).

3.  **Run with Docker:**

    ```bash
    docker-compose up --build
    ```

4.  **Run Locally:**
    Ensure Postgres and RabbitMQ are running, then:
    ```bash
    go run cmd/main.go
    ```

## Project Structure

- `cmd/`: Application entry point.
- `internal/`: Private application code.
  - `services/`: External service integrations (AI, MoMo, WhatsApp).
  - `handler/`: HTTP request handlers.
  - `consumer/`: RabbitMQ message consumers.
  - `database/`: Database connection and helpers.
  - `queue/`: RabbitMQ setup.
- `pkg/`: Public library code (types, utils).

## Note

This project is still in development. Production release planned for the near future.
