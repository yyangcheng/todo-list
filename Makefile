BUILD_DIR = $(CURDIR)/output

clean:
	rm -rf $(BUILD_DIR)

build:
	go build -o $(BUILD_DIR)/main main.go

compile:
    # 64-Bit
	# FreeBDS
	GOOS=freebsd GOARCH=amd64 go build -o $(BUILD_DIR)/main-freebsd-amd64 main.go
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/main-darwin-amd64 main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/main-linux-amd64 main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/main-windows-amd64 main.go

run:
	go run main.go