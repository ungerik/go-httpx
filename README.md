# go-httpx

Useful extensions around Go's `net/http` package providing simplified error handling, response writing, and graceful server shutdown.

[![Go Reference](https://pkg.go.dev/badge/github.com/ungerik/go-httpx.svg)](https://pkg.go.dev/github.com/ungerik/go-httpx)
[![Go Report Card](https://goreportcard.com/badge/github.com/ungerik/go-httpx)](https://goreportcard.com/report/github.com/ungerik/go-httpx)

## Installation

```bash
go get github.com/ungerik/go-httpx
```

## Features

- **Error Handling (`httperr`)**: Convert errors to HTTP responses with status codes
- **Response Writers (`respond`)**: Simplified response writing for JSON, XML, HTML, and plain text
- **Graceful Shutdown**: Handle server shutdown on OS signals
- **Content Types**: Constants for common MIME types
- **Panic Recovery**: Automatic panic handling in HTTP handlers

## Table of Contents

- [Quick Start](#quick-start)
- [Error Handling (httperr)](#error-handling-httperr)
  - [Basic Error Responses](#basic-error-responses)
  - [Custom Error Responses](#custom-error-responses)
  - [JSON Error Responses](#json-error-responses)
  - [HTTP Redirects as Errors](#http-redirects-as-errors)
  - [Error Logging](#error-logging)
  - [Sentinel Error Mapping](#sentinel-error-mapping)
- [Response Writers (respond)](#response-writers-respond)
  - [JSON Responses](#json-responses)
  - [HTML Responses](#html-responses)
  - [XML Responses](#xml-responses)
  - [Plain Text Responses](#plain-text-responses)
  - [Error-Only Handlers](#error-only-handlers)
- [Graceful Shutdown](#graceful-shutdown)
- [Content Types](#content-types)
- [Advanced Usage](#advanced-usage)

## Quick Start

### Basic Error Handling

```go
package main

import (
    "net/http"
    "github.com/ungerik/go-httpx/httperr"
    "github.com/ungerik/go-httpx/respond"
)

func main() {
    // Simple handler returning errors
    http.Handle("/api/user", respond.JSON(getUserHandler))

    http.ListenAndServe(":8080", nil)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) (any, error) {
    userID := r.URL.Query().Get("id")
    if userID == "" {
        return nil, httperr.BadRequest // Returns 400 Bad Request
    }

    user, err := fetchUser(userID)
    if err != nil {
        return nil, err // Errors are automatically handled
    }

    return user, nil // Automatically marshaled to JSON
}
```

## Error Handling (httperr)

The `httperr` package provides a clean way to return HTTP errors from handlers.

### Basic Error Responses

Pre-defined error responses for common HTTP status codes:

```go
import "github.com/ungerik/go-httpx/httperr"

func handler(w http.ResponseWriter, r *http.Request) error {
    // Use predefined errors
    return httperr.BadRequest       // 400
    return httperr.Unauthorized     // 401
    return httperr.PaymentRequired  // 402
    return httperr.Forbidden        // 403
    return httperr.NotFound         // 404
    return httperr.MethodNotAllowed // 405
}
```

### Custom Error Responses

Create custom error responses with specific status codes:

```go
// Simple status code
err := httperr.New(http.StatusTeapot) // 418 I'm a teapot

// With custom message
err := httperr.New(http.StatusBadRequest, "Invalid email format")

// With formatted message
err := httperr.Errorf(http.StatusBadRequest, "User %s not found", username)
```

### JSON Error Responses

Return structured error responses as JSON:

```go
type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Fields  []string `json:"fields,omitempty"`
}

func handler(w http.ResponseWriter, r *http.Request) error {
    detail := ErrorDetail{
        Code:    "VALIDATION_ERROR",
        Message: "Invalid input data",
        Fields:  []string{"email", "password"},
    }

    // Returns JSON error response with 400 status
    return httperr.JSON(http.StatusBadRequest, detail)
}
```

### HTTP Redirects as Errors

Use redirects as error values for control flow:

```go
func handler(w http.ResponseWriter, r *http.Request) error {
    if !isAuthenticated(r) {
        // Return redirect as error
        return httperr.TemporaryRedirect("/login")
    }

    // Or with specific status code
    return httperr.Redirect(http.StatusMovedPermanently, "/new-location")
    return httperr.PermanentRedirect("/new-location")
}
```

### Error Logging

Control which errors should be logged:

```go
import "github.com/ungerik/go-httpx/httperr"

// Check if error should be logged
if httperr.ShouldLog(err) {
    log.Printf("Error: %v", err)
}

// Wrap error to prevent logging (e.g., for expected errors like 404)
err := httperr.DontLog(httperr.NotFound)
httperr.ShouldLog(err) // Returns false
```

### Sentinel Error Mapping

Map standard errors to HTTP responses:

```go
import (
    "database/sql"
    "os"
    "github.com/ungerik/go-httpx/httperr"
)

func init() {
    // Default mappings (already configured):
    // os.ErrNotExist -> 404 Not Found
    // sql.ErrNoRows  -> 404 Not Found

    // Add custom sentinel mappings
    httperr.SentinelHandlers[sql.ErrConnDone] = http.HandlerFunc(
        func(w http.ResponseWriter, _ *http.Request) {
            http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
        },
    )
}

func handler(w http.ResponseWriter, r *http.Request) error {
    _, err := os.Open("nonexistent.txt")
    // Returns 404 automatically due to os.ErrNotExist mapping
    return err
}
```

### Error Handler Configuration

Customize the default error handler:

```go
import "github.com/ungerik/go-httpx/httperr"

func init() {
    // Show detailed errors in development
    httperr.DebugShowInternalErrorsInResponse = true

    // Customize error format
    httperr.DebugShowInternalErrorsInResponseFormat = "\n\nError: %+v"

    // Use custom error handler
    httperr.DefaultHandler = httperr.HandlerFunc(customErrorHandler)
}

func customErrorHandler(err error, w http.ResponseWriter, r *http.Request) bool {
    // Custom error handling logic
    log.Printf("Error on %s: %v", r.URL.Path, err)

    // Return false to let default handler process it
    return false
}
```

### Manual Error Handling

Handle errors manually without response writers:

```go
import "github.com/ungerik/go-httpx/httperr"

func handler(w http.ResponseWriter, r *http.Request) {
    err := processRequest(r)

    // Handle error (writes response if error is not nil)
    if httperr.Handle(err, w, r) {
        return // Error was handled
    }

    // Continue with normal response
    w.Write([]byte("Success"))
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if recovered := recover(); recovered != nil {
            // Handle panic as error
            httperr.HandlePanic(recovered, w, r)
        }
    }()

    // Code that might panic
    riskyOperation()
}
```

## Response Writers (respond)

The `respond` package provides type-safe handler functions that automatically serialize responses.

### JSON Responses

```go
import "github.com/ungerik/go-httpx/respond"

// Handler returns data to be marshaled as JSON
http.Handle("/api/users", respond.JSON(listUsers))

func listUsers(w http.ResponseWriter, r *http.Request) (any, error) {
    users, err := fetchUsers()
    if err != nil {
        return nil, err // Error automatically handled
    }

    // Automatically marshaled to JSON with Content-Type: application/json
    return users, nil
}

// Direct JSON writing
func customHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]any{"status": "ok", "count": 42}
    respond.WriteJSON(w, data)
}

// Configure JSON formatting
func init() {
    respond.PrettyPrint = true           // Enable pretty printing
    respond.PrettyPrintIndent = "  "     // Set indentation
}
```

### HTML Responses

```go
import "github.com/ungerik/go-httpx/respond"

// Dynamic HTML handler
http.Handle("/page", respond.HTML(renderPage))

func renderPage(w http.ResponseWriter, r *http.Request) ([]byte, error) {
    html := []byte("<html><body><h1>Hello World</h1></body></html>")
    return html, nil
}

// Static HTML handler
http.Handle("/static", respond.StaticHTML(`
    <!DOCTYPE html>
    <html>
        <body>
            <h1>Static Page</h1>
        </body>
    </html>
`))

// Direct HTML writing
func customHandler(w http.ResponseWriter, r *http.Request) {
    html := []byte("<h1>Custom</h1>")
    respond.WriteHTML(w, html)
}
```

### XML Responses

```go
import "github.com/ungerik/go-httpx/respond"

type User struct {
    XMLName xml.Name `xml:"user"`
    ID      int      `xml:"id"`
    Name    string   `xml:"name"`
}

// Handler returns data to be marshaled as XML
http.Handle("/api/user.xml", respond.XML(getUser))

func getUser(w http.ResponseWriter, r *http.Request) (any, error) {
    user := User{ID: 1, Name: "John Doe"}

    // Automatically marshaled to XML with Content-Type: application/xml
    return user, nil
}

// Direct XML writing
func customHandler(w http.ResponseWriter, r *http.Request) {
    user := User{ID: 1, Name: "Jane"}
    respond.WriteXML(w, user)
}
```

### Plain Text Responses

```go
import "github.com/ungerik/go-httpx/respond"

// Handler returns plain text string
http.Handle("/status", respond.Plaintext(getStatus))

func getStatus(w http.ResponseWriter, r *http.Request) (string, error) {
    return "System operational", nil
}

// Static plain text handler
http.Handle("/version", respond.StaticPlaintext("v1.2.3"))

// Direct plain text writing
func customHandler(w http.ResponseWriter, r *http.Request) {
    respond.WritePlaintext(w, "Hello World")
}
```

### Error-Only Handlers

For handlers that don't return data:

```go
import "github.com/ungerik/go-httpx/respond"

// Handler only returns error
http.Handle("/action", respond.Error(performAction))

func performAction(w http.ResponseWriter, r *http.Request) error {
    err := doSomething()
    if err != nil {
        return httperr.Errorf(http.StatusBadRequest, "Action failed: %v", err)
    }

    // Write success response manually
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Success"))

    return nil
}
```

### Panic Recovery

All `respond` handlers automatically recover from panics:

```go
import "github.com/ungerik/go-httpx/respond"

func init() {
    // Enable panic recovery (enabled by default)
    respond.CatchPanics = true
}

func riskyHandler(w http.ResponseWriter, r *http.Request) (any, error) {
    // If this panics, it will be caught and converted to 500 error
    riskyOperation()

    return map[string]string{"status": "ok"}, nil
}
```

## Graceful Shutdown

Handle graceful server shutdown on OS signals:

```go
package main

import (
    "log"
    "net/http"
    "os"
    "syscall"
    "time"

    "github.com/ungerik/go-httpx"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handler)

    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Setup graceful shutdown
    logger := log.New(os.Stdout, "", log.LstdFlags)
    httpx.GracefulShutdownServerOnSignal(
        server,
        logger,      // Signal logger (or nil)
        logger,      // Error logger (or nil)
        30*time.Second, // Shutdown timeout
        syscall.SIGINT, syscall.SIGTERM, // Signals (optional, defaults to SIGHUP, SIGINT, SIGTERM)
    )

    log.Println("Server starting on :8080")
    if err := server.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("Server error: %v", err)
    }
    log.Println("Server stopped")
}

func handler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World"))
}
```

Features:
- Listens for OS signals (default: SIGHUP, SIGINT, SIGTERM)
- Gracefully shuts down server with configurable timeout
- Optional logging of signals and errors
- Zero timeout disables timeout (waits indefinitely)

## Content Types

Constants for common MIME content types:

```go
import (
    "net/http"
    "github.com/ungerik/go-httpx/contenttype"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Text types (all with charset=utf-8)
    w.Header().Set("Content-Type", contenttype.PlainText)
    w.Header().Set("Content-Type", contenttype.HTML)
    w.Header().Set("Content-Type", contenttype.JavaScript)
    w.Header().Set("Content-Type", contenttype.CSV)

    // Data formats
    w.Header().Set("Content-Type", contenttype.JSON)
    w.Header().Set("Content-Type", contenttype.XML)

    // Binary formats
    w.Header().Set("Content-Type", contenttype.PDF)
    w.Header().Set("Content-Type", contenttype.Zip)
    w.Header().Set("Content-Type", contenttype.OctetStream)

    // Form data
    w.Header().Set("Content-Type", contenttype.WWWFormURLEncoded)
    w.Header().Set("Content-Type", contenttype.MultipartFormData)

    // Images
    w.Header().Set("Content-Type", contenttype.PNG)
    w.Header().Set("Content-Type", contenttype.GIF)
    w.Header().Set("Content-Type", contenttype.JPEG)
    w.Header().Set("Content-Type", contenttype.TIFF)
}
```

## Advanced Usage

### Custom Response Types

Implement the `httperr.Response` interface for custom error responses:

```go
import (
    "net/http"
    "github.com/ungerik/go-httpx/httperr"
)

type CustomError struct {
    StatusCode int
    Message    string
    Details    map[string]any
}

// Implement error interface
func (e CustomError) Error() string {
    return e.Message
}

// Implement http.Handler interface
func (e CustomError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(e.StatusCode)

    // Write custom JSON response
    json.NewEncoder(w).Encode(map[string]any{
        "error":   e.Message,
        "details": e.Details,
    })
}

// Now it can be returned as an error
func handler(w http.ResponseWriter, r *http.Request) error {
    return CustomError{
        StatusCode: http.StatusBadRequest,
        Message:    "Validation failed",
        Details:    map[string]any{"field": "email"},
    }
}
```

### Composing Error Handlers

Chain multiple error handlers:

```go
import "github.com/ungerik/go-httpx/httperr"

func handler(w http.ResponseWriter, r *http.Request) {
    err := processRequest()

    // Try multiple handlers in sequence
    handled := httperr.ForEachHandler(err, w, r,
        customHandler1,
        customHandler2,
        httperr.DefaultHandler,
    )

    if !handled {
        // Fallback handling
        http.Error(w, "Unhandled error", http.StatusInternalServerError)
    }
}

var customHandler1 = httperr.HandlerFunc(func(err error, w http.ResponseWriter, r *http.Request) bool {
    // Handle specific error types
    if errors.Is(err, ErrRateLimited) {
        http.Error(w, "Rate limited", http.StatusTooManyRequests)
        return true
    }
    return false
})
```

### Converting Panics to Errors

```go
import "github.com/ungerik/go-httpx/httperr"

func handler(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if recovered := recover(); recovered != nil {
            // Convert panic to error
            err := httperr.AsError(recovered)

            // Log the error
            log.Printf("Panic recovered: %v", err)

            // Handle as error
            httperr.Handle(err, w, r)
        }
    }()

    // Code that might panic
    riskyOperation()
}

// AsError converts various types to errors:
// - nil -> nil
// - error -> error (unchanged)
// - string -> errors.New(string)
// - fmt.Stringer -> errors.New(x.String())
// - other -> fmt.Errorf("%+v", val)
```

### Logger Interface

Implement the `httpx.Logger` interface for custom logging:

```go
import "github.com/ungerik/go-httpx"

type CustomLogger struct {
    // Your logger fields
}

func (l *CustomLogger) Printf(format string, args ...any) {
    // Your logging implementation
    log.Printf(format, args...)
}

// Use with graceful shutdown
var logger httpx.Logger = &CustomLogger{}
httpx.GracefulShutdownServerOnSignal(server, logger, logger, 30*time.Second)
```

## Best Practices

### 1. Use Appropriate Error Types

```go
// ✅ Good: Use specific error types
if userID == "" {
    return nil, httperr.BadRequest
}

if !hasPermission {
    return nil, httperr.Forbidden
}

// ❌ Avoid: Generic errors
return nil, errors.New("bad request")
```

### 2. Leverage DontLog for Expected Errors

```go
// ✅ Good: Don't log expected 404s
func getUser(id string) error {
    user, err := db.FindUser(id)
    if err == sql.ErrNoRows {
        return httperr.DontLog(httperr.NotFound)
    }
    return err
}
```

### 3. Use Structured Error Responses

```go
// ✅ Good: Structured errors for APIs
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

return httperr.JSON(http.StatusBadRequest, APIError{
    Code:    "INVALID_EMAIL",
    Message: "Email format is invalid",
})

// ❌ Avoid: Plain text errors in JSON APIs
return httperr.New(http.StatusBadRequest, "invalid email")
```

### 4. Configure Debug Mode Appropriately

```go
// Development
if os.Getenv("ENV") == "development" {
    httperr.DebugShowInternalErrorsInResponse = true
}

// Production
if os.Getenv("ENV") == "production" {
    httperr.DebugShowInternalErrorsInResponse = false
}
```

### 5. Use Response Handlers for Clean Code

```go
// ✅ Good: Clean separation of concerns
http.Handle("/api/users", respond.JSON(func(w http.ResponseWriter, r *http.Request) (any, error) {
    users, err := service.GetUsers(r.Context())
    return users, err
}))

// ❌ Avoid: Manual serialization
http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
    users, err := service.GetUsers(r.Context())
    if err != nil {
        // Manual error handling...
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
})
```

## Examples

See the `/examples` directory (if available) for complete working examples.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See [LICENSE](LICENSE) file for details.

## Related Projects

- [go-fs](https://github.com/ungerik/go-fs) - File system abstraction
- [go-dry](https://github.com/ungerik/go-dry) - DRY utilities for Go

## Support

For bugs and feature requests, please open an issue on GitHub.
