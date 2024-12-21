package web

import (
	"net/http"
	"time"

	"github.com/a-h/rest"
	"github.com/germandv/apio/internal/errs"
)

type CreateNoteReq struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	TagIDs  []string `json:"tag_ids"`
}

type CreateNoteRes struct {
	ID string `json:"id"`
}

func (api *API) handleCreateNote() func(http.ResponseWriter, *http.Request) error {
	api.oas.Post("/notes").
		HasTags([]string{"Notes"}).
		HasOperationID("createNote").
		HasDescription("Create a new Note.").
		HasResponseModel(http.StatusCreated, rest.ModelOf[CreateNoteRes]()).
		HasResponseModel(http.StatusBadRequest, rest.ModelOf[errs.ErrResp]()).
		HasRequestModel(rest.ModelOf[CreateNoteReq]())

	return func(w http.ResponseWriter, r *http.Request) error {
		user, err := GetUser(r.Context())
		if err != nil {
			return err
		}

		body := CreateNoteReq{}
		err = ReadJSON(w, r, &body)
		if err != nil {
			return err
		}

		id, err := api.noteSvc.Create(body.Title, body.Content, body.TagIDs, user.ID, time.Now().UTC())
		if err != nil {
			return err
		}

		WriteJSON(w, CreateNoteRes{ID: id.String()}, http.StatusCreated)
		return nil
	}
}
