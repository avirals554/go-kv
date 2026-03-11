# Go-KV

A simple, in-memory key-value store built with Go. This project demonstrates the evolution from a basic map implementation to a thread-safe version using Go's concurrency primitives.

## Features

- **Basic Operations**: Support for `Set`, `Get`, and `Delete` operations.
- **Concurrency Safety**: Uses `sync.RWMutex` to ensure safe access across multiple goroutines.
- **In-Memory**: Fast, volatile storage for quick data operations.

## Installation

To use this project, ensure you have Go installed on your machine.

```bash
git clone [https://github.com/avirals554/go-kv.git](https://github.com/avirals554/go-kv.git)
cd go-kv
