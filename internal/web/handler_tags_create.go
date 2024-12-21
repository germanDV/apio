package web

import (
	"net/http"

	"github.com/a-h/rest"
	"github.com/germandv/apio/internal/errs"
)

type CreateTagReq struct {
	Name string `json:"name"`
}

type CreateTagRes struct {
	ID string `json:"id"`
}

func (api *API) handleCreateTag() func(http.ResponseWriter, *http.Request) error {
	api.oas.Post("/tags").
		HasTags([]string{"Tags"}).
		HasOperationID("createTag").
		HasDescription("Create a new Tag.").
		HasResponseModel(http.StatusCreated, rest.ModelOf[CreateTagRes]()).
		HasResponseModel(http.StatusBadRequest, rest.ModelOf[errs.ErrResp]()).
		HasResponseModel(http.StatusConflict, rest.ModelOf[errs.ErrResp]()).
		HasRequestModel(rest.ModelOf[CreateTagReq]())

	return func(w http.ResponseWriter, r *http.Request) error {
		body := CreateTagReq{}
		err := ReadJSON(w, r, &body)
		if err != nil {
			return err
		}

		id, err := api.tagSvc.Create(body.Name)
		if err != nil {
			return err
		}

		WriteJSON(w, CreateTagRes{ID: id.String()}, http.StatusCreated)
		return nil
	}
}
