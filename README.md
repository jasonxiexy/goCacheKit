# Distributed Cache System in Golang

A lightweight distributed caching system implemented in Go, designed for performance, scalability, and simplicity.

## 🔧 Features

- **Local + HTTP-based Distributed Caching**
- **LRU (Least Recently Used) Eviction Policy**
- **Cache Breakdown Protection with Go Locks**
- **Consistent Hashing for Node Selection & Load Balancing**
- **Efficient Inter-Node Communication using Protobuf**

## Single Node Cache Retrieval Logic Flow

1. Receive a key  
   ↓  
2. Check if the key is cached  
   - **Yes** → Return cached value (①)  
   - **No** →  
     - Should we fetch from a remote node?  
       - **Yes** → Interact with remote node → Return cached value (②)  
       - **No** → Call a `callback function` to get value → Add to cache → Return cached value (③)

## 🚀 Getting Started

1. **Clone the repo:**
   ```bash
   git clone https://github.com/jasonxiexy/goCacheKit.git
   cd goCacheKit
   ```

2. **Run test in each block**
   ```bash
   cd xxx
   go test -v
   ```

3. **Test HTTP method to interact with Cache**
   
   In one terminal:
   ```bash
   go run main.go
   ```
   Another terminal:
   ```bash
   curl http://localhost:9999/_gocache/scores/Tom
   curl http://localhost:9999/_gocache/scores/kkk
   ```
   And expected results:
   ```bash
   PS D:\goCacheKit> go run main.go
   2025/05/28 15:31:11 gocache is running at localhost:9999
   2025/05/28 15:31:16 [Server localhost:9999] GET /_gocache/scores/Tom
   2025/05/28 15:31:16 [SlowDB] search key Tom
   2025/05/28 15:32:06 [Server localhost:9999] GET /_gocache/scores/kkk
   2025/05/28 15:32:06 [SlowDB] search key kkk
   ```