package util

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/KforG/discord-price-bot/config"
	"github.com/KforG/discord-price-bot/logging"
	"github.com/bwmarrin/discordgo"
)

func StartDiscordSession(token string) interface{} {
	//Create Discord session
	dg, err := discordgo.New(config.BotPrefix + config.Token)
	if err != nil {
		logging.Errorf("error creating new Discord sesssion,", err)
	}
	logging.Info("Successfully created new discord session")

	return dg
}

func GetJson(url string, target interface{}) error {
	//Fetch data from api
	resp, err := http.Get(url)
	if err != nil {
		logging.Errorf("Error fetching data from %s\n", url)
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		logging.Errorf("Error parsing response from %s\n", url)
		return err
	}
	return nil
}

func TimeBeforeHalving(blocksRemain int) string {
	seconds := blocksRemain * 150

	days := seconds / (3600 * 24)
	seconds = seconds % (3600 * 24)
	hours := seconds / (3600 * 24)
	seconds = seconds % 3600
	minutes := seconds / 60
	seconds = seconds % 60

	if days > 0 {
		return fmt.Sprintf("%d days %d hours %d minutes %d seconds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%d hours %d minutes %d seconds", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d minutes %d seconds", minutes, seconds)
}

type NetStatsResponse struct {
	BackendTipHeight int64
	Difficulty       float64
	TipHeight        int64
}

func GetNetworkStats() (int64, float64, error) {
	jsonPayload := NetStatsResponse{}
	err := GetJson(config.NetStatsApi, &jsonPayload)
	if err != nil {
		return jsonPayload.TipHeight, jsonPayload.Difficulty, err
	}
	return jsonPayload.TipHeight, jsonPayload.Difficulty, nil
}

// func taken from OCM, small modifications
func GetNethash(diff float64) uint64 {
	difficulty := big.NewFloat(diff)
	factor := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(48), nil)
	netHash := difficulty.Mul(difficulty, big.NewFloat(0).SetInt(factor))
	u, _ := netHash.Quo(netHash, big.NewFloat(9830250)).Uint64() // 0xffff * blocktime in seconds

	return u
}
