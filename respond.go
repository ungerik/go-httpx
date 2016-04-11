package httpx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

// RespondJSON marshals response as JSON, sets "application/json" as content type
// and writes the marshalled JSON to writer
func RespondJSON(response interface{}, writer http.ResponseWriter) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	return err
}

// RespondPrettyJSON marshals response as indented JSON, sets "application/json" as content type
// and writes the marshalled JSON to writer
func RespondPrettyJSON(response interface{}, indent string, writer http.ResponseWriter) error {
	data, err := json.MarshalIndent(response, "", indent)
	if err != nil {
		return err
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	return err
}

// RespondXML marshals response as XML, sets "application/json" as content type
// and writes the XML standard header and the marshalled XML to writer.
// If rootElement is not empty, then an additional root element with this name will be wrapped around the content.
func RespondXML(response interface{}, rootElement string, writer http.ResponseWriter) error {
	data, err := xml.Marshal(response)
	if err != nil {
		return err
	}
	writer.Header().Set("Content-Type", "application/xml")
	if rootElement == "" {
		_, err = fmt.Fprintf(writer, "%s%s", xml.Header, data)
	} else {
		_, err = fmt.Fprintf(writer, "%s<%s>%s</%s>", xml.Header, rootElement, data, rootElement)
	}
	return err
}

// RespondPrettyXML marshals response as indented XML, sets "application/xml" as content type
// and writes the XML standard header and the marshalled XML to writer.
// If rootElement is not empty, then an additional root element with this name will be wrapped around the content.
func RespondPrettyXML(response interface{}, rootElement string, indent string, writer http.ResponseWriter) error {
	data, err := xml.MarshalIndent(response, "", indent)
	if err != nil {
		return err
	}
	writer.Header().Set("Content-Type", "application/xml")
	if rootElement == "" {
		_, err = fmt.Fprintf(writer, "%s%s", xml.Header, data)
	} else {
		_, err = fmt.Fprintf(writer, "%s<%s>%s</%s>", xml.Header, rootElement, data, rootElement)
	}
	return err
}
