package keeper

import (
	"context"
	"fmt"
	"log/slog"

	retrievercoins "github.com/mbydanov/simple-miniapp-backend/internal/collector"
	"github.com/mbydanov/simple-miniapp-backend/internal/models"
	"gorm.io/gorm"
)

type KeeperImpl interface {
	Save(data retrievercoins.DataImpl)
}

type OrmImpl interface {
	Db() *gorm.DB
}

type Keeper struct {
	ctx context.Context
	db  OrmImpl
}

func NewKeeper(ctx context.Context, db OrmImpl) KeeperImpl {
	k := &Keeper{
		ctx: ctx,
		db:  db,
	}

	return k
}

func (k *Keeper) Save(data retrievercoins.DataImpl) {
	coinStat := &models.CoinStat{
		IdCoin:      data.GetId(),
		CmcRank:     data.GetCmcRank(),
		Name:        data.GetName(),
		Symbol:      data.GetSymbol(),
		Price:       data.GetPrice(),
		Volume24h:   data.GetVolume24h(),
		MarketCap:   data.GetMarketCap(),
		Currency:    data.GetCurrency(),
		LastUpdated: data.GetLastUpdated(),
	}
	if err := k.db.Db().Create(coinStat).Error; err != nil {
		slog.Error(fmt.Errorf("%w", err).Error())
	}
}
