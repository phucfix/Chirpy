# 🐤Chirpy
Chirpy is a mini fake social network similar to Twitter (or X). This project help me learn HTTP server in Golang.
## 🏴Goal
The goal of project is building a HTTP serve that people can register, login and post some messages. They also can view the other user message.
## ⚙️ Installation:
Clone this project:
```bash
git clone github.com/phucfix/chirpy
```
Or using go get tool:
```
```bash
go get github.com/phucfix/chirpy
```
## 🏃 How to run:
In root of the project, run this command:
```bash
go run .
```
## 🧪 How to test:
```bash
go test ./...
```
## API Usage
| Method | Endpoint | Description                   | Authentication |
|--------|----------|-------------------------------|----------------|
|GET     | /app/    | Get files served by the server| No             |
