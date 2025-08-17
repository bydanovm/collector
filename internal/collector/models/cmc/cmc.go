package models_cmc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	retrievercoins "github.com/mbydanov/simple-miniapp-backend/internal/collector"
	"github.com/mbydanov/simple-miniapp-backend/internal/utils"
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

func (m *MarketCMC) GetTopLatest(datac chan<- []retrievercoins.DataImpl, errc chan<- error) {
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		errc <- fmt.Errorf("%s:%w", utils.GetFunctionName(), err)
		return
	}
	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("API_CMC"))
	resp, err := m.client.Do(req)
	if err != nil {
		errc <- fmt.Errorf("%s:%w", utils.GetFunctionName(), err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	listingLatest := &ListingLatest{}
	if err = json.Unmarshal([]byte(respBody), listingLatest); err != nil {
		errc <- fmt.Errorf("%s:%w", utils.GetFunctionName(), err)
		return
	}
	if listingLatest.Status.Error_code != 0 {
		errc <- fmt.Errorf("%s:%s", utils.GetFunctionName(), listingLatest.Status.Error_message)
		return
	}

	var dataImplSlice []retrievercoins.DataImpl
	for _, item := range listingLatest.Data {
		dataImplSlice = append(dataImplSlice, item)
	}
	datac <- dataImplSlice
}

func (m *MarketCMC) GetSelectiveLatest(coins []string, datac chan<- []retrievercoins.DataImpl, errc chan<- error) {

	coinsStr := strings.Join(coins, ",")
	q := url.Values{}
	q.Add("symbol", coinsStr)
	q.Add("convert", "USD")

	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		errc <- fmt.Errorf("%s:%w", utils.GetFunctionName(), err)
		return
	}
	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("API_CMC"))
	req.URL.RawQuery = q.Encode()

	resp, err := m.client.Do(req)
	if err != nil {
		errc <- fmt.Errorf("%s:%w", utils.GetFunctionName(), err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	quotesLatest := &QuotesLatest{}
	if err = json.Unmarshal([]byte(respBody), quotesLatest); err != nil {
		errc <- fmt.Errorf("%s:%w", utils.GetFunctionName(), err)
		return
	}
	if quotesLatest.Status.Error_code != 0 {
		errc <- fmt.Errorf("%s:%s", utils.GetFunctionName(), quotesLatest.Status.Error_message)
		return
	}

	var dataImplSlice []retrievercoins.DataImpl
	for _, item := range quotesLatest.Data {
		for _, item := range item {
			dataImplSlice = append(dataImplSlice, item)
		}
	}
	datac <- dataImplSlice
}
