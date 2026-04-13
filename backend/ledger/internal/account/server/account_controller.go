package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rafaeldepontes/ledger/internal/account"
	"github.com/rafaeldepontes/ledger/internal/util"
)

type controller struct {
	service account.Service
}

func NewController(svc account.Service) account.Controller {
	return &controller{
		service: svc,
	}
}

// GetAccountBalance godoc
// @Summary Get account balance
// @Description Returns the computed balance for a specific account
// @Tags accounts
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} model.BalanceRes
// @Failure 400 {object} map[string]string
// @Router /accounts/{id}/balance [get]
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
