package web

import (
	"net/http"
	"time"

	"github.com/a-h/rest"
	"github.com/germandv/apio/internal/notes"
	"github.com/germandv/apio/internal/tools"
)

type NoteTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []NoteTag `json:"tags"`
}

type GetNotesRes struct {
	Notes []Note `json:"notes"`
}

func (api *Api) handleGetNotes() func(http.ResponseWriter, *http.Request) error {
	api.oas.Get("/notes").
		HasTags([]string{"Notes"}).
		HasOperationID("getNotes").
		HasDescription("Get all Notes").
		HasResponseModel(http.StatusOK, rest.ModelOf[GetNotesRes]())

	return func(w http.ResponseWriter, r *http.Request) error {
		items, err := api.noteSvc.GetAll()
		if err != nil {
			return err
		}

		entries := make([]Note, 0, len(items))
		for _, n := range items {
			tags := tools.Map(n.Tags, func(t notes.NoteTagEntity) NoteTag {
				return NoteTag{ID: t.ID.String(), Name: t.Name}
			})

			entries = append(entries, Note{
				ID:        n.ID.String(),
				Title:     n.Title.String(),
				Content:   n.Content.String(),
				CreatedBy: n.CreatedBy.String(),
				CreatedAt: n.CreatedAt,
				UpdatedAt: n.UpdatedAt,
				Tags:      tags,
			})
		}

		WriteJSON(w, GetNotesRes{Notes: entries}, http.StatusOK)
		return nil
	}
}
