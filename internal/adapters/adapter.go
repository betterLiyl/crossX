package adapters

type Quote struct {
	Symbol  string
	BestBid float64
	BestAsk float64
	Source  string
}

type MarketDataAdapter interface {
	Name() string
	BestBidAsk(symbol string) (Quote, error)
}
