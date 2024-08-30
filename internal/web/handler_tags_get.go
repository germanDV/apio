package web

import (
	"net/http"

	"github.com/a-h/rest"
)

type Tag struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	NoteCount int    `json:"note_count"`
}

type GetTagsRes struct {
	Tags []Tag `json:"tags"`
}

func (api *Api) handleGetTags() func(http.ResponseWriter, *http.Request) error {
	api.oas.Get("/tags").
		HasTags([]string{"Tags"}).
		HasOperationID("getTags").
		HasDescription("Get all Tags with note count.").
		HasResponseModel(http.StatusOK, rest.ModelOf[GetTagsRes]())

	return func(w http.ResponseWriter, r *http.Request) error {
		items, err := api.tagSvc.GetAll()
		if err != nil {
			return err
		}

		tags := make([]Tag, 0, len(items))
		for _, tag := range items {
			tags = append(tags, Tag{
				ID:        tag.ID.String(),
				Name:      tag.Name.String(),
				NoteCount: tag.NoteCount,
			})
		}

		WriteJSON(w, GetTagsRes{Tags: tags}, http.StatusOK)
		return nil
	}
}
