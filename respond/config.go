// Package respond provides type-safe HTTP handler functions that automatically
// serialize responses to various formats (JSON, XML, HTML, plain text).
//
// The package simplifies response handling by allowing handlers to return
// data and errors, which are then automatically serialized and written
// with appropriate content types.
//
// Key features:
//   - Automatic response serialization (JSON, XML, HTML, plain text)
//   - Built-in panic recovery
//   - Automatic error handling via httperr
//   - Pretty-printing support for JSON and XML
//   - Type-safe handler function types
//
// Example usage:
//
//	http.Handle("/api/users", respond.JSON(func(w http.ResponseWriter, r *http.Request) (any, error) {
//	    users, err := fetchUsers()
//	    return users, err // Automatically marshaled to JSON
//	}))
package respond

var (
	// CatchPanics controls whether response handlers automatically recover from panics.
	// When true (default), panics are caught and converted to 500 Internal Server Error responses.
	// Set to false in testing to see panic stack traces.
	CatchPanics = true

	// PrettyPrint controls whether JSON and XML responses are formatted with indentation.
	// When true (default), responses are pretty-printed for better readability.
	// Set to false in production to reduce response size.
	PrettyPrint = true

	// PrettyPrintIndent is the string used for each indentation level when pretty-printing.
	// Default is two spaces ("  ").
	PrettyPrintIndent = "  "
)
