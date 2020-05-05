package extractors

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func GetUUID(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	return uuid.FromString(vars["id"])
}
