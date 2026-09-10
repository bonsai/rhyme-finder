package rhymfinder

import (
    "encoding/json"
    "net/http"
)

type Server struct { Finder *Finder }

func (s Server) Handler() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, http.StatusOK, map[string]string{"status":"ok"}) })
    mux.HandleFunc("POST /v1/rhyme", func(w http.ResponseWriter, r *http.Request) {
        var q Query
        if err := json.NewDecoder(r.Body).Decode(&q); err != nil { writeJSON(w,http.StatusBadRequest,map[string]string{"error":err.Error()}); return }
        if q.Mode == "" { q.Mode = ModeAggressive }
        if s.Finder == nil { writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"finder is nil"}); return }
        writeJSON(w,http.StatusOK,map[string]any{"query":q,"candidates":s.Finder.Find(q)})
    })
    return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) { w.Header().Set("Content-Type","application/json"); w.WriteHeader(status); _=json.NewEncoder(w).Encode(v) }
