package common

import "testing"

func TestNormalizeAndValidate(t *testing.T) {
    m := map[string]any{ "order_id":"o1", "symbol":"btc-usdt", "price":100.0, "size":1.0, "side":"buy", "ts":"2025-01-01T12:00:00Z", "platform":"mock", "seq":1 }
    e, err := Normalize(m)
    if err != nil { t.Fatalf("normalize err: %v", err) }
    if e.Symbol != "BTC-USDT" || e.Side != "BUY" { t.Fatalf("normalize upper failed: %v %v", e.Symbol, e.Side) }
    if err := Validate(e); err != nil { t.Fatalf("validate err: %v", err) }
}

