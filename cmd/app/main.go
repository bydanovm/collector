package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	retrievercoins "github.com/mbydanov/simple-miniapp-backend/internal/collector"
	models_cmc "github.com/mbydanov/simple-miniapp-backend/internal/collector/models/cmc"
	"github.com/mbydanov/simple-miniapp-backend/internal/db/pgsql"
	"github.com/mbydanov/simple-miniapp-backend/internal/keeper"
	logger "github.com/mbydanov/simple-miniapp-backend/internal/log"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// PgSql
	connStrPgSql := pgsql.NewConnStr().
		Host(os.Getenv("POSTGRES_HOST")).
		Port(os.Getenv("POSTGRES_PORT")).
		User(os.Getenv("POSTGRES_USER")).
		Password(os.Getenv("POSTGRES_PASSWORD")).
		Dbname(os.Getenv("POSTGRES_DB")).
		Sslmode(os.Getenv("POSTGRES_SSLMODE"))
	if err := connStrPgSql.GetError(); err != nil {
		log.Panic(err)
		return
	}
	schemaName := os.Getenv("POSTGRES_SCHEMA")
	dbPgSql := pgsql.NewPgSQL(connStrPgSql, &schemaName)
	if err := dbPgSql.GetError(); err != nil {
		log.Panic(err)
		return
	}
	dbPgSql.AutoMigrate()

	logger := slog.New(logger.NewDBHandler(dbPgSql))
	slog.SetDefault(logger)

	marketCMC := models_cmc.NewMarketCMC()
	retriever := retrievercoins.NewRetriever(ctx, marketCMC)

	keeper := keeper.NewKeeper(ctx, dbPgSql)

	go func() {
		<-time.After(5 * time.Second)
		go retriever.Run()
	}()

	go func() {
		for err := range retriever.GetError() {
			slog.Error(fmt.Errorf("%w", err).Error())
		}
	}()

	go func() {
		for {
			for _, data := range retriever.GetData() {
				keeper.Save(data)
			}

		}
	}()

	for range ctx.Done() {
		os.Exit(1)
	}

}
