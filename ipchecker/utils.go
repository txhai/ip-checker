package ipchecker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func responseError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("error %v", err)})
	return
}
