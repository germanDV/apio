package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadJSON(t *testing.T) {
	type TestStruct struct {
		Foo int
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		dst := TestStruct{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", strings.NewReader("{\"foo\": 42}"))

		err := ReadJSON(w, r, &dst)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if dst.Foo != 42 {
			t.Fatalf("expected 42, got %v", dst.Foo)
		}
	})

	t.Run("wrong field type", func(t *testing.T) {
		t.Parallel()
		dst := TestStruct{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", strings.NewReader("{\"foo\": \"NaN\"}"))

		err := ReadJSON(w, r, &dst)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		want := "body contains incorrect JSON type for field \"Foo\""
		got := err.Error()
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		dst := TestStruct{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)

		err := ReadJSON(w, r, &dst)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		want := "body must not be empty"
		got := err.Error()
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})

	t.Run("badly-formed", func(t *testing.T) {
		t.Parallel()
		dst := TestStruct{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", strings.NewReader("{\"foo\": \"bar"))

		err := ReadJSON(w, r, &dst)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		want := "body contains badly-formed JSON"
		got := err.Error()
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		t.Parallel()

		disallowUnknownFields = false
		dst := TestStruct{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", strings.NewReader("{\"bar\": true}"))

		err := ReadJSON(w, r, &dst)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		disallowUnknownFields = true
		dst = TestStruct{}
		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/", strings.NewReader("{\"bar\": true}"))

		err = ReadJSON(w, r, &dst)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		want := "body contains unknown key \"bar\""
		got := err.Error()
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})

	t.Run("extra data after json", func(t *testing.T) {
		t.Parallel()
		dst := TestStruct{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", strings.NewReader("{\"foo\": 42} wtf"))

		err := ReadJSON(w, r, &dst)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		want := "body must only contain a single JSON value"
		got := err.Error()
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("write a successful response", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		data := Envelope{"message": "foo"}
		WriteJSON(w, data, http.StatusOK)
		if w.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
		}
		if w.Body.String() != `{"message":"foo"}` {
			t.Fatalf("expected %q, got %q", `{"message":"foo"}`, w.Body.String())
		}
		if w.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("expected %q, got %q", "application/json", w.Header().Get("Content-Type"))
		}
	})

	t.Run("write an error response", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		data := Envelope{"error": "bar"}
		WriteJSON(w, data, http.StatusBadRequest)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
		if w.Body.String() != `{"error":"bar"}` {
			t.Fatalf("expected %q, got %q", `{"error":"bar"}`, w.Body.String())
		}
		if w.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("expected %q, got %q", "application/json", w.Header().Get("Content-Type"))
		}
	})
}
