package pull

import (
    "context"
    "time"
)

type Adapter interface {
    Name() string
    FetchChanges(ctx context.Context, since Cursor) ([]Event, Cursor, error)
}

type Cursor struct{ Seq int64; Block int64 }

type Event struct{ Data map[string]any }

type Scheduler struct{
    Adapters []Adapter
    Interval time.Duration
}

func (s *Scheduler) Run(ctx context.Context) {
    if s.Interval <= 0 { s.Interval = 30 * time.Second }
    ticker := time.NewTicker(s.Interval)
    defer ticker.Stop()
    cursors := map[string]Cursor{}
    for {
        select{
        case <-ctx.Done(): return
        case <-ticker.C:
            for _, a := range s.Adapters {
                cur := cursors[a.Name()]
                evs, next, _ := a.FetchChanges(ctx, cur)
                _ = evs
                cursors[a.Name()] = next
            }
        }
    }
}

