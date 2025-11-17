package common

type OrderEvent struct {
    OrderID  string  `json:"order_id"`
    Symbol   string  `json:"symbol"`
    Price    float64 `json:"price"`
    Size     float64 `json:"size"`
    Side     string  `json:"side"`
    TS       string  `json:"ts"`
    Platform string  `json:"platform"`
    Signature string `json:"signature"`
    Seq      int64   `json:"seq"`
}

