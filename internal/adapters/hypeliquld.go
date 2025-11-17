package adapters

import "errors"

type HypeliquldAdapter struct{}

func (HypeliquldAdapter) Name() string { return "Hypeliquld" }

func (HypeliquldAdapter) BestBidAsk(symbol string) (Quote, error) {
	return Quote{Symbol: symbol, BestBid: 0, BestAsk: 0, Source: "Hypeliquld"}, errors.New("not_implemented")
}
