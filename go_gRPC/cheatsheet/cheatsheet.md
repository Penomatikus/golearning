# Protoc (go'ish) Cheatsheet

## Project compiler options

The compiler for this project was called as followed in the `go_gRPC` directory: `protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  chat.proto`

---

## Compiler options

- _Most of the compiler snippets only showing a specific flag and will not work when used in the terminal._
- _compiler is used as short for protocol buffer compiler_
- _All examples are run from `./go_gRPC` as root directory_

The `proto_path` flag causes the compiler to search for imported files in a set of directories specified. Use this, when you have imports in your proto files. Alternative is `-I`. If not provided, the compiler takes the directory where it was called. [Official Documentation here](https://developers.google.com/protocol-buffers/docs/proto3#importing_definitions)

```
protoc --proto_path=<any path>
```

</br>

The `--go_out` flag causes the compiler to generate a go output file at the given path. The name of the .proto file will be the name of the outout as well.

```
protoc --go_out=<any path>
```

</br>

---

### The difference: --go_opt VS --go-grpc_opt

The following arguments will create `*.pb.go` files. **These contain the types** like protobuffer messages. In every example below `--go_opt` can be replaced with `--go-grpc_opt`, which will create `_grpc.pb.go` files. **These contain the actual services and stubs**.

---

</br>

The `--go_opt=paths=import` flag causes the compiler to use the path provided in `option go_package` from the .proto file to determine the output path of the compiled proto file. Use `--go_out=` to specify the start directory from where the compiler should create the out.

- **Example:** When `protoc --go_out=server --go_opt=paths=import chat.proto` is used, the .proto of this project would be generated under `./go_gRPC/server/github.com/penomatikus/golearning/go_gRPC`

```
protoc --go_opt=paths=import
```

<br>

The `--go_opt=paths=source_relative` flag causes the compiler to use the path to the .proto file from the argument list. Use `--go_out=` to specify the start directory from where the compiler should create the out. The directory must exist.

- **Example:** When `./go_gRPC/foo/bar/chat.proto` would exist, `protoc --go_out=server --go_opt=paths=source_relative foo/bar/chat.proto` would create `./go_gRPC/server/foo/bar/chat.pb.go`.
- **Example II:** For this project `protoc --go_out=. --go_opt=paths=source_relative chat.proto` would create `chat.pb.proto` in the root directory `./go_gRPC`

```
protoc --go_opt=paths=source_relative
```

</br>
