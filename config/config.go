package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/KforG/discord-price-bot/logging"
)

var (
	Token       string
	BotPrefix   string
	PriceApi    string
	NetStatsApi string
)

type configStruct struct {
	Token       string
	BotPrefix   string
	PriceApi    string
	NetStatsApi string
}

func ReadConfig() error {
	logging.Infof("Reading config file...\n")
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		logging.Errorf("error reading config file! %s\n", err)
		return err
	}

	config := configStruct{}

	err = json.Unmarshal(file, &config)

	if err != nil {
		logging.Errorf("error unmarshaling file %s\n", err)
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix
	PriceApi = config.PriceApi
	NetStatsApi = config.NetStatsApi

	return nil
}
