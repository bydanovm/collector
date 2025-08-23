package models

import "time"

type CoinStat struct {
	IdCoin          int       `gorm:"column:id_coin"`
	Name            string    `gorm:"column:name"`
	Symbol          string    `gorm:"column:symbol"`
	CmcRank         int       `gorm:"column:cmc_rank"`
	Price           float32   `gorm:"column:price"`
	Volume24h       float32   `gorm:"column:volume_24h"`
	MarketCap       float32   `gorm:"column:market_cap"`
	Currency        string    `gorm:"column:currency"`
	LastUpdated     time.Time `gorm:"column:last_updated"`
	LastUpdatedUnix int64     `gorm:"column:last_updated_unix"`
}
