package render

import (
	"encoding/json"
	"net/http"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

// JSONRender is the interface for rendering Json
type JSONRender interface {
	Render(http.ResponseWriter) error
	WriteContentType(w http.ResponseWriter)
}

// Json contains the given interface object.
type Json struct {
	Data interface{}
}

func (r Json) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (r Json) Render(w http.ResponseWriter) error {
	return WriteJSON(w, r.Data)
}

func WriteJSON(w http.ResponseWriter, obj any) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
