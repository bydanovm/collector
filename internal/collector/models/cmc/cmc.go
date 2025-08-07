package models_cmc

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type MarketCMCImpl interface {
}

type MarketCMC struct {
	client *http.Client
}

func NewMarketCMC() *MarketCMC {
	return &MarketCMC{
		client: &http.Client{},
	}
}

func (m *MarketCMC) GetLatest(coins []string) {

	q := url.Values{}
	q.Add("symbol", "BTC,ETH")
	q.Add("convert", "USD")

	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		// s = append(s, "Возвращена ошибка:\n"+err.Error())
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("API_CMC"))
	req.URL.RawQuery = q.Encode()

	resp, err := m.client.Do(req)
	if err != nil {
		// s = append(s, "Возвращена ошибка:\n"+err.Error())
	}
	respBody, _ := io.ReadAll(resp.Body)
	qla := &QuotesLatestAnswer{}
	if err = json.Unmarshal([]byte(respBody), qla); err != nil {
		// s = append(s, "Возвращена ошибка:\n"+err.Error())
	}
	if qla.Error_code != 0 {
		// s = append(s, "Возвращена ошибка:\n"+qla.Error_message)
	}
	for range qla.QuotesLatestAnswerResults {
		// dateTime, err := models.ConvertDateTimeToMSK(qla.QuotesLatestAnswerResults[i].Last_updated)
		// if err != nil {
		// 	s = append(s, fmt.Sprintf("getAndSaveFromAPI:"+err.Error()))
		// }

		// str := fmt.Sprintf("Криптовалюта: %s\nЦена: %.9f %s\nОбновлено: %s",
		// 	qla.QuotesLatestAnswerResults[i].Symbol,
		// 	qla.QuotesLatestAnswerResults[i].Price,
		// 	qla.QuotesLatestAnswerResults[i].Currency,
		// 	dateTime,
		// )
		// s = append(s, str)

		// Добавление найденной валюты в БД текущих цен и справочник валют
		// cryptoprices := map[string]string{
		// 	"CryptoId":     fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Id),
		// 	"CryptoPrice":  fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Price),
		// 	"CryptoUpdate": dateTime,
		// }
		// dictCryptos := map[string]string{
		// 	"CryptoId":        fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Id),
		// 	"CryptoName":      fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Symbol),
		// 	"CryptoLastPrice": fmt.Sprintf("%v", qla.QuotesLatestAnswerResults[i].Price),
		// 	"CryptoUpdate":    dateTime,
		// }
		// if err := database.WriteData("dictcrypto", dictCryptos); err != nil {
		// 	s = append(s, fmt.Sprintf("GetLatest:"+err.Error()))
		// } else {
		// 	d := database.DictCrypto{
		// 		CryptoId:        qla.QuotesLatestAnswerResults[i].Id,
		// 		CryptoName:      qla.QuotesLatestAnswerResults[i].Symbol,
		// 		CryptoLastPrice: qla.QuotesLatestAnswerResults[i].Price,
		// 		CryptoCounter:   1,
		// 	}
		// 	database.DCCache[d.CryptoId] = d
		// }
		// if err := database.WriteData("cryptoprices", cryptoprices); err != nil {
		// 	s = append(s, "Возвращена ошибка:\n"+err.Error())
		// }
		// Поиск индекса найденной валюты и её удаление из массива needFind
		// needFind = models.FindCellAndDelete(needFind, qla.QuotesLatestAnswerResults[i].Symbol)

	}
	// Есть не найденная криптовалюта
	// if len(needFind) != 0 {
	// 	s = append(s, "Криптовалюта "+strings.Join(needFind, `, `)+" не найдена")
	// }
	// return s
}

type QuotesLatest struct {
	Status status
	Data   map[string][]cryptocurrencObject
}

type status struct {
	Timestamp     string
	ErrorCode     int
	Error_message string
	Elapsed       int
	CreditCount   int
	Notice        string
}

type cryptocurrencObject struct {
	Id                 int
	Name               string
	Symbol             string
	Slug               string
	Num_market_pairs   int
	Date_added         string
	Max_supply         float32
	Circulating_supply float32
	Total_supply       float32
	Is_active          int
	Infinite_supply    bool
	Cmc_rank           int
	Is_fiat            int
	Last_updated       string
	Quote              map[string]currency
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

type QuotesLatestAnswer struct {
	Error_code                int
	Error_message             string
	QuotesLatestAnswerResults []QuotesLatestAnswerResult
}

func (qla *QuotesLatestAnswer) UnmarshalJSON(bs []byte) error {
	var quotesLatest QuotesLatest
	if err := json.Unmarshal(bs, &quotesLatest); err != nil {
		return err
	}
	qla.Error_code = quotesLatest.Status.ErrorCode
	qla.Error_message = quotesLatest.Status.Error_message
	for _, value0 := range quotesLatest.Data {
		if len(value0) > 0 {
			qla.QuotesLatestAnswerResults = append(qla.QuotesLatestAnswerResults, QuotesLatestAnswerResult{
				Id:           value0[0].Id,
				Name:         value0[0].Name,
				Symbol:       value0[0].Symbol,
				Cmc_rank:     value0[0].Cmc_rank,
				Price:        value0[0].Quote["USD"].Price,
				Currency:     "USD",
				Last_updated: value0[0].Quote["USD"].Last_updated,
			})
		}
	}
	return nil
}

type QuotesLatestAnswerResult struct {
	Id           int
	Name         string
	Symbol       string
	Cmc_rank     int
	Price        float32
	Currency     string
	Last_updated time.Time
}
