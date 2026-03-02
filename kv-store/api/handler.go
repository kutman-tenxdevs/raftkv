package api

import (
	"encoding/json"
	"kv-store/store"
	"net/http"
	"strings"
)

// Handler holds a reference to the store and serves HTTP requests.
type Handler struct {
	Store *store.Store
}

// NewHandler creates a new Handler with the given store.
func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s}
}

// HandleKey dispatches GET, PUT, DELETE on /keys/{key}.
func (h *Handler) HandleKey(w http.ResponseWriter, r *http.Request) {
	// Parse key from URL path: /keys/{key}
	path := strings.TrimPrefix(r.URL.Path, "/keys/")
	key := strings.TrimRight(path, "/")

	if key == "" {
		http.Error(w, `{"error":"key is required"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getKey(w, key)
	case http.MethodPut:
		h.putKey(w, r, key)
	case http.MethodDelete:
		h.deleteKey(w, key)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	}
}

// HandleKeys handles GET /keys — returns all key-value pairs.
func (h *Handler) HandleKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Store.GetAll())
}

func (h *Handler) getKey(w http.ResponseWriter, key string) {
	val, ok := h.Store.Get(key)
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "key not found"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"key": key, "value": val})
}

func (h *Handler) putKey(w http.ResponseWriter, r *http.Request, key string) {
	var body struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body"})
		return
	}

	h.Store.Set(key, body.Value)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"key": key, "value": body.Value})
}

func (h *Handler) deleteKey(w http.ResponseWriter, key string) {
	ok := h.Store.Delete(key)
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "key not found"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
}
