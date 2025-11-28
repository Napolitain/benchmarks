# Hello World HTTP 1.1 API Benchmarks

This repository contains multiple Hello World HTTP 1.1 implementations for benchmarking purposes.

## Applications

All applications run on **port 8080** in **single-threaded mode** (run one at a time for benchmarking).

### go-http (Go net/http)
```bash
cd go-http
go run main.go
```

### go-fasthttp (Go fasthttp)
```bash
cd go-fasthttp
go mod download
go run main.go
```

### python-fastapi (Python FastAPI with Uvicorn)
```bash
cd python-fastapi
pip install -r requirements.txt
python main.py
```

### python-flask (Python Flask)
```bash
cd python-flask
pip install -r requirements.txt
python main.py
```

### node-http (Node.js http)
```bash
cd node-http
npm start
```

### nginx-static (Nginx static file)
```bash
cd nginx-static
nginx -p . -c nginx.conf
# To stop: nginx -p . -c nginx.conf -s stop
```

### cpp-uwebsockets (C++ uWebSockets)
```bash
# Built by benchrunner, run from bin directory
../../bin/HelloWorldBenchmark
```

### java-springboot (Java Spring Boot)
```bash
cd java-springboot
mvn spring-boot:run
```

## Testing

Each API responds to GET requests at the root endpoint (`/`) with:
```json
{"message": "Hello, World!"}
```

Test the API:
```bash
curl http://localhost:8080/
```
