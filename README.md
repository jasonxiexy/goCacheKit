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
   git clone https://github.com/yourusername/go-distributed-cache.git
   cd goCacheKit/goCache

2. **Run test in each block**
   ```bash
   cd xxx
   go test -v
