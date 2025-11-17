package push

import (
    "encoding/json"
    "net/http"
    "time"
    "crossx/internal/ingest/common"
)

type WSUpgrader interface {
    Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (WSConn, error)
}

type WSConn interface {
    ReadMessage() (int, []byte, error)
    WriteMessage(int, []byte) error
    SetReadDeadline(t time.Time) error
    SetWriteDeadline(t time.Time) error
    Close() error
}

type WSServer struct { Upgrader WSUpgrader }

func (s *WSServer) Handle(w http.ResponseWriter, r *http.Request) {
    conn, err := s.Upgrader.Upgrade(w, r, nil)
    if err != nil { http.Error(w, "upgrade", http.StatusBadRequest); return }
    defer conn.Close()
    _ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil { return }
        e, err := decodeEvent(msg)
        if err != nil { continue }
        _ = e
        _ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
        _ = conn.WriteMessage(1, []byte("ok"))
        _ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
    }
}

func decodeEvent(b []byte) (common.OrderEvent, error) { var e common.OrderEvent; err := json.Unmarshal(b, &e); if err != nil { return e, err }; if err := common.Validate(e); err != nil { return e, err }; return e, nil }
