# Hilo

> Minimal real-time chat with WebSocket.

## Architecture

**Design Pattern**: Domain-Driven Design + Clean Architecture  
**Dependency Injection**: Google Wire  
**Real-time Communication**: WebSocket (Gorilla)  
**State Management**: Riverpod  

### Layers
```
Presentation → Application → Domain ← Infrastructure
    (HTTP)      (Use Case)   (Entity)   (Postgres/WS)
```

## Features

- JWT authentication
- User search
- One-to-one chat rooms
- Real-time messaging via WebSocket
- Message persistence
- Read receipts

## Tech Stack

- Golang 1.25, PostgreSQL 15, WebSocket
- Flutter 3.35, Riverpod, web_socket_channel
