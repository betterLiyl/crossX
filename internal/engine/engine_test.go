package engine

import "testing"

func TestMatchFullFill(t *testing.T) {
    e := New()
    b := PlaceOrderRequest{Symbol: "BTC-USDT", Side: Buy, Price: 100, Size: 1}
    s := PlaceOrderRequest{Symbol: "BTC-USDT", Side: Sell, Price: 90, Size: 1}
    bo := e.PlaceOrder(b)
    so := e.PlaceOrder(s)
    bo2, _ := e.GetOrderStatus(bo.ID)
    so2, _ := e.GetOrderStatus(so.ID)
    if bo2.Filled != 1 || so2.Filled != 1 {
        t.Fatalf("filled mismatch: buy %v sell %v", bo2.Filled, so2.Filled)
    }
    if bo2.Status != "FILLED" || so2.Status != "FILLED" {
        t.Fatalf("status mismatch: buy %v sell %v", bo2.Status, so2.Status)
    }
}

func TestPartialFill(t *testing.T) {
    e := New()
    b := PlaceOrderRequest{Symbol: "ETH-USDT", Side: Buy, Price: 50, Size: 1}
    s := PlaceOrderRequest{Symbol: "ETH-USDT", Side: Sell, Price: 49, Size: 2}
    bo := e.PlaceOrder(b)
    so := e.PlaceOrder(s)
    bo2, _ := e.GetOrderStatus(bo.ID)
    so2, _ := e.GetOrderStatus(so.ID)
    if bo2.Filled != 1 || so2.Filled != 1 {
        t.Fatalf("filled mismatch: buy %v sell %v", bo2.Filled, so2.Filled)
    }
    if bo2.Status != "FILLED" || so2.Status != "OPEN" {
        t.Fatalf("status mismatch: buy %v sell %v", bo2.Status, so2.Status)
    }
}

func TestTopOfBook(t *testing.T) {
    e := New()
    e.PlaceOrder(PlaceOrderRequest{Symbol: "SOL-USDT", Side: Buy, Price: 20, Size: 1})
    e.PlaceOrder(PlaceOrderRequest{Symbol: "SOL-USDT", Side: Buy, Price: 22, Size: 1})
    e.PlaceOrder(PlaceOrderRequest{Symbol: "SOL-USDT", Side: Sell, Price: 25, Size: 1})
    e.PlaceOrder(PlaceOrderRequest{Symbol: "SOL-USDT", Side: Sell, Price: 24, Size: 1})
    q := e.Aggregate("SOL-USDT")
    if q.BestBid != 22 || q.BestAsk != 24 {
        t.Fatalf("quote mismatch: bid %v ask %v", q.BestBid, q.BestAsk)
    }
}
