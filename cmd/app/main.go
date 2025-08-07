package main

import (
	"context"
	"os"
	"time"

	retrievercoins "github.com/mbydanov/simple-miniapp-backend/internal/collector"
	models_cmc "github.com/mbydanov/simple-miniapp-backend/internal/collector/models/cmc"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	marketCMC := models_cmc.NewMarketCMC()
	retriever := retrievercoins.NewRetriever(ctx, marketCMC)

	go func() {
		<-time.After(5 * time.Second)
		go retriever.Run()
	}()

	for range ctx.Done() {
		os.Exit(1)
	}

}
