# Compiler and flags
CXX      := g++
CXXFLAGS := -std=c++17 -Wall -Iinclude -lsystemc -lm

# Directories
SRC_DIR  := src
INC_DIR  := include
BUILD_DIR := build

# Find all source and header files recursively
SRC_FILES := $(shell find $(SRC_DIR) -name '*.cpp')
HDR_FILES := $(shell find $(INC_DIR) -name '*.h' -o -name '*.hpp')

# Object files go in build/, mirroring src/ structure
OBJ_FILES := $(patsubst $(SRC_DIR)/%.cpp, $(BUILD_DIR)/%.o, $(SRC_FILES))

# Output binary
TARGET := $(BUILD_DIR)/main

# Default rule
all: $(TARGET)

# Link step
$(TARGET): $(OBJ_FILES)
	@mkdir -p $(dir $@)
	$(CXX) $(OBJ_FILES) -o $@ -lsystemc -lm

# Compile step
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.cpp $(HDR_FILES)
	@mkdir -p $(dir $@)
	$(CXX) -c $< -o $@ -I$(INC_DIR) -Wall -std=c++17

# Clean rule
clean:
	rm -rf $(BUILD_DIR)

.PHONY: all clean
