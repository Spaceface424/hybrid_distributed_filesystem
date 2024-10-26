# Makefile for compiling Protocol Buffers

# The protoc command to run
PROTOC = protoc
PINGACK_PROTO_FILE = shared/pingack.proto
GREP_PROTO_FILE = grep/grep.proto
OUTPUT_DIR = .

# The default target
all: compile_proto

# Compile the Protocol Buffer file
compile_proto:
	$(PROTOC) --go_out=$(OUTPUT_DIR) --go-grpc_out=$(OUTPUT_DIR) $(PINGACK_PROTO_FILE)
	$(PROTOC) --go_out=$(OUTPUT_DIR) --go-grpc_out=$(OUTPUT_DIR) $(GREP_PROTO_FILE)

# Clean generated files (optional)
clean:
	rm -f $(OUTPUT_DIR)/*.pb.go

.PHONY: all compile_proto clean
