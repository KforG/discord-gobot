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
type CoingeckoResponse struct {
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
	Price      string
	MarketCap  string
	NetHash    string
	Difficulty string
}

func UpdateChannelStats(dg *discordgo.Session) {
	for {
		channelPriceRefresh(dg)
		channelNetworkStatsRefresh(dg)
		time.Sleep(301 * time.Second) //Discord has rate-limited the number of times you can change channel names. Currently this is 1 request per 300 seconds.
	}

}

func channelNetworkStatsRefresh(dg *discordgo.Session) error {
	channel := ChannelID{
		NetHash:    "925295761887997992",
		Difficulty: "925295762257117184",
	}

	_, diff, err := util.GetNetworkStats()
	if err != nil {
		logging.Errorf("Error fetching netstatsinfo, skipping iteration... %s\n", err)
		return err
	}
	logging.Infof("Fetched new netstats, updating channel names\n")

	_, err = dg.ChannelEditComplex(channel.NetHash, &discordgo.ChannelEdit{
		Name:     fmt.Sprintf("NetHash: %0.2f GH/s", convertNethashToGiga(util.GetNethash(diff))),
		Position: 2,
	})
	if err != nil {
		logging.Errorf("Error updating NetHash channel %s\n", err)
	}
	_, err = dg.ChannelEditComplex(channel.Difficulty, &discordgo.ChannelEdit{
		Name:     fmt.Sprintf("NetDiff: %0.2f", diff),
		Position: 3,
	})
	if err != nil {
		logging.Errorf("Error updating Difficulty channel %s\n", err)
	}

	logging.Infof("NetStats channelnames update completed.. Sleeping for 5 minutes... \n")
	return nil
}

func channelPriceRefresh(dg *discordgo.Session) error {
	channel := ChannelID{
		Price:     "919741562295037982",
		MarketCap: "919741591407722626",
	}
	jsonPayload := CoingeckoResponse{}
	err := util.GetJson(config.PriceApi, &jsonPayload)
	if err != nil {
		logging.Errorf("Error fetching new data, skipping iteration... %s\n", err)
		return err
	}
	logging.Infof("Fetched new price info, updating channel names\n")

	_, err = dg.ChannelEditComplex(channel.Price, &discordgo.ChannelEdit{
		Name:     fmt.Sprintf("%s | %s", convertFiatPrice(jsonPayload.MarketData.CurrentPrice.Usd), convertBTCtoSats(jsonPayload.MarketData.CurrentPrice.Btc)),
		Position: 0,
	})
	if err != nil {
		logging.Errorf("Error updating price channel %s\n", err)
	}
	_, err = dg.ChannelEditComplex(channel.MarketCap, &discordgo.ChannelEdit{
		Name:     fmt.Sprintf("%s | ₿%0.0f", convertFiatMC(jsonPayload.MarketData.MarketCap.Usd), jsonPayload.MarketData.MarketCap.Btc),
		Position: 1,
	})
	if err != nil {
		logging.Errorf("Error updating marketcap channel %s\n", err)
	}

	logging.Infof("PriceStats channelnames update completed.. Sleeping for 5 minutes... \n")
	return nil
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

func convertNethashToGiga(u uint64) float64 {
	k := float64(u) / 1000000000
	return k
}
