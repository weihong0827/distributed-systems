# Set the name of the binary
BIN_DIR=bin

# The first rule is run by default when you execute "make" with no arguments
all: build

# Default value for NUM_CLIENTS
NUM_CLIENTS?=10

# Create the binary directory
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Compile the program and place the executable in the binary directory
build: $(BIN_DIR)
	go build -o $(BIN_DIR) main.go server.go client.go
	chmod +x $(BIN_DIR)/main

# Run the program
run3: build
	./$(BIN_DIR)/main --numClients $(NUM_CLIENTS) --question 1.3

run2: build
	./$(BIN_DIR)/main --numClients $(NUM_CLIENTS) --question 1.2

run1: build
	./$(BIN_DIR)/main --numClients $(NUM_CLIENTS) 
# Clean up
clean:
	rm -rf $(BIN_DIR)

