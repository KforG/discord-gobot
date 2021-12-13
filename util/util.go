package util

import (
	"encoding/json"
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
