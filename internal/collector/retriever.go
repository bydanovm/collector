package retrievercoins

import (
	"context"
	"time"
)

type MarketImpl interface {
	GetLatest(coins []string)
}

type RetrieverImpl interface {
	Run()
	End()
}

func NewRetriever(ctx context.Context, market MarketImpl) RetrieverImpl {
	r := &Retriever{
		ctx:     ctx,
		Market:  market,
		Timeout: 10,
		Start:   make(chan struct{}),
		Stop:    make(chan struct{}),
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
}

func (r *Retriever) Run() {
	close(r.Start)
}

func (r *Retriever) End() {
	r.Stop <- struct{}{}
}

func (r *Retriever) process() {
	go func() {
		for range r.Start {
		}
		ticker := time.NewTicker(r.Timeout * time.Second)
		for {
			select {
			case <-ticker.C:
				r.Market.GetLatest([]string{"BTC"})
			case <-r.Stop:
				return
			}
		}
	}()
}
