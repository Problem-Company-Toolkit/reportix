# Reportix

Reportix is a Golang package designed to help developers create meaningful error payloads and translate gRPC error messages to JSON via an interceptor. With Reportix, you can enhance error tracing and enable better error handling across different teams.

## Features

- **Error Interceptor**: Intercepts and processes errors returned by gRPC methods.
- **Error Payload Creation**: Constructs detailed error payloads containing reason, domain, code, message, and metadata.
- **Error Translation**: Converts gRPC errors to a JSON format, aiding frontend developers in error handling.
- **Callback Support**: Allows custom error processing through callback functions.

## Installation

```bash
go get github.com/problem-company-toolkit/reportix
```

## Usage

Here's a brief example of how you can use Reportix:

```go
import (
    "context"
    "fmt"
    "github.com/problem-company-toolkit/reportix"
    "google.golang.org/grpc/codes"
    "google.golang.org/genproto/googleapis/rpc/errdetails"
)

// Initialize error with details and metadata
err := reportix.NewError(
    codes.InvalidArgument,
    "Invalid UUID",
    &errdetails.ErrorInfo{
        Reason: "MISSING_FIELD",
        Domain: "example-service",
        Metadata: map[string]string{
            "field":  "id",
            "object": "user",
            "error":  "missing required field",
        },
    },
)

// Create a new error interceptor with a custom callback function
interceptor := reportix.NewErrInterceptor(reportix.ErrInterceptorOpts{
    Callback: func(ctx context.Context, grpcErr error) error {
        // Custom error processing
        return grpcErr
    },
})

// Use the UnaryServerInterceptor for gRPC method error processing
server := grpc.NewServer(
    grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()),
)

// ... Continue with server setup and registration of services
```