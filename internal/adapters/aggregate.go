package adapters

import (
	"sync"
	"time"
)

type AggResult struct {
	Symbol     string            `json:"symbol"`
	BestBid    float64           `json:"best_bid"`
	BestAsk    float64           `json:"best_ask"`
	BestBidSrc string            `json:"best_bid_src"`
	BestAskSrc string            `json:"best_ask_src"`
	Quotes     []Quote           `json:"quotes"`
	Errors     map[string]string `json:"errors"`
}

func AggregateBest(adapters []MarketDataAdapter, symbol string, timeout time.Duration) AggResult {
	res := AggResult{Symbol: symbol, Errors: map[string]string{}}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, a := range adapters {
		a := a
		wg.Add(1)
		go func() {
			defer wg.Done()
			q, err := a.BestBidAsk(symbol)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				res.Errors[a.Name()] = err.Error()
				return
			}
			res.Quotes = append(res.Quotes, q)
			if q.BestBid > res.BestBid || res.BestBid == 0 {
				res.BestBid = q.BestBid
				res.BestBidSrc = a.Name()
			}
			if (q.BestAsk < res.BestAsk || res.BestAsk == 0) && q.BestAsk > 0 {
				res.BestAsk = q.BestAsk
				res.BestAskSrc = a.Name()
			}
		}()
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(timeout):
	}
	return res
}
