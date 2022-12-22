# Protoc (go'ish) Cheatsheet

## Compiler Installation

- Download latest protocol buffer compiler: [Github releases](https://github.com/protocolbuffers/protobuf/releases/tag/v21.12)
- Add the compiler to your PATH
- Download the go plugins:
  - `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28`
  - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2`

---

## Compile projects protocol buffer

The compiler for this project was called as followed in the `go_gRPC/api` directory:  
 `protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  chat.proto`

---

## Generated files

- `*.pb.go` contains all the protocol buffer code to populate, serialize, and retrieve request and response message types
- `_grpc.pb.go` contains the following:
  - An interface type (or stub) for clients to call with the methods defined in the service
  - An interface type for servers to implement, also with the methods defined in the service
  
---

## Compiler options

### **--go_out**

The `--go_out` flag causes the compiler to generate a `*.pb.go` at the given path. **These contain the types** like protobuffer messages. The name of the .proto file will be the name of the outout as well.

```bash
protoc --go_out=<any path>
```

### **--go_opt=paths=source_relative**

The `--go_opt=paths=source_relative` flag causes the compiler to use the path to the .proto file from the argument list. Use `--go_out=` to specify the start directory from where the compiler should create the out. The directory must exist.

- **Example:** When `./go_gRPC/foo/bar/chat.proto` would exist, `protoc --go_out=server --go_opt=paths=source_relative foo/bar/chat.proto` would create `./go_gRPC/server/foo/bar/chat.pb.go`.
- **Example II:** For this project `protoc --go_out=. --go_opt=paths=source_relative api/chat.proto` would create `chat.pb.proto` in the root directory `go_gRPC/api`

```bash
protoc --go_opt=paths=source_relative
```

### **--go-grpc_out**

The `--go-grpc_out` flag causes the compiler to generate a `*_grpc.pb.go` at the given path. **These contain the actual services and stubs**. The name of the .proto file will be used as prefix for the compiled grpc.pb.go file.

```bash
protoc --go-grpc_out=<any path>
```

### **--go-grpc_opt=paths=source_relative**

The `--go_opt=paths=source_relative` flag causes the compiler to use the path to the .proto file from the argument list. Use `--go_out=` to specify the start directory from where the compiler should create the out. The directory must exist.

- **Example:** When `./go_gRPC/foo/bar/chat.proto` would exist, `protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative foo/bar/chat.proto` would create `./go_gRPC/foo/bar/chat.pb.go`.
- **Example II:** For this project `protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative api/chat.proto` would create `chat_grpc.pb.proto` in the root directory `go_gRPC/api`

```bash
protoc --go-grpc_opt=paths=source_relative
```

---

## Other options not used in this project

### **--go_opt=paths=import**

The `--go_opt=paths=import` flag causes the compiler to use the path provided in `option go_package` from the .proto file to determine the output path of the compiled proto file. Use `--go_out=` to specify the start directory from where the compiler should create the out.

- **Example:** When `protoc --go_out=. --go_opt=paths=import chat.proto` is used, the .proto of this project would be generated under `./go_gRPC/api/github.com/penomatikus/golearning/go_gRPC` because the chat.proto has `option go_package = "github.com/penomatikus/golearning/go_gRPC/api";` as go_package.

```bash
protoc --go_opt=paths=import
```

### **--go-grpc_opt=paths=import**

Same as above but for grpc services and stub.

### **--proto_path**

The `proto_path` flag causes the compiler to search for imported files in a set of directories specified. Use this, when you have imports in your proto files. Alternative is `-I`. If not provided, the compiler takes the directory where it was called. [Official Documentation here](https://developers.google.com/protocol-buffers/docs/proto3#importing_definitions)

```
protoc --proto_path=<any path>
```
