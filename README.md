# TCP-UPD-SERVERS-CLIENTS

# Go Chat Server Project
https://www.youtube.com/watch?v=fzW6bC9AXPQ

This project contains TCP and UDP chat servers and clients implemented in Go.

## Files

- `tcp_server.go`: TCP chat server
- `tcp_client.go`: TCP chat client
- `udp_server.go`: UDP chat server
- `udp_client.go`: UDP chat client

## Requirements

- Go 1.18 or higher

## Build and Run

### TCP Server

```bash
go run tcp_server.go
```

### TCP Client

In a separate terminal:

```bash
go run tcp_client.go
```

### UDP Server

```bash
go run udp_server.go
```

### UDP Client

In a separate terminal:

```bash
go run udp_client.go
```

## Features

- Clients can connect and send/receive messages
- Messages are broadcast to all connected clients
- Disconnection is handled gracefully
- Uses goroutines for concurrency
- TCP server logs client messages to files
- UDP server manages clients by address and broadcasts messages
- Commands supported in TCP server:
  - `/quit` or `bye`: disconnect
  - `/time`: show server time
  - `/date`: show server date
  - `/joke`: tell u to go study
  - `/clients`: show number of connected clients
  - `/help`: show available commands

## Notes

- Each server and client is a standalone Go program.
- Compile and run each program separately.
- UDP is connectionless; client management is based on message receipt.
- TCP server logs client messages in `client_logs` directory.

## Testing & Optimization

- You can simulate different networking conditions using tools like `tc` on Linux.
- Monitor latency, packet loss, and throughput using network tools.
- Test edge cases like sudden client disconnects and high message volume.
