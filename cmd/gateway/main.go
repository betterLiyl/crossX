package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"crossx/internal/adapters"
	"crossx/internal/engine"
	"crossx/internal/ingest/metrics"
	push "crossx/internal/ingest/push"
)

type server struct {
	mux    *http.ServeMux
	secret string
	eng    *engine.Engine
}

func main() {
	s := &server{mux: http.NewServeMux(), secret: os.Getenv("JWT_SECRET"), eng: engine.New()}
	s.routes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{Addr: ":" + port, Handler: s.authMiddleware(s.mux)}
	log.Fatal(srv.ListenAndServe())
}

func (s *server) routes() {
	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	s.mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/swagger.json")
	})
	s.mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, metrics.SnapshotJSON())
	})
    s.mux.HandleFunc("/v1/orders/push/http", func(w http.ResponseWriter, r *http.Request) {
        push.HTTPHandler(w, r)
    })
    ws := &push.WSServer{Upgrader: newDefaultUpgrader()}
    s.mux.HandleFunc("/v1/orders/push/ws", func(w http.ResponseWriter, r *http.Request) {
        ws.Handle(w, r)
    })
	s.mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		var req engine.PlaceOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		res := s.eng.PlaceOrder(req)
		writeJSON(w, http.StatusOK, res)
	})
	s.mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/orders/")
		if id == "" {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		res, ok := s.eng.GetOrderStatus(id)
		if !ok {
			http.Error(w, "notfound", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
	s.mux.HandleFunc("/aggregate/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		symbol := strings.TrimPrefix(r.URL.Path, "/aggregate/")
		if symbol == "" {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		agg := adapters.AggregateBest([]adapters.MarketDataAdapter{adapters.OKXAdapter{}, adapters.HypeliquldAdapter{}}, symbol, 800000000)
		writeJSON(w, http.StatusOK, agg)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" || r.URL.Path == "/swagger.json" {
			next.ServeHTTP(w, r)
			return
		}
		h := r.Header.Get("Authorization")
		if h == "" {
			http.Error(w, "auth", http.StatusUnauthorized)
			return
		}
		p := strings.SplitN(h, " ", 2)
		if len(p) != 2 || strings.ToLower(p[0]) != "bearer" {
			http.Error(w, "auth", http.StatusUnauthorized)
			return
		}
		if !verifyToken(s.secret, p[1]) {
			http.Error(w, "auth", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func verifyToken(secret, token string) bool {
	if secret == "" {
		return true
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(parts[0] + "." + parts[1]))
	sig := mac.Sum(nil)
	want := base64.RawURLEncoding.EncodeToString(sig)
	return hmac.Equal([]byte(want), []byte(parts[2]))
}
type defaultUpgrader struct{}

func newDefaultUpgrader() *defaultUpgrader { return &defaultUpgrader{} }

func (u *defaultUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (push.WSConn, error) {
    return &noopConn{}, nil
}

type noopConn struct{}

func (n *noopConn) ReadMessage() (int, []byte, error) { return 1, []byte("{}"), nil }
func (n *noopConn) WriteMessage(t int, b []byte) error { return nil }
func (n *noopConn) SetReadDeadline(t time.Time) error { return nil }
func (n *noopConn) SetWriteDeadline(t time.Time) error { return nil }
func (n *noopConn) Close() error { return nil }
