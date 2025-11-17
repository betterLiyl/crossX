package adapters

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"
)

type OKXAdapter struct{}

func (OKXAdapter) Name() string { return "OKX" }

func (OKXAdapter) BestBidAsk(symbol string) (Quote, error) {
	base := os.Getenv("OKX_BASE_URL")
	if base == "" {
		base = "https://www.okx.com"
	}
	url := base + "/api/v5/market/ticker?instId=" + symbol
	c := &http.Client{Timeout: 800 * time.Millisecond}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Quote{}, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return Quote{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Quote{}, errors.New("bad_status")
	}
	var out struct {
		Data []struct {
			BidPx string `json:"bidPx"`
			AskPx string `json:"askPx"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Quote{}, err
	}
	if len(out.Data) == 0 {
		return Quote{}, errors.New("empty")
	}
	bp, err := strconv.ParseFloat(out.Data[0].BidPx, 64)
	if err != nil {
		return Quote{}, err
	}
	ap, err := strconv.ParseFloat(out.Data[0].AskPx, 64)
	if err != nil {
		return Quote{}, err
	}
	return Quote{Symbol: symbol, BestBid: bp, BestAsk: ap, Source: "OKX"}, nil
}
