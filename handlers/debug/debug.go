package debug

import (
	"encoding/json"
	"net/http"
	"os"

	"go.uber.org/zap"
)

// Struct with binded signature funcs
// Uses log and build params
type Handler struct {
	Log          *zap.SugaredLogger
	BuildVersion string
	BuildCommit  string
	BuildDate    string
}

// Heavy check if our main systems are live
func (h Handler) Readiness(rw http.ResponseWriter, r *http.Request) {
	code := http.StatusOK

	data := struct {
		Status string
	}{
		Status: "OK",
	}

	if err := process(rw, data, code); err != nil {
		h.Log.Errorf("[Readiness]", "[Process]", "[Failed]", err)
	}
	h.Log.Infof("[Readiness]", "[Code]", code, "[Method]", r.Method, "[Path]", r.URL.Path)
}

// Heavy check if our main systems are live
func (h Handler) Liveness(rw http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unavailable"
	}

	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Commit    string `json:"commit,omitempty"`
		Date      string `json:"date,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     h.BuildVersion,
		Commit:    h.BuildCommit,
		Date:      h.BuildDate,
		Host:      hostname,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	if err := process(rw, data, http.StatusOK); err != nil {
		h.Log.Errorf("[Liveness]", "[Process]", "[Failed]", err)
	}
	h.Log.Infof("[Liveness]", "[Code]", http.StatusOK, "[Method]", r.Method, "[Path]", r.URL.Path)
}

// process all necessery requests
func process(rw http.ResponseWriter, data interface{}, statusCode int) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rw.WriteHeader(statusCode)

	rw.Header().Add("Content-Type", "application/json")

	_, err = rw.Write(msg)
	if err != nil {
		return err
	}

	return nil
}
