<p align="center"> <b> Gedis - Redis-compatible Cache Server  </b> </p>


<p align="center">
  <img src="https://github.com/user-attachments/assets/631e6f23-86de-46d6-ad4b-9ed7abccc769" alt="Gedis Logo" width="500" height="500">
</p>

<p align="center">
  <b>A high-performance in-memory caching solution written in Go</b>
</p>

<p align="center">
  <a href="#key-features">Features</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#command-reference">Commands</a> •
  <a href="#client-project">Client</a>
</p>

## Overview

Gedis is a lightweight, high-performance in-memory caching server that implements the Redis RESP (REdis Serialization Protocol). It's designed to be simple, fast, and compatible with existing Redis clients while providing core functionality needed for modern application caching.

## Key Features

✅ **Redis Protocol Compatible** - Works with existing Redis clients  
✅ **Performant** - Written in Go for high concurrency and speed  
✅ **Comprehensive Data Structure Support**:
  - Key-value strings
  - Lists for queue/stack operations
  - Sets for unique collections
  - Hashes for structured data
  - Sorted sets for priority queues

✅ **Key Expiration** - TTL support for automatic cache invalidation  
✅ **Atomic Operations** - Increment/decrement and other atomic commands  
✅ **Key Watching** - Monitor changes for transaction safety

## Architecture

Gedis follows a simple client-server architecture:

<p align="center">
  <img src="https://github.com/user-attachments/assets/94f93cd9-6a73-43b6-b301-d2206b417f96" alt="Gedis Logo" width="1000" height="500">
</p>


The server handles connections using Go's concurrency primitives, with each client connection processed in its own goroutine. Commands are parsed using the RESP protocol implementation and executed against the in-memory data store.

## Installation

### Prerequisites

- Go 1.24+

### From Source

```bash
# Clone the repository
git clone https://github.com/GedisCaching/Gedis.git
cd Gedis

# Build and run
make run
```

## Quick Start

### Connect with Gedis CLI

```bash
- Start the server
- telnet 127.0.0.1 7000
```

### Basic Operations

```
SET name John
OK
GET name
"John"
```

### Using with a Go Application

```go
package main

import (
  gedis "github.com/GedisCaching/Gedis/gedis"
)

func main() {
  GedisClient, err := gedis.NewGedis(gedis.Config{
    Address:  Address,
    Password: Password,
  })
}
```

## Command Reference

Gedis supports a subset of Redis commands, organized by data type:

### String Operations
- `GET key` - Get the value of a key
- `SET key value` - Set the value of a key
- `DEL key` - Delete a key

### List Operations
- `LPUSH key value [value ...]` - Add values to the head of a list
- `RPUSH key value [value ...]` - Add values to the tail of a list
- `LPOP key` - Remove and get the first element of a list
- `RPOP key` - Remove and get the last element of a list
- `LLEN key` - Get the length of a list
- `LRANGE key start stop` - Get elements from a list

### Hash Operations
- `HSET key field value` - Set the value of a hash field
- `HGET key field` - Get the value of a hash field
- `HDEL key field [field ...]` - Delete fields from a hash
- `HGETALL key` - Get all fields and values from a hash

### Set Operations
- `SADD key member [member ...]` - Add members to a set
- `SREM key member [member ...]` - Remove members from a set
- `SMEMBERS key` - Get all members of a set
- `SISMEMBER key member` - Check if a member exists in a set

### Sorted Set Operations
- `ZADD key score member` - Add members to a sorted set
- `ZREM key member [member ...]` - Remove members from a sorted set
- `ZRANGE key start stop` - Get elements from a sorted set

### Key Management
- `TTL key` - Get the time to live of a key
- `EXPIRE key seconds` - Set the expiration time of a key
- `DEL key [key ...]` - Delete keys

### Numeric Operations
- `INCR key` - Increment the value of a key
- `DECR key` - Decrement the value of a key

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add some amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).

---

## Client-Project

[Gedis-Client](https://github.com/GedisCaching/Gedis-Client) is a Go-based client for interacting with the Gedis caching server, a Redis-like in-memory data store. This client provides a simple and intuitive way to interact with Gedis server, allowing you to perform various operations like basic key-value operations, list operations, hash operations, sorted sets, and more.
