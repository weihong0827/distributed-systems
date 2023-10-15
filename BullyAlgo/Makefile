OUT_DIR := bin
APP := bully

$(shell mkdir -p $(OUT_DIR))

build:
	go build -o $(OUT_DIR)/$(APP) main.go

run: build
	./$(OUT_DIR)/$(APP) -case=1
	./$(OUT_DIR)/$(APP) -case=2
	./$(OUT_DIR)/$(APP) -case=3
	./$(OUT_DIR)/$(APP) -case=4
	./$(OUT_DIR)/$(APP) -case=5
	./$(OUT_DIR)/$(APP) -case=6
	./$(OUT_DIR)/$(APP) -case=7

case1: build
	./$(OUT_DIR)/$(APP) -case=1

case2: build
	./$(OUT_DIR)/$(APP) -case=2

case3.a: build
	./$(OUT_DIR)/$(APP) -case=3

case3.b: build
	./$(OUT_DIR)/$(APP) -case=4

case4: build
	./$(OUT_DIR)/$(APP) -case=5

case5.a: build
	./$(OUT_DIR)/$(APP) -case=6

case5.b: build
	./$(OUT_DIR)/$(APP) -case=7

# Additional target to clean the build
clean:
	rm -rf $(OUT_DIR)

.PHONY: build run case1 case2 case3.a case3.b case4 case5.a case5.b clean
