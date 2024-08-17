BIN_DIR = bin
BIN_FILE = server
CMD_DIR = cmd
CMD_FILE = main.go

all: build

run: build
	@./$(BIN_DIR)/$(BIN_FILE)

build:
	@echo "> start building the server..."
	@go build -o $(BIN_DIR)/$(BIN_FILE) $(CMD_DIR)/$(CMD_FILE)
	@echo "> finished building"

clean:
	@echo "> cleaning bin dir..."
	@rm -rf $(BIN_DIR)
	@echo "> finished cleaning"
