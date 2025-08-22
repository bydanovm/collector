package broker

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	logger "github.com/mbydanov/simple-miniapp-backend/internal/log"
	"github.com/mbydanov/simple-miniapp-backend/internal/models"
	"github.com/mbydanov/simple-miniapp-backend/internal/utils"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

var broker Broker

func NewBroker(db *gorm.DB) {
	broker = *initBroker(db)
}
func GetBroker() *Broker {
	return &broker
}

type Broker struct {
	Broker *nats.Conn
	db     *gorm.DB
	log    logger.LoggerImpl
	err    error
}

func initBroker(db *gorm.DB) *Broker {
	natsConn, err := nats.Connect("nats://"+os.Getenv("NATS_HOST")+":"+os.Getenv("NATS_PORT")+"", nats.Token(os.Getenv("NATS_TOKEN")))
	if err != nil {
		slog.Error(fmt.Errorf("%s:%w", utils.GetFunctionName(), err).Error())
	}

	broker := &Broker{Broker: natsConn, db: db}
	broker.NewSubscribes()
	return broker
}

func (b *Broker) NewSubscribes() {
	b.Broker.Subscribe("tgbot.collector", GetCoin(b))
}

type MsgStruct struct {
	Coin     string
	Quantity int
	Time     uint64
}

func GetCoin(b *Broker) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		data := &MsgStruct{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			slog.Error(fmt.Errorf("%s:%w", utils.GetFunctionName(), err).Error())
		}

		if data.Quantity >= 0 && data.Quantity <= 1 {
			coinStat := models.CoinStat{}
			tx := b.db.Where("symbol = ?", data.Coin).Order("last_updated desc").First(&coinStat)
			if tx.Error != nil {
				slog.Error(fmt.Errorf("%s:%w", utils.GetFunctionName(), tx.Error).Error())
			}

			data, err := json.Marshal(&coinStat)
			if err != nil {
				slog.Error(fmt.Errorf("%s:%w", utils.GetFunctionName(), err).Error())
			}

			b.Broker.Publish("collector.tgbot", data)
		} else {
			coinStat := make([]models.CoinStat, 0, data.Quantity)
			tx := b.db.Where("symbol = ?", data.Coin).Order("last_updated desc").Limit(data.Quantity).Find(&coinStat)
			if tx.Error != nil {
				slog.Error(fmt.Errorf("%s:%w", utils.GetFunctionName(), tx.Error).Error())
			}

			data, err := json.Marshal(&coinStat)
			if err != nil {
				slog.Error(fmt.Errorf("%s:%w", utils.GetFunctionName(), err).Error())
			}

			b.Broker.Publish("collector.tgbot", data)
		}
	}
}
