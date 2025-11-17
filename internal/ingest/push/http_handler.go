package push

import (
	"crossx/internal/ingest/common"
	"encoding/json"
	"net/http"
	"os"
)

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	var m map[string]any
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "bad", http.StatusBadRequest)
		return
	}
	sig := r.Header.Get("X-Signature")
	secret := os.Getenv("PUSH_SECRET")
	if err := common.VerifyHMACSHA256(m, secret, sig); err != nil {
		http.Error(w, "sig", http.StatusUnauthorized)
		return
	}
	e, err := common.Normalize(m)
	if err != nil {
		http.Error(w, "bad", http.StatusBadRequest)
		return
	}
	if err := common.Validate(e); err != nil {
		http.Error(w, "bad", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
