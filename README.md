# Go Microservices Demo

This is a simple microservices demo that showcases database synchronization between two services using NATS for event communication. The system demonstrates key microservices patterns including event-driven architecture, database per service, and eventual consistency.

## System Architecture

### Overview
The system implements a user management system with two separate services that maintain synchronized data using events:

```
┌──────────────┐         ┌─────────┐         ┌────────────────────┐
│  User Service├────────►│  NATS   ├────────►│Notification Service│
│  (Postgres)  │  events │ Server  │  sub    │     (MongoDB)      │
└──────────────┘         └─────────┘         └────────────────────┘
```

The system consists of two microservices:

1. **User Service** (Port 8080)
   - Uses PostgreSQL database
   - Handles user creation and updates
   - Publishes events when user data changes

2. **Notification Service** (Port 8081)
   - Uses MongoDB database
   - Subscribes to user events
   - Maintains a synchronized copy of user data

## Prerequisites

- Go 1.16 or later
- PostgreSQL
- MongoDB
- NATS Server

## Setup

1. Start the NATS Server:
```bash
nats-server
```

2. Create PostgreSQL database and table:
```bash
createdb users_db
psql users_db < user-service/schema.sql
```

3. Start MongoDB:
```bash
mongod
```

4. Configure the services:
   - Update database connection strings in both services if needed
   - Default PostgreSQL connection: localhost:5432
   - Default MongoDB connection: mongodb://localhost:27017
   - Default NATS connection: localhost:4222

## Running the Services

1. Start the User Service:
```bash
cd user-service
go run main.go
```

2. Start the Notification Service:
```bash
cd notification-service
go run main.go
```

## API Endpoints

### User Service (localhost:8080)
- `GET /users` - List all users
  ```bash
  curl http://localhost:8080/users
  ```
  Response:
  ```json
  [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "updated_at": "2025-05-25T10:00:00Z"
    }
  ]
  ```

- `POST /users` - Create a new user
  ```bash
  curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
  ```
  Response:
  ```json
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "updated_at": "2025-05-25T10:00:00Z"
  }
  ```

- `PUT /users/{id}` - Update a user
  ```bash
  curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe Updated",
    "email": "john.updated@example.com"
  }'
  ```
  Response:
  ```json
  {
    "id": 1,
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "updated_at": "2025-05-25T10:01:00Z"
  }
  ```

### Notification Service (localhost:8081)
- `GET /notifications` - List all notifications
  ```bash
  curl http://localhost:8081/notifications
  ```
  Response:
  ```json
  [
    {
      "user_id": 1,
      "name": "John Doe Updated",
      "email": "john.updated@example.com",
      "updated_at": "2025-05-25T10:01:00Z"
    }
  ]
  ```

- `GET /notifications/user/{id}` - Get user notification data
  ```bash
  curl http://localhost:8081/notifications/user/1
  ```
  Response:
  ```json
  {
    "user_id": 1,
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "updated_at": "2025-05-25T10:01:00Z"
  }
  ```

## How it Works

### Event-Driven Communication
1. When a user is created or updated in the User Service:
   - Data is first saved in PostgreSQL
   - An event is published to NATS with the following format:
     ```json
     {
       "type": "user_created",  // or "user_updated"
       "payload": {
         "id": 1,
         "name": "John Doe",
         "email": "john@example.com",
         "updated_at": "2025-05-25T10:00:00Z"
       }
     }
     ```

2. The Notification Service:
   - Subscribes to the "user.events" topic on NATS
   - Receives the events in real-time
   - Updates its MongoDB database to maintain a synchronized copy

### Database Synchronization
- **User Service (PostgreSQL)**
  - Primary source of user data
  - Handles CRUD operations
  - Maintains data integrity with constraints (e.g., unique email)
  - Schema:
    ```sql
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
    ```

- **Notification Service (MongoDB)**
  - Maintains a denormalized copy of user data
  - Optimized for read operations
  - Uses upsert operations for atomic updates
  - Collection structure:
    ```json
    {
      "user_id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "updated_at": "2025-05-25T10:00:00Z"
    }
    ```

### Benefits of this Architecture
1. **Loose Coupling**: Services communicate through events, not direct calls
2. **Independent Scaling**: Each service can be scaled independently
3. **Technology Freedom**: Each service uses the best database for its needs
4. **Resilience**: Services can continue working even if others are temporarily down
5. **Eventually Consistent**: Data stays synchronized across services automatically

### Error Handling
- Duplicate email addresses are prevented by PostgreSQL constraints
- Failed events can be replayed from NATS if needed
- Each service handles its own errors independently

This implementation demonstrates key microservices patterns including:
- Event-Driven Architecture
- Database per Service
- Eventually Consistent Data
- Service Independence
- Polyglot Persistence (PostgreSQL + MongoDB)

## Monitoring and Troubleshooting

### Service Status
You can check if services are running:
```bash
# Check NATS Server
curl http://localhost:8222/varz

# Check PostgreSQL
psql -U postgres -d users_db -c "SELECT version();"

# Check MongoDB
mongosh --eval "db.serverStatus()"
```

### Common Issues and Solutions

1. **Service Won't Start**
   - Check if the required ports are available (8080, 8081, 4222)
   - Verify database connections
   - Check NATS server is running

2. **Data Not Syncing**
   - Verify NATS connection in both services
   - Check event publishing in User Service logs
   - Check event subscription in Notification Service logs

3. **Database Connection Issues**
   - PostgreSQL: Check credentials and database existence
   - MongoDB: Verify MongoDB service is running
   - NATS: Ensure NATS server is running and accessible

### Debugging Tips
1. Watch NATS traffic:
   ```bash
   nats sub "user.events" -s http://localhost:4222
   ```

2. Monitor PostgreSQL queries:
   ```bash
   psql -U postgres users_db -c "SELECT * FROM users;"
   ```

3. Check MongoDB data:
   ```bash
   mongosh "mongodb://localhost:27017/notifications_db" --eval "db.user_notifications.find()"
   ```

## Development and Extension

### Adding New Features
1. Add new fields:
   - Update PostgreSQL schema
   - Modify User model
   - Update event payload
   - Update MongoDB schema
   - Update handlers in both services

2. Add new events:
   - Define new event type in shared package
   - Add publisher in User Service
   - Add subscriber in Notification Service

### Best Practices
1. Always validate input data
2. Handle database errors gracefully
3. Implement proper logging
4. Use transactions where necessary
5. Follow idempotency patterns
6. Implement retry mechanisms for failed operations
