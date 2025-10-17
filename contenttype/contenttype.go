// Package contenttype provides constants for common MIME content types.
//
// All text-based types include charset UTF-8 for proper character encoding.
// These constants can be used when setting Content-Type headers in HTTP responses.
//
// Example usage:
//
//	w.Header().Set("Content-Type", contenttype.JSON)
//	w.Header().Set("Content-Type", contenttype.HTML)
package contenttype

const (
	// Text content types (all with charset=utf-8)
	PlainText  = "text/plain; charset=utf-8"       // Plain text content
	JavaScript = "text/javascript; charset=utf-8"  // JavaScript code
	HTML       = "text/html; charset=utf-8"        // HTML documents
	CSV        = "text/csv; charset=utf-8"         // Comma-separated values

	// Data serialization formats
	XML  = "application/xml"                // XML documents
	JSON = "application/json; charset=utf-8" // JSON data

	// Binary formats
	PDF         = "application/pdf"         // PDF documents
	Zip         = "application/zip"         // ZIP archives
	OctetStream = "application/octet-stream" // Generic binary data

	// Form data types
	WWWFormURLEncoded = "application/x-www-form-urlencoded" // URL-encoded form data
	MultipartFormData = "multipart/form-data"               // Multipart form data (file uploads)

	// Image formats
	PNG  = "image/png"  // PNG images
	GIF  = "image/gif"  // GIF images
	JPEG = "image/jpeg" // JPEG images
	TIFF = "image/tiff" // TIFF images
)
