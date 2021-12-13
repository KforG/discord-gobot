package functions

import (
	"fmt"
	"time"

	"github.com/KforG/discord-price-bot/config"
	"github.com/KforG/discord-price-bot/logging"
	"github.com/KforG/discord-price-bot/util"
	"github.com/bwmarrin/discordgo"
)

// some of the values are not used currently, so they're commented out to not waste memory space.
type Response struct {
	MarketData struct {
		CurrentPrice struct {
			Btc float64 `json:"btc"`
			//	Eur float64 `json:"eur"`
			//	Gbp float64 `json:"gbp"`
			Usd float64 `json:"usd"`
		} `json:"current_price"`
		MarketCap struct {
			Btc float64 `json:"btc"`
			//	Eur int     `json:"eur"`
			//	Gbp int     `json:"gbp"`
			Usd int `json:"usd"`
		} `json:"market_cap"`
		//MarketCapRank            int     `json:"market_cap_rank"`
		PriceChangePercentage24H float64 `json:"price_change_percentage_24h"`
		//PriceChangePercentage7D  float64 `json:"price_change_percentage_7d"`
		//PriceChangePercentage30D float64 `json:"price_change_percentage_30d"`
		PriceChange24HInCurrency struct {
			Btc float64 `json:"btc"`
			//	Eur float64 `json:"eur"`
			//	Gbp float64 `json:"gbp"`
			Usd float64 `json:"usd"`
		} `json:"price_change_24h_in_currency"`
		/*PriceChangePercentage1HInCurrency struct {
			Btc float64 `json:"btc"`
			Eur float64 `json:"eur"`
			Gbp float64 `json:"gbp"`
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_1h_in_currency"`
		PriceChangePercentage24HInCurrency struct {
			Btc float64 `json:"btc"`
		//	Eur float64 `json:"eur"`
		//	Gbp float64 `json:"gbp"`
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_24h_in_currency"`
		PriceChangePercentage7DInCurrency struct {
			Btc float64 `json:"btc"`
			Eur float64 `json:"eur"`
			Gbp float64 `json:"gbp"`
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_7d_in_currency"`   */
		//CirculatingSupply float64 `json:"circulating_supply"`
	} `json:"market_data"`
}

type ChannelID struct {
	Price     string
	MarketCap string
}

func ChannelPriceRefresh(dg *discordgo.Session) {

	channel := ChannelID{
		Price:     "919741562295037982",
		MarketCap: "919741591407722626",
	}
	for {
		jsonPayload := Response{}
		err := util.GetJson(config.Api, &jsonPayload)
		if err != nil {
			logging.Errorf("Error fetching new data, skipping iteration... %s\n", err)
			continue
		}
		logging.Infof("Fetched new price info, updating channel names\n")

		_, err = dg.ChannelEdit(channel.Price, fmt.Sprintf("%s | %s", convertFiatPrice(jsonPayload.MarketData.CurrentPrice.Usd), convertBTCtoSats(jsonPayload.MarketData.CurrentPrice.Btc)))
		if err != nil {
			logging.Errorf("Error updating price channel %s\n", err)
		}
		_, err = dg.ChannelEdit(channel.MarketCap, fmt.Sprintf("%s | ₿%0.0f", convertFiatMC(jsonPayload.MarketData.MarketCap.Usd), jsonPayload.MarketData.MarketCap.Btc))
		if err != nil {
			logging.Errorf("Error updating marketcap channel %s\n", err)
		}

		logging.Infof("Channelname update completed.. Sleeping for 5 minutes... \n")
		time.Sleep(301 * time.Second) //Discord has rate-limited the number of times you can change channel names. Currently this is 1 request per 300 seconds.
	}
}

func convertBTCtoSats(btc float64) string {
	sats := btc * 100000000
	return fmt.Sprintf("%d sats", int(sats))
}

func convertFiatPrice(fiat float64) string {
	return fmt.Sprintf("Price: $%0.2f", fiat)
}

func convertFiatMC(fiat int) string {
	fiatMil := fiat / 1000000 //I'm just assuming the MC won't go above or below a billion
	return fmt.Sprintf("MktCap: $%0.2dM", fiatMil)
}
