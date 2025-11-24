// Package engine implements a minimal in-memory order matching engine
// with per-symbol order books and immediate (continuous) matching.
package engine

import (
    "sort"
    "sync"
    "time"
)

// Side denotes whether an order is a buy or sell.
type Side string

// Supported order sides.
const (
    Buy  Side = "BUY"
    Sell Side = "SELL"
)

// PlaceOrderRequest is the input payload to place an order
// into the matching engine for a specific symbol.
type PlaceOrderRequest struct {
    Symbol string  `json:"symbol"`
    Side   Side    `json:"side"`
    Price  float64 `json:"price"`
    Size   float64 `json:"size"`
}

// OrderStatus represents the current state of an order
// tracked by the engine, including partial fills.
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

// book holds the buy and sell sides for a single symbol.
// Orders are kept in slices and sorted by price priority.
type book struct {
    buys  []OrderStatus
    sells []OrderStatus
}

// Engine maintains per-symbol books and a flat index of orders by ID.
// It is safe for concurrent access via internal locking.
type Engine struct {
    mu     sync.RWMutex
    books  map[string]*book
    orders map[string]OrderStatus
}

// New constructs an empty matching engine instance.
func New() *Engine {
    return &Engine{books: make(map[string]*book), orders: make(map[string]OrderStatus)}
}

// PlaceOrder inserts the order into the appropriate side of the book,
// re-sorts by price priority, performs immediate matching, and returns
// the up-to-date status of the newly placed order.
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

// GetOrderStatus retrieves an order by ID from the engine.
func (e *Engine) GetOrderStatus(id string) (OrderStatus, bool) {
    e.mu.RLock()
    defer e.mu.RUnlock()
    os, ok := e.orders[id]
    return os, ok
}

// AggregateQuote summarizes the current best bid/ask for a symbol.
type AggregateQuote struct {
    Symbol  string  `json:"symbol"`
    BestBid float64 `json:"best_bid"`
    BestAsk float64 `json:"best_ask"`
}

// Aggregate returns the current best bid and best ask snapshot for the symbol.
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

// indexOrders rebuilds the ID -> OrderStatus index for a book.
// This is an O(n) operation over the current slices.
func (e *Engine) indexOrders(b *book) {
    for _, o := range b.buys {
        e.orders[o.ID] = o
    }
    for _, o := range b.sells {
        e.orders[o.ID] = o
    }
}

// getBook returns the book for the symbol, creating it if absent.
func (e *Engine) getBook(symbol string) *book {
    if b, ok := e.books[symbol]; ok {
        return b
    }
    b := &book{buys: []OrderStatus{}, sells: []OrderStatus{}}
    e.books[symbol] = b
    return b
}

// match repeatedly crosses the top of the book while best bid >= best ask,
// updating filled quantities and statuses. Filled orders remain in-place.
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
            if buy.Filled >= buy.Size {
                i++
            }
            if sell.Filled >= sell.Size {
                j++
            }
            continue
        }
        buy.Filled += qty
        sell.Filled += qty
        b.buys[i] = buy
        b.sells[j] = sell
        if buy.Filled >= buy.Size {
            b.buys[i].Status = "FILLED"
            i++
        }
        if sell.Filled >= sell.Size {
            b.sells[j].Status = "FILLED"
            j++
        }
    }
}

// min returns the smaller of two float64 values.
func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

// seq provides a simple process-wide incrementing counter for ID suffixes.
var seq struct {
    mu sync.Mutex
    n  int64
}

// makeID builds a unique order ID using symbol, timestamp, and a counter.
func makeID(sym string) string {
    seq.mu.Lock()
    defer seq.mu.Unlock()
    seq.n++
    return sym + "-" + time.Now().Format("20060102150405") + "-" + fmtI(seq.n)
}

// fmtI formats a positive int64 into its base-10 string representation.
func fmtI(n int64) string {
    b := make([]byte, 0, 20)
    if n == 0 {
        return "0"
    }
    for n > 0 {
        b = append(b, byte('0'+n%10))
        n /= 10
    }
    for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
        b[i], b[j] = b[j], b[i]
    }
    return string(b)
}
