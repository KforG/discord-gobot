package util

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/KforG/discord-price-bot/config"
	"github.com/KforG/discord-price-bot/logging"
)

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

func TimeBeforeHalving(blocksRemain int64) string {
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

func GetBlockReward(blockHeight int) (currentBlockReward, nextBlockReward float64) {
	halvingInterval := 840000
	initialBlockReward := 50

	halvings := blockHeight / halvingInterval

	if halvings >= 32 { //Vertcoin and Bitcoin will undergo 32 halvings before zero emission
		return 0, 0
	}
	if halvings == 0 {
		return float64(initialBlockReward), (float64(initialBlockReward) / 2)
	}

	currentBlockReward = float64(initialBlockReward) / 2

	for i := 1; i < halvings; i++ {
		currentBlockReward = currentBlockReward / 2
	}

	return currentBlockReward, (currentBlockReward / 2)
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
