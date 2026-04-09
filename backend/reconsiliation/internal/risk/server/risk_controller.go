package server

import (
	"encoding/json"
	"net/http"

	"github.com/rafaeldepontes/reconsiliation/internal/risk"
	rs "github.com/rafaeldepontes/reconsiliation/internal/risk/service"
	"github.com/rafaeldepontes/reconsiliation/internal/util"
)

type controller struct {
	service risk.Service
}

func NewController() risk.Controller {
	return &controller{
		service: rs.NewService(),
	}
}

func (c *controller) GetRiskEvaluation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("transaction_id")
	res, err := c.service.GetResult(id)
	if err != nil {
		util.HandleError(w, "risk evaluation not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
