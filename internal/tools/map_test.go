package tools

import (
	"fmt"
	"slices"
	"testing"
)

func TestMap(t *testing.T) {
	type From struct {
		FirstName string
		LastName  string
	}

	type To struct {
		Name string
	}

	from := []From{
		{FirstName: "Jane", LastName: "Doe"},
		{FirstName: "John", LastName: "Doe"},
	}

	to := Map(from, func(e From) To {
		return To{Name: fmt.Sprintf("%s %s", e.FirstName, e.LastName)}
	})

	if len(to) != len(from) {
		t.Fatalf("expected %d elements, got %d", len(from), len(to))
	}

	if !slices.ContainsFunc(to, func(e To) bool {
		return e.Name == "Jane Doe"
	}) {
		t.Fatalf("Jane Doe not found, got %v", to)
	}

	if !slices.ContainsFunc(to, func(e To) bool {
		return e.Name == "John Doe"
	}) {
		t.Fatalf("John Doe not found, got %v", to)
	}
}
