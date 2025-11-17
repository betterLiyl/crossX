package common

import (
	"errors"
	"time"
)

func Validate(e OrderEvent) error {
	if e.OrderID == "" || e.Symbol == "" || e.Side == "" || e.Platform == "" || e.TS == "" {
		return errors.New("missing_fields")
	}
	if e.Price <= 0 || e.Size <= 0 {
		return errors.New("bad_values")
	}
	if e.Side != "BUY" && e.Side != "SELL" {
		return errors.New("bad_side")
	}
	if _, err := time.Parse(time.RFC3339Nano, e.TS); err != nil {
		return errors.New("bad_ts")
	}
	return nil
}
