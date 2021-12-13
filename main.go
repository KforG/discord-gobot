package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KforG/discord-price-bot/config"
	"github.com/KforG/discord-price-bot/functions"
	"github.com/KforG/discord-price-bot/logging"
	"github.com/bwmarrin/discordgo"
)

func main() {
	//init logging
	logFile, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logging.SetLogFile(logFile)

	logging.Infof("Discord bot started up! \n")

	//init config
	err := config.ReadConfig()
	if err != nil {
		logging.Errorf("error loading config.json %s\n", err)
		panic(err)
	}

	//Init session
	dg, err := discordgo.New(config.BotPrefix + config.Token)
	if err != nil {
		logging.Errorf("error creating new Discord sesssion, %s\n", err)
	}
	logging.Infof("Successfully created new discord session \n")

	//Open websocket to Discord
	err = dg.Open()
	if err != nil {
		logging.Errorf("error opening websocket connection, %s\n", err)
		panic(err)
	}

	//Display username of bot
	me, err := dg.User("@me")
	if err != nil {
		logging.Warnf("error obtaining account details, %s\n", err)
	}
	logging.Infof("We've logged in as %s\n", me)

	go functions.ChannelPriceRefresh(dg) //Update channelname price/MC/etc

	// Hold up program for Go routines and exit gracefully
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	logging.Infof("Input detected\n")

	//Shutdown
	dg.Close()
	logging.Infof("Closed connection.\n")
}
