package risk

import "net/http"

type Controller interface {
	GetRiskEvaluation(w http.ResponseWriter, r *http.Request)
}
