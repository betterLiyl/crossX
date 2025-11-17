package common

import (
	"errors"
	"strings"
)

func Normalize(m map[string]any) (OrderEvent, error) {
	var e OrderEvent
	if v, ok := m["order_id"].(string); ok {
		e.OrderID = v
	} else {
		return e, errors.New("order_id")
	}
	if v, ok := m["symbol"].(string); ok {
		e.Symbol = strings.ToUpper(v)
	} else {
		return e, errors.New("symbol")
	}
	switch v := m["price"].(type) {
	case float64:
		e.Price = v
	case int64:
		e.Price = float64(v)
	case int:
		e.Price = float64(v)
	default:
		return e, errors.New("price")
	}
	switch v := m["size"].(type) {
	case float64:
		e.Size = v
	case int64:
		e.Size = float64(v)
	case int:
		e.Size = float64(v)
	default:
		return e, errors.New("size")
	}
	if v, ok := m["side"].(string); ok {
		e.Side = strings.ToUpper(v)
	} else {
		return e, errors.New("side")
	}
	if v, ok := m["ts"].(string); ok {
		e.TS = v
	} else {
		return e, errors.New("ts")
	}
	if v, ok := m["platform"].(string); ok {
		e.Platform = v
	} else {
		return e, errors.New("platform")
	}
	if v, ok := m["signature"].(string); ok {
		e.Signature = v
	}
	switch v := m["seq"].(type) {
	case float64:
		e.Seq = int64(v)
	case int64:
		e.Seq = v
	case int:
		e.Seq = int64(v)
	}
	return e, nil
}
