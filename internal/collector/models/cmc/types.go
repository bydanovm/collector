package models_cmc

import (
	"time"
)

type QuotesLatest struct {
	Status status
	Data   map[string][]*CryptocurrencObject
}

type ListingLatest struct {
	Status status
	Data   []*CryptocurrencObject
}

type status struct {
	Timestamp     string
	Error_code    int
	Error_message string
	Elapsed       int
	Credit_count  int
	Notice        string
	Total_count   int
}

type CryptocurrencObject struct {
	Infinite_supply    bool
	Id                 int
	Num_market_pairs   int
	Is_active          int
	Cmc_rank           int
	Is_fiat            int
	Name               string
	Symbol             string
	Slug               string
	Date_added         time.Time
	Last_updated       time.Time
	Max_supply         float32
	Circulating_supply float32
	Total_supply       float32
	Quote              map[string]currency
}

func (c *CryptocurrencObject) GetId() int {
	return c.Id
}

func (c *CryptocurrencObject) GetNumMarketPairs() int {
	return c.Num_market_pairs
}

func (c *CryptocurrencObject) GetIsActive() int {
	return c.Is_active
}

func (c *CryptocurrencObject) GetCmcRank() int {
	return c.Cmc_rank
}

func (c *CryptocurrencObject) GetName() string {
	return c.Name
}

func (c *CryptocurrencObject) GetSymbol() string {
	return c.Symbol
}

func (c *CryptocurrencObject) GetSlug() string {
	return c.Slug
}

func (c *CryptocurrencObject) GetDateAdded() time.Time {
	return c.Date_added
}

func (c *CryptocurrencObject) GetLastUpdated() time.Time {
	return c.Last_updated
}

func (c *CryptocurrencObject) GetMaxSupply() float32 {
	return c.Max_supply
}

func (c *CryptocurrencObject) GetCirculatingSupply() float32 {
	return c.Circulating_supply
}

func (c *CryptocurrencObject) GetTotalSupply() float32 {
	return c.Total_supply
}

func (c *CryptocurrencObject) GetQuote() map[string]currency {
	return c.Quote
}

func (c *CryptocurrencObject) GetCurrency() string {
	return "USD"
}

func (c *CryptocurrencObject) GetPrice() float32 {
	return c.Quote["USD"].Price
}

func (c *CryptocurrencObject) GetVolume24h() float32 {
	return c.Quote["USD"].Volume_24h
}

func (c *CryptocurrencObject) GetMarketCap() float32 {
	return c.Quote["USD"].Market_cap
}

type currency struct {
	Price                    float32
	Volume_24h               float32
	Volume_change_24h        float32
	Percent_change_1h        float32
	Percent_change_24h       float32
	Percent_change_7d        float32
	Percent_change_30d       float32
	Percent_change_60d       float32
	Percent_change_90d       float32
	Market_cap               float32
	Market_cap_dominance     float32
	Fully_diluted_market_cap float32
	Last_updated             time.Time
}
