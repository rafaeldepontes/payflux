package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rafaeldepontes/ledger/internal/account"
	as "github.com/rafaeldepontes/ledger/internal/account/service"
	"github.com/rafaeldepontes/ledger/internal/util"
)

type controller struct {
	service account.Service
}

func NewController() account.Controller {
	return &controller{
		service: as.NewService(),
	}
}

func (c controller) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		util.HandleError(w, "missing account id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.HandleError(w, "invalid account id", http.StatusBadRequest)
		return
	}

	res, err := c.service.GetAccountBalance(id)
	if err != nil {
		util.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
