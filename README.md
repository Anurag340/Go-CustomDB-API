# 🛠️ Custom Database API

`GO-Custom-Database-API` is a lightweight, file-based key-value store written in Go, designed for simple database tasks, demos, and small-scale applications. It stores data as JSON files in a directory structure by collections and provides RESTful APIs using the Fiber web framework.

### 🧠 Why customDb?

- **No dependencies** beyond the Go standard library (and Fiber for the HTTP server)
- Uses **`sync.Mutex`** to protect collections from race conditions in concurrent environments, preventing data corruption or unexpected behavior
- Easily extensible for more advanced features like TTL, indexing, or persistence layers

---

## 📸 Demo

[![Watch the video](https://github.com/Anurag340/Go-CustomDB-API/blob/c8791a1f620d07c9fab7a592ffb95c03c2356da4/go-db.png)](https://www.youtube.com/watch?v=lGG62EoVuHw)

---

## 🎥 YouTube Walkthrough

For a full demonstration and explanation of the project, watch the detailed walkthrough:

🔗 [https://www.youtube.com/watch?v=lGG62EoVuHw](https://www.youtube.com/watch?v=lGG62EoVuHw)

---

## 🚀 Features

- ✅ Basic CRUD for user data via REST API
- ✅ File-based storage (each record as a JSON file)
- ✅ Concurrency-safe access with `sync.Mutex`
- ✅ Minimal external dependencies
- 🧩 Easy to extend and customize

---

## 📦 Installation

```bash
git clone https://github.com/yourusername/customDb.git
cd customDb
go mod tidy
go run main.go
