# Split Bill

A simple web application to split bills among friends. It allows users to create items, assign items to each bill participants, and calculate how much each person owes.

## Features
> to be inserted

## Requirements
- Golang 1.18 or later

## Installation
### Normal installation:
```bash
cp .env.example .env # then adjust values
go mod download
go run main.go
```
Also supports Air for live reloading during development:
```bash
# run previous commands first
go install github.com/air-verse/air@latest
air
```
Might need to adjust GOPATH and GOROOT on some cases
### Docker installation:
```bash
# setup env as in normal installation step 1
docker build -t split-bill-be:<commit-hash> .
docker run -d --name split-bill-be-container -p 8080:8080 split-bill-be:<commit-hash>
```

## Docs
API documentation is available at `/swagger/index.html` endpoint after running the application.