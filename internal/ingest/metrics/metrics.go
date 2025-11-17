package metrics

import (
    "sync"
)

type key struct{ Mode, Platform string }

type store struct{
    mu sync.Mutex
    success map[key]int64
    errors map[key]int64
}

var s = &store{success: map[key]int64{}, errors: map[key]int64{}}

func IncSuccess(mode, platform string) { s.mu.Lock(); s.success[key{mode, platform}]++; s.mu.Unlock() }
func IncError(mode, platform string) { s.mu.Lock(); s.errors[key{mode, platform}]++; s.mu.Unlock() }

type Snapshot struct{
    Success map[string]int64 `json:"success"`
    Errors map[string]int64 `json:"errors"`
}

func SnapshotJSON() Snapshot {
    s.mu.Lock(); defer s.mu.Unlock()
    succ := map[string]int64{}
    errs := map[string]int64{}
    for k, v := range s.success { succ[k.Mode+"|"+k.Platform] = v }
    for k, v := range s.errors { errs[k.Mode+"|"+k.Platform] = v }
    return Snapshot{Success: succ, Errors: errs}
}

