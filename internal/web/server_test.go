package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/a-h/rest"
	"github.com/germandv/apio/internal/cache/memorycache"
	"github.com/germandv/apio/internal/id"
	"github.com/germandv/apio/internal/memorydb"
	"github.com/germandv/apio/internal/notes"
	"github.com/germandv/apio/internal/tags"
	"github.com/germandv/apio/internal/tokenauth"
)

const (
	PrivateKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIALMOYNRFv2OqZltGQO0Gg34gQK8nb0W/aXzUS90Uve2oAoGCCqGSM49
AwEHoUQDQgAEUOu0Nc9/EiVSyBKyfvv38MlteRWA+6S8jpRIOC2eMn2kYSv1RCc7
uejvLVc0EYn2spObZjsMv4qvNz0XxYduDQ==
-----END EC PRIVATE KEY-----`

	PublicKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEUOu0Nc9/EiVSyBKyfvv38MlteRWA
+6S8jpRIOC2eMn2kYSv1RCc7uejvLVc0EYn2spObZjsMv4qvNz0XxYduDQ==
-----END PUBLIC KEY-----`
)

var api *Api

func TestMain(m *testing.M) {
	// Before tests.
	cacheClient, _ := memorycache.New()
	auth, _ := tokenauth.New(PrivateKey, PublicKey)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	oas := rest.NewAPI("test_api")
	tagSvc := tags.NewService(memorydb.NewTagsRepository())
	noteSvc := notes.NewService(memorydb.NewNotesRepository())
	api = New(logger, auth, oas, tagSvc, noteSvc, cacheClient)

	// Run tests.
	exitCode := m.Run()

	// After tests.
	os.Exit(exitCode)
}

func TestAPI(t *testing.T) {
	// Healthcheck.
	r := httptest.NewRequest("GET", "/healthcheck", nil)
	w := httptest.NewRecorder()
	handler := api.handleHealthcheck()
	err := handler(w, r)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}

	// Create a Tag.
	r = httptest.NewRequest("POST", "/tags", strings.NewReader(`{"name":"testtag"}`))
	w = httptest.NewRecorder()
	handler = api.handleCreateTag()

	err = handler(w, r)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
	}

	tagOneRes := CreateTagRes{}
	json.Unmarshal(w.Body.Bytes(), &tagOneRes)
	_, err = id.Parse(tagOneRes.ID)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	// Create a second Tag.
	r = httptest.NewRequest("POST", "/tags", strings.NewReader(`{"name":"zen"}`))
	w = httptest.NewRecorder()
	handler(w, r)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	tagTwoRes := CreateTagRes{}
	json.Unmarshal(w.Body.Bytes(), &tagTwoRes)
	_, err = id.Parse(tagTwoRes.ID)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	// Retrieve Tags.
	r = httptest.NewRequest("GET", "/tags", nil)
	w = httptest.NewRecorder()
	handler = api.handleGetTags()

	err = handler(w, r)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}

	getTagsRes := GetTagsRes{}
	json.Unmarshal(w.Body.Bytes(), &getTagsRes)
	if len(getTagsRes.Tags) != 2 {
		t.Fatalf("expected 2 tag, got %d", len(getTagsRes.Tags))
	}

	// Create Note without Tag.
	body := `{"title": "title_one", "content": "content_one"}`
	r = httptest.NewRequest("POST", "/notes", strings.NewReader(body))
	w = httptest.NewRecorder()
	handler = api.handleCreateNote()

	err = handler(w, r)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
	}

	noteWithoutTagsRes := CreateNoteRes{}
	json.Unmarshal(w.Body.Bytes(), &noteWithoutTagsRes)
	_, err = id.Parse(noteWithoutTagsRes.ID)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	// Create Note with 2 Tags.
	createNoteReq := CreateNoteReq{
		Title:   "title_two",
		Content: "content_two",
		TagIDs:  []string{tagOneRes.ID, tagTwoRes.ID},
	}

	jsonBody, err := json.Marshal(createNoteReq)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	r = httptest.NewRequest("POST", "/notes", bytes.NewBuffer(jsonBody))
	w = httptest.NewRecorder()
	handler = api.handleCreateNote()

	err = handler(w, r)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
	}

	noteWithTagsRes := CreateNoteRes{}
	json.Unmarshal(w.Body.Bytes(), &noteWithTagsRes)
	_, err = id.Parse(noteWithTagsRes.ID)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	// Retrieve Notes.
	r = httptest.NewRequest("GET", "/notes", nil)
	w = httptest.NewRecorder()
	handler = api.handleGetNotes()

	err = handler(w, r)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}

	getNotesRes := GetNotesRes{}
	json.Unmarshal(w.Body.Bytes(), &getNotesRes)
	if len(getNotesRes.Notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(getNotesRes.Notes))
	}

	for _, note := range getNotesRes.Notes {
		if note.Title == "title_one" {
			if len(note.Tags) != 0 {
				t.Fatalf("expected 0 tags, got %d", len(note.Tags))
			}
		}

		if note.Title == "title_two" {
			if len(note.Tags) != 2 {
				t.Fatalf("expected 2 tags, got %d", len(note.Tags))
			}
			for _, tag := range note.Tags {
				fmt.Println(tag.Name)
				if tag.Name != "testtag" && tag.Name != "zen" {
					t.Fatalf("expected tags testtag and zen, got %s", tag.Name)
				}
			}
		}
	}
}
