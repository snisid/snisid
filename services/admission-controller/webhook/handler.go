package webhook

import (
	"encoding/json"
	"net/http"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type AdmissionReview struct {
	Request  *AdmissionRequest  `json:"request"`
	Response *AdmissionResponse `json:"response"`
}

type AdmissionRequest struct {
	UID  string          `json:"uid"`
	Kind map[string]string `json:"kind"`
}

type AdmissionResponse struct {
	UID     string `json:"uid"`
	Allowed bool   `json:"allowed"`
	Result  *Result `json:"status,omitempty"`
}

type Result struct {
	Message string `json:"message"`
}

func HandleAdmission(w http.ResponseWriter, r *http.Request) {
	var review AdmissionReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	logger.Info("ADMISSION: Validating incoming request via AI Policy Engine...")

	// Default Allow logic
	review.Response = &AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: true,
	}

	// Policy check example
	if review.Request.Kind["kind"] == "Pod" {
		// Mock AI risk check
		riskScore := 0.1
		if riskScore > 0.8 {
			review.Response.Allowed = false
			review.Response.Result = &Result{Message: "Suspicious pod deployment blocked by SNISID AI"}
		}
	}

	json.NewEncoder(w).Encode(review)
}
