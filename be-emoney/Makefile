.PHONY: run build test deps clean

# Jalankan server
run:
	go run main.go

# Build binary
build:
	go build -o bin/emoney-server.exe main.go

# Download dependencies & tidy
deps:
	go mod tidy

# Jalankan semua test
test:
	go test ./...

# Bersihkan build artifacts
clean:
	rm -rf bin/
