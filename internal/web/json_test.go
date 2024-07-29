package web

import (
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

// TODO: TestWriteJSON
