package api

import (
	"auth/extractors"
	"auth/inout"
	"net/http"
)

func (api *API) GetSecretV1(w http.ResponseWriter, r *http.Request) {

	id, _ := extractors.GetUUID(r)
	secret := api.Controller.GetSecret(r.Context(), id)

	if secret == nil {
		api.Renderer.Render(w, r, http.StatusNotFound, nil)
	} else {
		api.Renderer.Render(w, r, http.StatusOK, &inout.GetSecretResponseV1{
			Data: &inout.Secret{
				Id:      secret.Id.Bytes(),
				Created: secret.Created,
				Value:   secret.Value.Bytes(),
			}})
	}
}
