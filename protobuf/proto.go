package protobuf

//go:generate go install -v google.golang.org/protobuf/cmd/protoc-gen-go@latest
//go:generate go install -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
// should locate in the root of projects.
//go:generate go run ./infra/vprotogen/
