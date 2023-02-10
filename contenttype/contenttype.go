// Package contenttype exports MIME content-type constants.
// All text based types use charset UTF-8.
package contenttype

const (
	PlainText         = "text/plain; charset=utf-8"
	JavaScript        = "text/javascript; charset=utf-8"
	HTML              = "text/html; charset=utf-8"
	CSV               = "text/csv; charset=utf-8"
	XML               = "application/xml"
	JSON              = "application/json; charset=utf-8"
	PDF               = "application/pdf"
	Zip               = "application/zip"
	OctetStream       = "application/octet-stream"
	WWWFormURLEncoded = "application/x-www-form-urlencoded"
	MultipartFormData = "multipart/form-data"
	PNG               = "image/png"
	GIF               = "image/gif"
	JPEG              = "image/jpeg"
	TIFF              = "image/tiff"
)
