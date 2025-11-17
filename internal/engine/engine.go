package engine

import (
    "sort"
    "sync"
    "time"
)

type Side string

const (
    Buy  Side = "BUY"
    Sell Side = "SELL"
)

type PlaceOrderRequest struct {
    Symbol string  `json:"symbol"`
    Side   Side    `json:"side"`
    Price  float64 `json:"price"`
    Size   float64 `json:"size"`
}

type OrderStatus struct {
    ID        string  `json:"id"`
    Symbol    string  `json:"symbol"`
    Side      Side    `json:"side"`
    Price     float64 `json:"price"`
    Size      float64 `json:"size"`
    Filled    float64 `json:"filled"`
    Status    string  `json:"status"`
    CreatedAt int64   `json:"created_at"`
}

type book struct {
    buys  []OrderStatus
    sells []OrderStatus
}

type Engine struct {
    mu    sync.RWMutex
    books map[string]*book
    orders map[string]OrderStatus
}

func New() *Engine {
    return &Engine{books: make(map[string]*book), orders: make(map[string]OrderStatus)}
}

func (e *Engine) PlaceOrder(req PlaceOrderRequest) OrderStatus {
    e.mu.Lock()
    defer e.mu.Unlock()
    b := e.getBook(req.Symbol)
    id := makeID(req.Symbol)
    os := OrderStatus{ID: id, Symbol: req.Symbol, Side: req.Side, Price: req.Price, Size: req.Size, Status: "OPEN", CreatedAt: time.Now().UnixMilli()}
    if req.Side == Buy {
        b.buys = append(b.buys, os)
        sort.Slice(b.buys, func(i, j int) bool { return b.buys[i].Price > b.buys[j].Price })
    } else {
        b.sells = append(b.sells, os)
        sort.Slice(b.sells, func(i, j int) bool { return b.sells[i].Price < b.sells[j].Price })
    }
    e.match(b)
    e.indexOrders(b)
    return e.orders[id]
}

func (e *Engine) GetOrderStatus(id string) (OrderStatus, bool) {
    e.mu.RLock()
    defer e.mu.RUnlock()
    os, ok := e.orders[id]
    return os, ok
}

type AggregateQuote struct {
    Symbol string  `json:"symbol"`
    BestBid float64 `json:"best_bid"`
    BestAsk float64 `json:"best_ask"`
}

func (e *Engine) Aggregate(symbol string) AggregateQuote {
    e.mu.RLock()
    defer e.mu.RUnlock()
    b := e.getBook(symbol)
    bid := 0.0
    ask := 0.0
    if len(b.buys) > 0 {
        bid = b.buys[0].Price
    }
    if len(b.sells) > 0 {
        ask = b.sells[0].Price
    }
    return AggregateQuote{Symbol: symbol, BestBid: bid, BestAsk: ask}
}

func (e *Engine) indexOrders(b *book) {
    for _, o := range b.buys {
        e.orders[o.ID] = o
    }
    for _, o := range b.sells {
        e.orders[o.ID] = o
    }
}

func (e *Engine) getBook(symbol string) *book {
    if b, ok := e.books[symbol]; ok {
        return b
    }
    b := &book{buys: []OrderStatus{}, sells: []OrderStatus{}}
    e.books[symbol] = b
    return b
}

func (e *Engine) match(b *book) {
    i := 0
    j := 0
    for i < len(b.buys) && j < len(b.sells) {
        buy := b.buys[i]
        sell := b.sells[j]
        if buy.Price < sell.Price {
            break
        }
        qty := min(buy.Size-buy.Filled, sell.Size-sell.Filled)
        if qty <= 0 {
            if buy.Filled >= buy.Size { i++ }
            if sell.Filled >= sell.Size { j++ }
            continue
        }
        buy.Filled += qty
        sell.Filled += qty
        b.buys[i] = buy
        b.sells[j] = sell
        if buy.Filled >= buy.Size { b.buys[i].Status = "FILLED"; i++ }
        if sell.Filled >= sell.Size { b.sells[j].Status = "FILLED"; j++ }
    }
}

func min(a, b float64) float64 { if a < b { return a }; return b }

var seq struct{ mu sync.Mutex; n int64 }

func makeID(sym string) string {
    seq.mu.Lock(); defer seq.mu.Unlock(); seq.n++
    return sym + "-" + time.Now().Format("20060102150405") + "-" + fmtI(seq.n)
}

func fmtI(n int64) string {
    b := make([]byte, 0, 20)
    if n == 0 { return "0" }
    for n > 0 { b = append(b, byte('0'+n%10)); n/=10 }
    for i, j := 0, len(b)-1; i<j; i, j = i+1, j-1 { b[i], b[j] = b[j], b[i] }
    return string(b)
}

