# Notification System
This project is a microservices-based notification system built with Go, Fiber, and Kafka. It is designed to handle user authentication, friend requests, and real-time notifications efficiently. The system is divided into two main services: `user_service` and `notification_service`, each responsible for specific functionalities.

The `user_service` manages user-related operations such as registration, login, and friend management. It ensures secure authentication and authorization using JWT tokens. Users can send, accept, and manage friend requests seamlessly. This service also handles user sessions and maintains a blacklist of JWT tokens to ensure secure logout operations.

The `notification_service` is responsible for handling notifications. It listens to events from Kafka and stores notifications in the database. Users receive real-time notifications for various activities, ensuring they are always up-to-date with the latest events. This service processes incoming Kafka messages, stores them as notifications, and provides endpoints for users to retrieve and manage their notifications.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)

## Features

- User authentication and authorization with JWT
- Friend request management
- Real-time notifications using Kafka
- RESTful API with Fiber
- MySQL database integration with GORM

## Architecture

The system is divided into two main services:

1. **User Service**: Manages user authentication, friend requests, and user-related data.
2. **Notification Service**: Listens to Kafka events and stores notifications in the database.

## Folder Structure

## Installation

### Prerequisites

- Go 1.23.5 or higher
- MySQL
- Kafka

### Steps

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/notification_system.git
    cd notification_system
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Set up the MySQL database and update the DSN in the `.env` files for both services.

4. Start Kafka.

5. Run the services:
    ```sh
    go run user_service/main.go
    go run notification_service/main.go
    go run main.go
    ```

## Usage

### Running the Services

- **User Service**: Runs on port `8000`
- **Notification Service**: Runs on port `6000`

### API Gateway

The API Gateway runs on port `8080` and forwards requests to the appropriate service based on the URL path.

## API Endpoints

### User Service

- `POST /user/auth/signup`: Register a new user
- `POST /user/auth/login`: Login and receive a JWT token
- `POST /user/auth/logout`: Logout and blacklist the JWT token
- `POST /user/friend/send`: Send a friend request
- `GET /user/friend/requests`: View friend requests
- `POST /user/friend/accept`: Accept a friend request
- `GET /user/friend/list`: View friends
- `POST /user/friend/unfriend`: Unfriend a user
- `POST /user/friend/unfollow`: Unfollow a user
- `POST /user/friend/follow-again`: Follow a user again

### Notification Service

- `GET /notifications`: Get all notifications for a user
- `PUT /notifications/:notification_id/read`: Mark a notification as read
