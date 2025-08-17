package retrievercoins

import (
	"context"
	"time"
)

type MarketImpl interface {
	GetTopLatest(datac chan<- []DataImpl, errc chan<- error)
	GetSelectiveLatest(coins []string, datac chan<- []DataImpl, errc chan<- error)
}

type DataImpl interface {
	GetId() int
	GetCmcRank() int
	GetName() string
	GetSymbol() string
	GetPrice() float32
	GetVolume24h() float32
	GetMarketCap() float32
	GetCurrency() string
	GetLastUpdated() time.Time
}

type RetrieverImpl interface {
	Run()
	End()
	GetError() <-chan error
	GetData() []DataImpl
}

func NewRetriever(ctx context.Context, market MarketImpl) RetrieverImpl {
	r := &Retriever{
		ctx:     ctx,
		Market:  market,
		Timeout: 240,
		Start:   make(chan struct{}),
		Stop:    make(chan struct{}),
		Errors:  make(chan error, 1),
		Data:    make(chan []DataImpl, 1),
	}
	go r.process()
	return r
}

type Retriever struct {
	ctx     context.Context
	Market  MarketImpl
	Timeout time.Duration
	Start   chan struct{}
	Stop    chan struct{}
	Errors  chan error
	Data    chan []DataImpl
}

func (r *Retriever) Run() {
	close(r.Start)
}

func (r *Retriever) End() {
	close(r.Stop)
}

func (r *Retriever) GetError() <-chan error {
	return r.Errors
}

func (r *Retriever) process() {
	go func() {
		for range r.Start {
		}
		ticker := time.NewTicker(r.Timeout * time.Second)
		for {
			select {
			case <-ticker.C:
				r.Market.GetTopLatest(r.Data, r.Errors)
			case <-r.Stop:
				return
			case <-r.ctx.Done():
				return
			}
		}
	}()
}

func (r *Retriever) GetData() []DataImpl {
	return <-r.Data
}
