# go-kv

A persistent, concurrent, replicated key-value store built from scratch in Go. Inspired by Redis.

Built as a learning project to understand distributed systems from the ground up — TCP networking, concurrency, persistence, replication, and automatic failover.

---

## What it does

- Accepts TCP connections from clients
- Supports `SET key value` and `GET key` commands
- Handles multiple clients concurrently using goroutines
- Persists data to disk using a Write Ahead Log (WAL)
- Survives restarts — data is replayed from disk on startup
- Supports leader/follower replication across two nodes
- Automatic failover via a watchdog and proxy — if the leader dies, traffic is redirected to the follower

---

## How to run

**Always start in this order:**

**1. Start the follower first:**
```bash
go run main.go 5380 Follower
```

**2. Start the leader:**
```bash
go run main.go 5379 Leader
```

**3. Start the watchdog:**
```bash
go run watchdog/watchdog.go
```

**4. Start the proxy:**
```bash
go run proxy/proxy.go
```

**5. Connect a client through the proxy:**
```bash
nc localhost 5378
```

> The follower must be started before the leader so the leader can establish its replication connection on startup.

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
The server can run in two modes — Leader or Follower. The leader accepts client connections normally. Every time it processes a `SET` command, it forwards the same command to the follower over a persistent TCP connection. The follower maintains its own independent copy of the data. If the leader goes down, the follower already has everything.

### Phase 5 — Automatic Failover
Two additional programs handle failover:

**Watchdog** — pings the leader every 3 seconds. After 3 consecutive failures it declares the leader dead and writes the follower's port to `active_node.txt`.

**Proxy** — sits in front of everything and listens on port `5378`. Before each client connection it reads `active_node.txt` to find the currently active node and forwards all traffic there using `io.Copy`. Clients always connect to the proxy and never need to know which node is active.

---

## Architecture

```
Client
  │
  ▼
Proxy (port 5378)
  │  reads active_node.txt
  │  forwards all traffic
  │
  ▼
Leader (port 5379) ←── Watchdog monitors this
  │  accepts connections
  │  handles GET/SET
  │  writes to val.log
  │  forwards SET to follower
  ▼
Follower (port 5380)
  │  accepts forwarded SET commands
  │  maintains its own map
  │  writes to its own val.log
```

**On leader failure:**
```
Watchdog detects failure → writes 5380 to active_node.txt
Proxy reads active_node.txt → forwards new connections to follower
Follower takes over → clients notice nothing on reconnect
```

---

## Project structure

```
go-kv/
├── main.go              # leader/follower server
├── val.log              # write ahead log (auto created)
├── active_node.txt      # current active node port (auto created by watchdog)
├── watchdog/
│   └── watchdog.go      # monitors leader, triggers failover
└── proxy/
    └── proxy.go         # forwards client traffic to active node
```

---

## Concepts learned

- TCP networking — how raw bytes travel between machines
- Goroutines — Go's lightweight concurrency primitive
- Mutex — mutual exclusion to prevent race conditions
- Write Ahead Log — how databases persist data to disk
- Leader/Follower replication — how distributed systems survive node failures
- Watchdog pattern — health checking and failure detection
- Proxy pattern — transparent traffic forwarding and failover

---

## What's next

- Persistent follower reconnection — leader retries connection if follower restarts
- Partitioning — split data across multiple nodes
- Authentication — per-client data isolation
- WAL compaction — clean up the log file over time
- Leader election — automatic promotion without manual watchdog
