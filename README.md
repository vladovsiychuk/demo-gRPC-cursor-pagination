# Overview
This repository demonstrates a discovery-service for a social media platform, designed to handle and paginate user posts. In this system, the post-service sends new posts asynchronously via gRPC to the discovery-service, which maintains all posts in sorted order. Clients can then fetch a paginated feed of posts from the discovery-service using cursor-based pagination. This project focuses on showcasing the implementation and features of the discovery-service.


## Project Structure

```
.
├── cmd/                # Main entrypoint and integration tests
│   ├── main.go         # Application bootstrap
│   └── integration_test.go
├── discovery/          # Core service logic
│   ├── models.go       # Data models
│   ├── server.go       # gRPC server implementation
│   └── service.go      # Service logic
├── protob/             # Generated protobuf Go files
│   └── discovery/v1/   # Generated code for v1 API
├── protodef/           # Protobuf definitions
│   └── discovery/v1/   # .proto files for v1 API
├── util/               # Utility functions
│   └── id.go           # ID generation helpers
├── go.mod, go.sum      # Go module files
├── Makefile            # Common build/test commands
└── README.md           # Project documentation
```

## Getting Started

### Prerequisites
- Go 1.19+
- `protoc` (Protocol Buffers compiler)
- `buf` (for proto code generation)

### Install Dependencies
```bash
go mod tidy
```

### Generate Protobuf Code
```bash
cd protodef
buf generate
```

### Build
```bash
make build
```

### Run
```bash
make run
```

### Test
```bash
make test
```

## Development Notes
- Protobuf definitions are in `protodef/discovery/v1/`.
- Generated Go code is in `protob/discovery/v1/`.
- Main service logic is in the `discovery/` package.
- Use the Makefile for common tasks.

