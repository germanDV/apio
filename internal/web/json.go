package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	maxBodySize           = 1_048_576 // 1MB
	disallowUnknownFields = false
)

type Envelope map[string]any

// WriteJSON sends the response in JSON format,
// if an encoding error occurs, it sends an empty 500 response.
func WriteJSON(w http.ResponseWriter, data Envelope, status int) {
	encoded, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(encoded)
}

// ReadJSON decodes r.Body into the destination pointed to by `dst` and handles possible errors.
func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Limit the size of the request body to prevent malicious requests
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBodySize))

	dec := json.NewDecoder(r.Body)

	// Disallow unknown fields, if an unknown field is found in the body, we return an error.
	// Without the DisallowUnknownFields, unknown fields would simply be ignored.
	if disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	// Decode request body to the provided destination
	err := dec.Decode(dst)

	// Handle errors
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBodySize)

		default:
			return err
		}
	}

	// Call Decode() again.
	// If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return an error.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
