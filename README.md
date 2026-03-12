# go-kv

A persistent, concurrent, replicated key-value store built from scratch in Go. Inspired by Redis.

Built as a learning project to understand distributed systems from the ground up — TCP networking, concurrency, persistence, and replication.

---

## What it does

- Accepts TCP connections from clients
- Supports `SET key value` and `GET key` commands
- Handles multiple clients concurrently using goroutines
- Persists data to disk using a Write Ahead Log (WAL)
- Survives restarts — data is replayed from disk on startup
- Supports leader/follower replication across two nodes

---

## How to run

### Single node
```bash
go run main.go 5379 Follower
```

### Replication (two nodes)

Start the follower first:
```bash
go run main.go 5380 Follower
```

Then start the leader:
```bash
go run main.go 5379 Leader
```

### Connect a client
```bash
nc localhost 5379
```

---

## Commands

```
SET name aviral     # stores "aviral" under key "name"
GET name            # returns "aviral"
```

---

## How it works

### Phase 1 — TCP Server
The server listens on a TCP port and accepts incoming connections. Each connection gets a dedicated read loop that parses raw bytes into commands. Commands are split by spaces to extract the operation, key, and value.

### Phase 2 — Concurrency
Each incoming connection is handled in its own goroutine so multiple clients can connect simultaneously without blocking each other. A mutex protects the shared in-memory map from race conditions — without it, concurrent reads and writes corrupt the map and crash the program.

### Phase 3 — Persistence
Every `SET` command is appended as a line to `val.log` on disk. On startup, the server reads this file line by line and replays every command to rebuild the in-memory map. This means data survives server restarts.

### Phase 4 — Replication
The server can run in two modes — Leader or Follower. The leader accepts client connections normally. Every time it processes a `SET` command, it forwards the same command to the follower over a persistent TCP connection. The follower maintains its own copy of the data independently. If the leader goes down, the follower already has all the data.

---

## Architecture

```
Client
  │
  ▼
Leader (port 5379)
  │  accepts connections
  │  handles GET/SET
  │  writes to val.log
  │  forwards SET to follower
  ▼
Follower (port 5380)
  │  accepts forwarded commands
  │  maintains its own map
  │  writes to its own val.log
```

---

## Concepts learned

- TCP networking — how raw bytes travel between machines
- Goroutines — Go's lightweight concurrency primitive
- Mutex — mutual exclusion to prevent race conditions
- Write Ahead Log — how databases persist data to disk
- Leader/Follower replication — how distributed systems survive node failures

---

## What's next

- Automatic failover — detect when leader dies and promote follower
- Partitioning — split data across multiple nodes
- Authentication — per-client data isolation
- Compaction — clean up the WAL file over time
